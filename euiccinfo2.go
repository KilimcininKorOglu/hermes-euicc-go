// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

// EUICCInfo2 parser ported from the former fork's v2/euiccinfo2.go, adapted to
// the upstream damonto/euicc-go bertlv API. Upstream returns EUICCInfo2 as a raw
// TLV; this reconstructs the parsed struct the chip-info command consumes.

package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/damonto/euicc-go/bertlv"
	"github.com/damonto/euicc-go/bertlv/primitive"
)

var errUnexpectedTag = errors.New("unexpected tag")

// EUICCInfo2 represents the parsed EUICCInfo2 response from the eUICC.
type EUICCInfo2 struct {
	ProfileVersion                 string
	SVN                            string
	EUICCFirmwareVer               string
	ExtCardResource                ExtCardResource
	UICCCapability                 []string
	TS102241Version                string
	GlobalPlatformVersion          string
	RSPCapability                  []string
	EUICCCiPKIdListForVerification []string
	EUICCCiPKIdListForSigning      []string
	EUICCCategory                  string
	ForbiddenProfilePolicyRules    []string
	PPVersion                      string
	SASAccreditationNumber         string
	CertificationDataObject        CertificationDataObject
}

// ExtCardResource holds memory and application information from the eUICC.
type ExtCardResource struct {
	InstalledApplication  uint32
	FreeNonVolatileMemory uint32
	FreeVolatileMemory    uint32
}

// CertificationDataObject holds certification information for the eUICC.
type CertificationDataObject struct {
	PlatformLabel    string
	DiscoveryBaseURL string
}

// EUICCCategory represents the category of the eUICC.
type EUICCCategory int

const (
	EUICCCategoryOther       EUICCCategory = 0
	EUICCCategoryBasic       EUICCCategory = 1
	EUICCCategoryMedium      EUICCCategory = 2
	EUICCCategoryContactless EUICCCategory = 3
)

func (c EUICCCategory) String() string {
	switch c {
	case EUICCCategoryBasic:
		return "basicEuicc"
	case EUICCCategoryMedium:
		return "mediumEuicc"
	case EUICCCategoryContactless:
		return "contactlessEuicc"
	default:
		return "other"
	}
}

// UnmarshalBERTLV parses the EUICCInfo2 response from BER-TLV format.
func (e *EUICCInfo2) UnmarshalBERTLV(tlv *bertlv.TLV) error {
	// EUICCInfo2 is either tag 0xBF20 (32) or 0xBF22 (34).
	if !tlv.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 32) &&
		!tlv.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 34) {
		return errUnexpectedTag
	}

	for _, child := range tlv.Children {
		switch {
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 1): // 0x81 profileVersion
			e.ProfileVersion = parseVersion(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 2): // 0x82 svn
			e.SVN = parseVersion(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 3): // 0x83 euiccFirmwareVer
			e.EUICCFirmwareVer = parseVersion(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 4): // 0x84 extCardResource
			if err := e.ExtCardResource.UnmarshalBERTLV(child); err != nil {
				return err
			}
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 5): // 0x85 uiccCapability
			e.UICCCapability = parseUICCCapability(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 6): // 0x86 ts102241Version
			e.TS102241Version = parseVersion(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 7): // 0x87 globalplatformVersion
			e.GlobalPlatformVersion = parseVersion(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 8): // 0x88 rspCapability
			e.RSPCapability = parseRSPCapability(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 9): // 0xA9 euiccCiPKIdListForVerification
			e.EUICCCiPKIdListForVerification = parsePKIdList(child)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 10): // 0xAA euiccCiPKIdListForSigning
			e.EUICCCiPKIdListForSigning = parsePKIdList(child)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 11): // 0xAB euiccCategory
			if len(child.Value) > 0 {
				var cat int8
				child.UnmarshalValue(primitive.UnmarshalInt(&cat))
				e.EUICCCategory = EUICCCategory(cat).String()
			}
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 25): // 0x99 forbiddenProfilePolicyRules
			e.ForbiddenProfilePolicyRules = parseForbiddenPPR(child.Value)
		case child.Tag.If(bertlv.Universal, bertlv.Primitive, 4): // 0x04 ppVersion (OCTET STRING)
			e.PPVersion = parseVersion(child.Value)
		case child.Tag.If(bertlv.Universal, bertlv.Primitive, 12): // 0x0C sasAccreditationNumber (UTF8String)
			e.SASAccreditationNumber = string(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 12): // 0xAC certificationDataObject
			if err := e.CertificationDataObject.UnmarshalBERTLV(child); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnmarshalBERTLV parses the ExtCardResource from BER-TLV format.
