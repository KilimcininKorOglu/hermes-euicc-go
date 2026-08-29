// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

// Local reimplementation of the fork's chip-info convenience API on top of the
// upstream damonto/euicc-go lpa.Client. Upstream exposes only low-level ES10
// calls, so this aggregates EID, configured addresses, parsed EUICCInfo2, and
// the Rules Authorisation Table (RAT).

package main

import (
	"strings"

	"github.com/damonto/euicc-go/bertlv"
	"github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
)

// EUICCChipInfo aggregates comprehensive information about the eUICC chip.
type EUICCChipInfo struct {
	EID                     string
	ConfiguredAddresses     *lpa.EUICCConfiguredAddresses
	Info2                   *EUICCInfo2
	RulesAuthorisationTable []*RulesAuthorisationTable
}

// OperatorID represents a mobile network operator identifier.
type OperatorID struct {
	PLMN string
	GID1 string
	GID2 string
}

// UnmarshalBERTLV parses an OperatorID from BER-TLV format.
func (o *OperatorID) UnmarshalBERTLV(tlv *bertlv.TLV) error {
	for _, child := range tlv.Children {
		switch child.Tag.Value() {
		case 0x00: // mcc_mnc (PLMN), context-specific primitive 0 → 0x80
			o.PLMN = hexEncodeLower(child.Value)
		case 0x01: // gid1
			o.GID1 = hexEncodeLower(child.Value)
		case 0x02: // gid2
			o.GID2 = hexEncodeLower(child.Value)
		}
	}
	return nil
}

// RulesAuthorisationTable represents profile policy authorization rules.
type RulesAuthorisationTable struct {
	PPRIds           []string
	AllowedOperators []*OperatorID
	PPRFlags         []string
}

// UnmarshalBERTLV parses a RulesAuthorisationTable entry from BER-TLV format.
func (r *RulesAuthorisationTable) UnmarshalBERTLV(tlv *bertlv.TLV) error {
	pprNames := []string{"pprUpdateControl", "ppr1", "ppr2", "ppr3"}
	for _, child := range tlv.Children {
		switch child.Tag.Value() {
		case 0x00: // pprIds (0x80)
			r.PPRIds = parseBitString(child.Value, pprNames)
		case 0x01: // allowedOperators (0xA1)
			var operators []*OperatorID
			for _, opChild := range child.Children {
				op := new(OperatorID)
				if err := op.UnmarshalBERTLV(opChild); err != nil {
					return err
				}
				operators = append(operators, op)
			}
			r.AllowedOperators = operators
		case 0x02: // pprFlags (0x82)
			r.PPRFlags = parseBitString(child.Value, pprNames)
		}
	}
	return nil
}

// getRATRequest retrieves the Rules Authorisation Table (0xBF43).
type getRATRequest struct{}

func (r *getRATRequest) MarshalBERTLV() (*bertlv.TLV, error) {
	return bertlv.NewChildren(bertlv.ContextSpecific.Constructed(67)), nil // 0xBF43
}

func (r *getRATRequest) CardResponse() *getRATResponse {
	return new(getRATResponse)
}

type getRATResponse struct {
	RATList []*RulesAuthorisationTable
}

func (r *getRATResponse) UnmarshalBERTLV(tlv *bertlv.TLV) error {
	if !tlv.Tag.If(bertlv.ContextSpecific, bertlv.Constructed, 67) { // 0xBF43
		return errUnexpectedTag
	}
	ratContainer := tlv.First(bertlv.ContextSpecific.Constructed(0))
	if ratContainer == nil {
		return nil // Empty list is valid.
	}
	var rats []*RulesAuthorisationTable
	for _, child := range ratContainer.Children {
		rat := new(RulesAuthorisationTable)
		if err := rat.UnmarshalBERTLV(child); err != nil {
			return err
		}
		rats = append(rats, rat)
	}
	r.RATList = rats
	return nil
}

func (r *getRATResponse) Valid() error { return nil }

// getChipInfo aggregates comprehensive eUICC information. Only EID retrieval
// failure is fatal; optional data is best-effort.
func getChipInfo(client *lpa.Client) (*EUICCChipInfo, error) {
	var info EUICCChipInfo

	eidBytes, err := client.EID()
	if err != nil {
		return nil, err
	}
	info.EID = strings.ToUpper(hexEncodeLower(eidBytes))

	info.ConfiguredAddresses, _ = client.EUICCConfiguredAddresses()
	info.Info2, _ = getEUICCInfo2Parsed(client)
	info.RulesAuthorisationTable, _ = getRAT(client)

	return &info, nil
}

// getEUICCInfo2Parsed fetches the raw EUICCInfo2 TLV and parses it.
func getEUICCInfo2Parsed(client *lpa.Client) (*EUICCInfo2, error) {
	tlv, err := client.EUICCInfo2()
	if err != nil {
		return nil, err
	}
	var info EUICCInfo2
	if err := info.UnmarshalBERTLV(tlv); err != nil {
		return nil, err
	}
	return &info, nil
}

// getRAT retrieves the Rules Authorisation Table from the eUICC.
func getRAT(client *lpa.Client) ([]*RulesAuthorisationTable, error) {
	response, err := sgp22.InvokeAPDU(client.APDU, new(getRATRequest))
	if err != nil {
		return nil, err
	}
	return response.RATList, nil
}

// hexEncodeLower converts bytes to a lowercase hex string.
func hexEncodeLower(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	const hexChars = "0123456789abcdef"
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = hexChars[b>>4]
		result[i*2+1] = hexChars[b&0x0f]
	}
	return string(result)
}