func (e *ExtCardResource) UnmarshalBERTLV(tlv *bertlv.TLV) error {
	if !tlv.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 4) {
		return errUnexpectedTag
	}

	// ExtCardResource is an OCTET STRING wrapping nested TLVs; parse them manually.
	value := tlv.Value
	offset := 0
	for offset < len(value) {
		if offset >= len(value) {
			break
		}
		tag := value[offset]
		offset++
		if offset >= len(value) {
			break
		}
		length := int(value[offset])
		offset++
		if length&0x80 != 0 {
			numLengthBytes := length & 0x7F
			length = 0
			for i := 0; i < numLengthBytes && offset < len(value); i++ {
				length = (length << 8) | int(value[offset])
				offset++
			}
		}
		if offset+length > len(value) {
			break
		}
		fieldValue := value[offset : offset+length]
		offset += length

		switch tag {
		case 0x81: // installedApplication
			e.InstalledApplication = uint32(parseExtCardInt(fieldValue, 1))
		case 0x82: // freeNonVolatileMemory
			e.FreeNonVolatileMemory = uint32(parseExtCardInt(fieldValue, 2))
		case 0x83: // freeVolatileMemory
			e.FreeVolatileMemory = uint32(parseExtCardInt(fieldValue, 3))
		}
	}

	return nil
}

func parseExtCardInt(fieldValue []byte, tagNum uint8) int32 {
	if len(fieldValue) == 0 {
		return 0
	}
	var val int32
	tempTLV := bertlv.NewValue(bertlv.ContextSpecific.Primitive(uint64(tagNum)), fieldValue)
	tempTLV.UnmarshalValue(primitive.UnmarshalInt(&val))
	return val
}

// UnmarshalBERTLV parses the CertificationDataObject from BER-TLV format.
func (c *CertificationDataObject) UnmarshalBERTLV(tlv *bertlv.TLV) error {
	if !tlv.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 12) {
		return errUnexpectedTag
	}
	for _, child := range tlv.Children {
		switch {
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 0): // 0x80 platformLabel
			c.PlatformLabel = string(child.Value)
		case child.Tag.If(bertlv.ContextSpecific, bertlv.Primitive, 1): // 0x81 discoveryBaseURL
			c.DiscoveryBaseURL = string(child.Value)
		}
	}
	return nil
}

// parseVersion converts BER-encoded version bytes to a dotted string.
func parseVersion(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	if len(data) == 3 {
		return fmt.Sprintf("%d.%d.%d", data[0], data[1], data[2])
	}
	result := fmt.Sprintf("%d", data[0])
	for i := 1; i < len(data); i++ {
		result += fmt.Sprintf(".%d", data[i])
	}
	return result
}

func parseUICCCapability(data []byte) []string {
	return parseBitString(data, []string{
		"contactlessSupport", "usimSupport", "isimSupport", "csimSupport",
		"akaMilenage", "akaCave", "akaTuak128", "akaTuak256", "rfu1", "rfu2",
		"gbaAuthenUsim", "gbaAuthenISim", "mbmsAuthenUsim", "eapClient",
		"javacard", "multos", "multipleUsimSupport", "multipleIsimSupport",
		"multipleCsimSupport", "berTlvFileSupport", "dfLinkSupport", "catTp",
		"getIdentity", "profile-a-x25519", "profile-b-p256", "suciCalculatorApi",
	})
}

func parseRSPCapability(data []byte) []string {
	return parseBitString(data, []string{
		"additionalProfile", "crlSupport", "rpmSupport",
		"testProfileSupport", "deviceInfoExtensibilitySupport",
	})
}

func parseForbiddenPPR(data []byte) []string {
	return parseBitString(data, []string{"pprUpdateControl", "ppr1", "ppr2", "ppr3"})
}

// parseBitString converts a BER bit string to a list of set capability names.
func parseBitString(data []byte, names []string) []string {
	if len(data) < 2 {
		return nil
	}
	unusedBits := int(data[0])
	var result []string
	bitData := data[1:]
	bitIndex := 0
	totalBits := len(bitData)*8 - unusedBits
	for _, b := range bitData {
		for bitPos := 7; bitPos >= 0; bitPos-- {
			if bitIndex >= totalBits {
				break
			}
			if bitIndex < len(names) && (b&(1<<bitPos)) != 0 {
				result = append(result, names[bitIndex])
			}
			bitIndex++
		}
	}
	return result
}

// parsePKIdList parses a list of Public Key Identifiers as hex strings.
func parsePKIdList(tlv *bertlv.TLV) []string {
	if tlv == nil || len(tlv.Children) == 0 {
		return nil
	}
	var result []string
	for _, child := range tlv.Children {
		if len(child.Value) > 0 {
			result = append(result, hex.EncodeToString(child.Value))
		}
	}
	return result
}
