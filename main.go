// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/at"
	"github.com/KilimcininKorOglu/euicc-go/driver/mbim"
	"github.com/KilimcininKorOglu/euicc-go/driver/qmi"
	"github.com/KilimcininKorOglu/euicc-go/lpa"
	sgp22 "github.com/KilimcininKorOglu/euicc-go/v2"
)

// Response structures for JSON output
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type EIDResponse struct {
	EID string `json:"eid"`
}

type ProfileResponse struct {
	ICCID                string `json:"iccid"`
	ISDPAID              string `json:"isdp_aid,omitempty"`
	ProfileState         int    `json:"profile_state"`
	ProfileName          string `json:"profile_name,omitempty"`
	ProfileNickname      string `json:"profile_nickname,omitempty"`
	ServiceProviderName  string `json:"service_provider_name,omitempty"`
	ProfileClass         string `json:"profile_class,omitempty"`
	Icon                 string `json:"icon,omitempty"`
	IconFileType         string `json:"icon_file_type,omitempty"`
}

type NotificationResponse struct {
	SequenceNumber             int    `json:"sequence_number"`
	ProfileManagementOperation int    `json:"profile_management_operation"`
	Address                    string `json:"address,omitempty"`
	ICCID                      string `json:"iccid,omitempty"`
}

type DiscoveryResponse struct {
	EventID string `json:"event_id"`
	Address string `json:"address"`
}

type InfoResponse struct {
	EID        string `json:"eid"`
	EUICCInfo1 string `json:"euicc_info1"`
	EUICCInfo2 string `json:"euicc_info2"`
}

type DownloadResponse struct {
	ISDPAID      string `json:"isdp_aid"`
	Notification int    `json:"notification"`
}

type ConfiguredAddressesResponse struct {
	DefaultSMDPAddress string `json:"default_smdp_address,omitempty"`
	RootSMDSAddress    string `json:"root_smds_address,omitempty"`
}

type ProcessedNotification struct {
	SequenceNumber int  `json:"sequence_number"`
	Removed        bool `json:"removed"`
}

type FailedNotification struct {
	SequenceNumber int    `json:"sequence_number"`
	Error          string `json:"error"`
}

type AutoNotificationResponse struct {
	Message       string                   `json:"message"`
	Total         int                      `json:"total"`
	Processed     int                      `json:"processed"`
	Failed        int                      `json:"failed"`
	ProcessedList []ProcessedNotification  `json:"processed_list"`
	FailedList    []FailedNotification     `json:"failed_list"`
}

type ChipInfoResponse struct {
	EID                     string                       `json:"eid"`
	ConfiguredAddresses     *ConfiguredAddressesResponse `json:"configured_addresses,omitempty"`
	Info2                   *EUICCInfo2Response          `json:"euicc_info2,omitempty"`
	RulesAuthorisationTable []RATResponse                `json:"rules_authorisation_table,omitempty"`
}

type EUICCInfo2Response struct {
	// Version Information
	ProfileVersion        string `json:"profile_version,omitempty"`
	SVN                   string `json:"svn,omitempty"`
	EUICCFirmwareVer      string `json:"euicc_firmware_ver,omitempty"`
	TS102241Version       string `json:"ts102241_version,omitempty"`
	GlobalPlatformVersion string `json:"global_platform_version,omitempty"`
	PPVersion             string `json:"pp_version,omitempty"`

	// Memory/Storage Information
	ExtCardResource ExtCardResourceResponse `json:"ext_card_resource"`

	// Capabilities
	UICCCapability []string `json:"uicc_capability,omitempty"`
	RSPCapability  []string `json:"rsp_capability,omitempty"`

	// Security
	EUICCCiPKIdListForVerification []string `json:"euicc_ci_pkid_list_for_verification,omitempty"`
	EUICCCiPKIdListForSigning      []string `json:"euicc_ci_pkid_list_for_signing,omitempty"`
	ForbiddenProfilePolicyRules    []string `json:"forbidden_profile_policy_rules,omitempty"`

	// Classification
	EUICCCategory string `json:"euicc_category,omitempty"`

	// Certification
	SASAccreditationNumber string                          `json:"sas_accreditation_number,omitempty"`
	CertificationDataObject CertificationDataObjectResponse `json:"certification_data_object,omitempty"`
}

type ExtCardResourceResponse struct {
	InstalledApplication  uint32 `json:"installed_application"`
	FreeNonVolatileMemory uint32 `json:"free_non_volatile_memory"`
	FreeVolatileMemory    uint32 `json:"free_volatile_memory"`
}

type CertificationDataObjectResponse struct {
	PlatformLabel    string `json:"platform_label,omitempty"`
	DiscoveryBaseURL string `json:"discovery_base_url,omitempty"`
}

type RATResponse struct {
	PPRIds           []string                  `json:"ppr_ids,omitempty"`
	AllowedOperators []AllowedOperatorResponse `json:"allowed_operators,omitempty"`
}

type AllowedOperatorResponse struct {
	PLMN string `json:"plmn,omitempty"`
	GID1 string `json:"gid1,omitempty"`
	GID2 string `json:"gid2,omitempty"`
}

// Global flags
var (
	devicePath  = flag.String("device", "", "Device path (e.g., /dev/cdc-wdm0, /dev/ttyUSB2)")
	driverType  = flag.String("driver", "", "Driver type: qmi, mbim, at, ccid (auto-detect if not specified)")
	slotNumber  = flag.Int("slot", 1, "SIM slot number")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
	timeout     = flag.Int("timeout", 30, "HTTP timeout in seconds")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	// Handle commands that don't need client
	switch command {
	case "help":
		printUsage()
		return
	case "version":
		handleVersion()
		return
	}

	// Validate command before initializing client
	validCommands := map[string]bool{
		"eid":                   true,
		"info":                  true,
		"chip-info":             true,
		"list":                  true,
		"enable":                true,
		"disable":               true,
		"delete":                true,
		"nickname":              true,
		"download":              true,
		"discovery":             true,
		"discover-download":     true,
		"notifications":         true,
		"notification-remove":   true,
		"notification-handle":   true,
		"auto-notification":     true,
		"notification-process":  true,
		"configured-addresses":  true,
		"set-default-dp":        true,
		"challenge":             true,
		"memory-reset":          true,
	}

	if !validCommands[command] {
		outputError(fmt.Errorf("unknown command: %s", command))
		fmt.Fprintln(os.Stderr, "\nRun 'hermes-euicc help' for usage information.")
		os.Exit(1)
	}

	// Initialize LPA client
	client, err := initClient()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}
	defer client.Close()

	// Execute command
	switch command {
	case "eid":
		handleEID(client)
	case "info":
		handleInfo(client)
	case "chip-info":
		handleChipInfo(client)
	case "list":
		handleList(client)
	case "enable":
		handleEnable(client)
	case "disable":
		handleDisable(client)
	case "delete":
		handleDelete(client)
	case "nickname":
		handleNickname(client)
	case "download":
		handleDownload(client)
	case "discovery":
		handleDiscovery(client)
	case "discover-download":
		handleDiscoverDownload(client)
	case "notifications":
		handleNotifications(client)
	case "notification-remove":
		handleNotificationRemove(client)
	case "notification-handle":
		handleNotificationHandle(client)
	case "auto-notification":
		handleAutoNotification(client)
	case "notification-process":
		handleNotificationProcess(client)
	case "configured-addresses":
		handleConfiguredAddresses(client)
	case "set-default-dp":
		handleSetDefaultDP(client)
	case "challenge":
		handleChallenge(client)
	case "memory-reset":
		handleMemoryReset(client)
	default:
		outputError(fmt.Errorf("unknown command: %s", command))
		printUsage()
		os.Exit(1)
	}
}

func initClient() (*lpa.Client, error) {
	var channel apdu.SmartCardChannel
	var err error

	if *driverType != "" {
		// User specified driver
		channel, err = createDriver(*driverType, *devicePath, *slotNumber)
	} else {
		// Auto-detect driver
		channel, err = autoDetectDriver(*devicePath, *slotNumber)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize driver: %w", err)
	}

	opts := &lpa.Options{
		Channel: channel,
		Timeout: time.Duration(*timeout) * time.Second,
	}

	return lpa.New(opts)
}

func createDriver(driverName, device string, slot int) (apdu.SmartCardChannel, error) {
	switch driverName {
	case "qmi":
		if device == "" {
			device = "/dev/cdc-wdm0"
		}
		return qmi.New(device, uint8(slot))
	case "mbim":
		if device == "" {
			device = "/dev/cdc-wdm0"
		}
		return mbim.New(device, uint8(slot))
	case "at":
		if device == "" {
			return nil, fmt.Errorf("device path required for AT driver")
		}
		return at.New(device)
	case "ccid":
		return initCCIDDriver()
	default:
		return nil, fmt.Errorf("unknown driver type: %s", driverName)
	}
}

func autoDetectDriver(device string, slot int) (apdu.SmartCardChannel, error) {
	// Try QMI
	if device == "" || device == "/dev/cdc-wdm0" {
		if ch, err := qmi.New("/dev/cdc-wdm0", uint8(slot)); err == nil {
			if *verbose {
				log.Println("Auto-detected: QMI driver")
			}
			return ch, nil
		}
	}

	// Try MBIM
	if device == "" || device == "/dev/cdc-wdm0" {
		if ch, err := mbim.New("/dev/cdc-wdm0", uint8(slot)); err == nil {
			if *verbose {
				log.Println("Auto-detected: MBIM driver")
			}
			return ch, nil
		}
	}

	// Try AT on common devices
	atDevices := []string{"/dev/ttyUSB2", "/dev/ttyUSB3", "/dev/ttyUSB1"}
	if device != "" {
		atDevices = []string{device}
	}
	for _, dev := range atDevices {
		if ch, err := at.New(dev); err == nil {
			if *verbose {
				log.Printf("Auto-detected: AT driver on %s\n", dev)
			}
			return ch, nil
		}
	}

	// Try CCID (only on supported platforms)
	if ccidSupported {
		if ch, err := initCCIDDriver(); err == nil {
			if *verbose {
				log.Printf("Auto-detected: CCID driver\n")
			}
			return ch, nil
		}
	}

	return nil, fmt.Errorf("no compatible driver found")
}

// Command handlers

func handleVersion() {
	outputSuccess(map[string]string{
		"name":      "Hermes eUICC Manager",
		"version":   "1.0.0",
		"copyright": "Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>",
		"license":   "MIT",
	})
}

func handleEID(client *lpa.Client) {
	eid, err := client.EID()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(EIDResponse{
		EID: hex.EncodeToString(eid),
	})
}

func handleInfo(client *lpa.Client) {
	eid, err := client.EID()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	info1, err := client.EUICCInfo1()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	info2, err := client.EUICCInfo2()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(InfoResponse{
		EID:        hex.EncodeToString(eid),
		EUICCInfo1: hex.EncodeToString(info1.Bytes()),
		EUICCInfo2: hex.EncodeToString(info2.Bytes()),
	})
}

func handleChipInfo(client *lpa.Client) {
	// Get chip info using library's ChipInfo function
	chipInfo, err := client.ChipInfo()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	// Build response
	response := ChipInfoResponse{
		EID: chipInfo.EID,
	}

	// Add configured addresses if available
	if chipInfo.ConfiguredAddresses != nil {
		response.ConfiguredAddresses = &ConfiguredAddressesResponse{
			DefaultSMDPAddress: chipInfo.ConfiguredAddresses.DefaultSMDPAddress,
			RootSMDSAddress:    chipInfo.ConfiguredAddresses.RootSMDSAddress,
		}
	}

	// Add Info2 if available
	if chipInfo.Info2 != nil {
		response.Info2 = &EUICCInfo2Response{
			ProfileVersion:        chipInfo.Info2.ProfileVersion,
			SVN:                   chipInfo.Info2.SVN,
			EUICCFirmwareVer:      chipInfo.Info2.EUICCFirmwareVer,
			TS102241Version:       chipInfo.Info2.TS102241Version,
			GlobalPlatformVersion: chipInfo.Info2.GlobalPlatformVersion,
			PPVersion:             chipInfo.Info2.PPVersion,
			ExtCardResource: ExtCardResourceResponse{
				InstalledApplication:  chipInfo.Info2.ExtCardResource.InstalledApplication,
				FreeNonVolatileMemory: chipInfo.Info2.ExtCardResource.FreeNonVolatileMemory,
				FreeVolatileMemory:    chipInfo.Info2.ExtCardResource.FreeVolatileMemory,
			},
			UICCCapability:                 chipInfo.Info2.UICCCapability,
			RSPCapability:                  chipInfo.Info2.RSPCapability,
			EUICCCiPKIdListForVerification: chipInfo.Info2.EUICCCiPKIdListForVerification,
			EUICCCiPKIdListForSigning:      chipInfo.Info2.EUICCCiPKIdListForSigning,
			ForbiddenProfilePolicyRules:    chipInfo.Info2.ForbiddenProfilePolicyRules,
			EUICCCategory:                  chipInfo.Info2.EUICCCategory,
			SASAccreditationNumber:         chipInfo.Info2.SASAccreditationNumber,
			CertificationDataObject: CertificationDataObjectResponse{
				PlatformLabel:    chipInfo.Info2.CertificationDataObject.PlatformLabel,
				DiscoveryBaseURL: chipInfo.Info2.CertificationDataObject.DiscoveryBaseURL,
			},
		}
	}

	// Add RAT if available
	if len(chipInfo.RulesAuthorisationTable) > 0 {
		response.RulesAuthorisationTable = make([]RATResponse, 0, len(chipInfo.RulesAuthorisationTable))
		for _, rat := range chipInfo.RulesAuthorisationTable {
			ratResp := RATResponse{
				PPRIds: rat.PPRIds,
			}

			if len(rat.AllowedOperators) > 0 {
				ratResp.AllowedOperators = make([]AllowedOperatorResponse, 0, len(rat.AllowedOperators))
				for _, op := range rat.AllowedOperators {
					ratResp.AllowedOperators = append(ratResp.AllowedOperators, AllowedOperatorResponse{
						PLMN: op.PLMN,
						GID1: op.GID1,
						GID2: op.GID2,
					})
				}
			}

			response.RulesAuthorisationTable = append(response.RulesAuthorisationTable, ratResp)
		}
	}

	outputSuccess(response)
}

func handleList(client *lpa.Client) {
	profiles, err := client.ListProfile(nil, nil)
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	response := make([]ProfileResponse, 0, len(profiles))
	for _, p := range profiles {
		pr := ProfileResponse{
			ICCID:               p.ICCID.String(),
			ISDPAID:             p.ISDPAID.String(),
			ProfileState:        int(p.ProfileState),
			ProfileName:         p.ProfileName,
			ProfileNickname:     p.ProfileNickname,
			ServiceProviderName: p.ServiceProviderName,
			ProfileClass:        p.ProfileClass.String(),
		}
		if p.Icon.Valid() {
			pr.Icon = p.Icon.String()
			pr.IconFileType = p.Icon.FileType()
		}
		response = append(response, pr)
	}

	outputSuccess(response)
}

func handleEnable(client *lpa.Client) {
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("usage: enable <iccid>"))
		os.Exit(1)
	}

	iccid, err := sgp22.NewICCID(flag.Arg(1))
	if err != nil {
		outputError(fmt.Errorf("invalid ICCID: %w", err))
		os.Exit(1)
	}

	if err := client.EnableProfile(iccid, true); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"message": "profile enabled successfully",
		"iccid":   flag.Arg(1),
	})
}

func handleDisable(client *lpa.Client) {
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("usage: disable <iccid>"))
		os.Exit(1)
	}

	iccid, err := sgp22.NewICCID(flag.Arg(1))
	if err != nil {
		outputError(fmt.Errorf("invalid ICCID: %w", err))
		os.Exit(1)
	}

	if err := client.DisableProfile(iccid, true); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"message": "profile disabled successfully",
		"iccid":   flag.Arg(1),
	})
}

func handleDelete(client *lpa.Client) {
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("usage: delete <iccid>"))
		os.Exit(1)
	}

	iccid, err := sgp22.NewICCID(flag.Arg(1))
	if err != nil {
		outputError(fmt.Errorf("invalid ICCID: %w", err))
		os.Exit(1)
	}

	if err := client.DeleteProfile(iccid); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"message": "profile deleted successfully",
		"iccid":   flag.Arg(1),
	})
}

func handleNickname(client *lpa.Client) {
	if flag.NArg() < 3 {
		outputError(fmt.Errorf("usage: nickname <iccid> <nickname>"))
		os.Exit(1)
	}

	iccid, err := sgp22.NewICCID(flag.Arg(1))
	if err != nil {
		outputError(fmt.Errorf("invalid ICCID: %w", err))
		os.Exit(1)
	}

	nickname := flag.Arg(2)

	if err := client.SetNickname(iccid, nickname); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"message":  "nickname set successfully",
		"iccid":    flag.Arg(1),
		"nickname": nickname,
	})
}

func handleDownload(client *lpa.Client) {
	var (
		activationCode   = flag.String("code", "", "Activation code (LPA:1$smdp.io$MATCHING-ID)")
		confirmationCode = flag.String("confirmation-code", "", "Confirmation code")
		imei             = flag.String("imei", "", "IMEI")
		autoConfirm      = flag.Bool("confirm", false, "Auto-confirm download")
	)

	downloadFlags := flag.NewFlagSet("download", flag.ExitOnError)
	downloadFlags.StringVar(activationCode, "code", "", "Activation code")
	downloadFlags.StringVar(confirmationCode, "confirmation-code", "", "Confirmation code")
	downloadFlags.StringVar(imei, "imei", "", "IMEI")
	downloadFlags.BoolVar(autoConfirm, "confirm", false, "Auto-confirm download")
	downloadFlags.Parse(flag.Args()[1:])

	if *activationCode == "" {
		outputError(fmt.Errorf("activation code required: use --code"))
		os.Exit(1)
	}

	ac := &lpa.ActivationCode{}
	if err := ac.UnmarshalText([]byte(*activationCode)); err != nil {
		outputError(fmt.Errorf("invalid activation code: %w", err))
		os.Exit(1)
	}

	if *imei != "" {
		ac.IMEI = *imei
	}

	ctx := context.Background()
	opts := &lpa.DownloadOptions{
		OnProgress: func(stage lpa.DownloadStage) {
			if *verbose {
				log.Printf("Download stage: %v\n", stage)
			}
		},
		OnConfirm: func(metadata *sgp22.ProfileInfo) bool {
			return *autoConfirm
		},
		OnEnterConfirmationCode: func() string {
			return *confirmationCode
		},
	}

	result, err := client.DownloadProfile(ctx, ac, opts)
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	dr := DownloadResponse{
		ISDPAID: result.ISDPAID().String(),
	}
	if result.Notification != nil {
		dr.Notification = int(result.Notification.ProfileManagementOperation)
	}

	outputSuccess(dr)
}

func handleDiscovery(client *lpa.Client) {
	var (
		server = flag.String("server", "", "SM-DS server address (default: lpa.ds.gsma.com)")
		imei   = flag.String("imei", "", "IMEI for authentication")
	)

	discoveryFlags := flag.NewFlagSet("discovery", flag.ExitOnError)
	discoveryFlags.StringVar(server, "server", "", "SM-DS server address")
	discoveryFlags.StringVar(imei, "imei", "", "IMEI")
	discoveryFlags.Parse(flag.Args()[1:])

	// Prepare discovery options
	opts := &lpa.DiscoverProfilesOptions{}

	// Set SM-DS address if provided
	if *server != "" {
		opts.SMDSAddress = *server
	}

	// Set IMEI if provided
	if *imei != "" {
		imeiBytes, err := sgp22.NewIMEI(*imei)
		if err != nil {
			outputError(fmt.Errorf("invalid IMEI: %w", err))
			os.Exit(1)
		}
		opts.IMEI = imeiBytes
	}

	// Use library's DiscoverProfiles function
	profiles, err := client.DiscoverProfiles(opts)
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	// Convert to response format
	response := make([]DiscoveryResponse, len(profiles))
	for i, profile := range profiles {
		response[i] = DiscoveryResponse{
			EventID: profile.EventID,
			Address: profile.SMDPAddress,
		}
	}

	outputSuccess(response)
}

func handleDiscoverDownload(client *lpa.Client) {
	var (
		server = flag.String("server", "", "SM-DS server address (default: lpa.ds.gsma.com)")
		imei   = flag.String("imei", "", "IMEI for authentication")
	)

	discoveryFlags := flag.NewFlagSet("discover-download", flag.ExitOnError)
	discoveryFlags.StringVar(server, "server", "", "SM-DS server address")
	discoveryFlags.StringVar(imei, "imei", "", "IMEI")
	discoveryFlags.Parse(flag.Args()[1:])

	// Prepare discovery options
	discoveryOpts := &lpa.DiscoverProfilesOptions{}

	// Set SM-DS address if provided
	if *server != "" {
		discoveryOpts.SMDSAddress = *server
	}

	// Set IMEI if provided
	if *imei != "" {
		imeiBytes, err := sgp22.NewIMEI(*imei)
		if err != nil {
			outputError(fmt.Errorf("invalid IMEI: %w", err))
			os.Exit(1)
		}
		discoveryOpts.IMEI = imeiBytes
	}

	// Use library's DiscoverAndDownload function
	ctx := context.Background()
	result, err := client.DiscoverAndDownload(ctx, discoveryOpts, nil)
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	// Check if a profile was downloaded
	if result == nil {
		outputSuccess(map[string]interface{}{
			"message": "no profiles available for download",
		})
		return
	}

	outputSuccess(map[string]interface{}{
		"message": "profile downloaded successfully",
	})
}

func handleNotifications(client *lpa.Client) {
	notifications, err := client.ListNotification()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	response := make([]NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		response = append(response, NotificationResponse{
			SequenceNumber:             int(n.SequenceNumber),
			ProfileManagementOperation: int(n.ProfileManagementOperation),
			Address:                    n.Address,
			ICCID:                      n.ICCID.String(),
		})
	}

	outputSuccess(response)
}

func handleNotificationRemove(client *lpa.Client) {
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("usage: notification-remove <sequence-number>"))
		os.Exit(1)
	}

	var seqNum int
	if _, err := fmt.Sscanf(flag.Arg(1), "%d", &seqNum); err != nil {
		outputError(fmt.Errorf("invalid sequence number: %w", err))
		os.Exit(1)
	}

	if err := client.RemoveNotificationFromList(sgp22.SequenceNumber(seqNum)); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]interface{}{
		"message":         "notification removed successfully",
		"sequence_number": seqNum,
	})
}

func handleNotificationHandle(client *lpa.Client) {
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("usage: notification-handle <sequence-number>"))
		os.Exit(1)
	}

	var seqNum int
	if _, err := fmt.Sscanf(flag.Arg(1), "%d", &seqNum); err != nil {
		outputError(fmt.Errorf("invalid sequence number: %w", err))
		os.Exit(1)
	}

	notifications, err := client.RetrieveNotificationList(sgp22.SequenceNumber(seqNum))
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	if len(notifications) == 0 {
		outputError(fmt.Errorf("notification not found"))
		os.Exit(1)
	}

	if err := client.HandleNotification(notifications[0]); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]interface{}{
		"message":         "notification handled successfully",
		"sequence_number": seqNum,
	})
}

func handleAutoNotification(client *lpa.Client) {
	// Use the library's ProcessAllNotifications function
	results, err := client.ProcessAllNotifications(&lpa.ProcessNotificationsOptions{
		AutoRemove:      true,
		ContinueOnError: true,
	})
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	// Convert results to response format
	processed := make([]ProcessedNotification, 0)
	failed := make([]FailedNotification, 0)

	for _, result := range results {
		if result.Success {
			processed = append(processed, ProcessedNotification{
				SequenceNumber: int(result.SequenceNumber),
				Removed:        result.Removed,
			})
		} else {
			failed = append(failed, FailedNotification{
				SequenceNumber: int(result.SequenceNumber),
				Error:          result.Error.Error(),
			})
		}
	}

	outputSuccess(AutoNotificationResponse{
		Message:       "auto notification processing completed",
		Total:         len(results),
		Processed:     len(processed),
		Failed:        len(failed),
		ProcessedList: processed,
		FailedList:    failed,
	})
}

func handleNotificationProcess(client *lpa.Client) {
	// Get sequence numbers from arguments
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("sequence number(s) required"))
		os.Exit(1)
	}

	// Parse all sequence numbers from arguments
	var sequenceNumbers []sgp22.SequenceNumber
	for i := 1; i < flag.NArg(); i++ {
		seqNum, err := strconv.Atoi(flag.Arg(i))
		if err != nil {
			outputError(fmt.Errorf("invalid sequence number '%s': %w", flag.Arg(i), err))
			os.Exit(1)
		}
		sequenceNumbers = append(sequenceNumbers, sgp22.SequenceNumber(seqNum))
	}

	// Use the library's ProcessNotifications function
	results, err := client.ProcessNotifications(
		&lpa.ProcessNotificationsOptions{
			AutoRemove:      true,
			ContinueOnError: true,
		},
		sequenceNumbers...,
	)
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	// Convert results to response format
	processed := make([]ProcessedNotification, 0)
	failed := make([]FailedNotification, 0)

	for _, result := range results {
		if result.Success {
			processed = append(processed, ProcessedNotification{
				SequenceNumber: int(result.SequenceNumber),
				Removed:        result.Removed,
			})
		} else {
			failed = append(failed, FailedNotification{
				SequenceNumber: int(result.SequenceNumber),
				Error:          result.Error.Error(),
			})
		}
	}

	outputSuccess(AutoNotificationResponse{
		Message:       "notification processing completed",
		Total:         len(results),
		Processed:     len(processed),
		Failed:        len(failed),
		ProcessedList: processed,
		FailedList:    failed,
	})
}

func handleConfiguredAddresses(client *lpa.Client) {
	addresses, err := client.EUICCConfiguredAddresses()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(ConfiguredAddressesResponse{
		DefaultSMDPAddress: addresses.DefaultSMDPAddress,
		RootSMDSAddress:    addresses.RootSMDSAddress,
	})
}

func handleSetDefaultDP(client *lpa.Client) {
	if flag.NArg() < 2 {
		outputError(fmt.Errorf("usage: set-default-dp <address>"))
		os.Exit(1)
	}

	address := flag.Arg(1)
	if err := client.SetDefaultDPAddress(address); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"message": "default DP address set successfully",
		"address": address,
	})
}

func handleChallenge(client *lpa.Client) {
	challenge, err := client.EUICCChallenge()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"challenge": hex.EncodeToString(challenge),
	})
}

func handleMemoryReset(client *lpa.Client) {
	if err := client.MemoryReset(); err != nil {
		outputError(err)
		os.Exit(1)
	}

	outputSuccess(map[string]string{
		"message": "memory reset successfully",
	})
}

// Output helpers

func outputSuccess(data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	outputJSON(response)
}

func outputError(err error) {
	response := Response{
		Success: false,
		Error:   err.Error(),
	}
	outputJSON(response)
}

func outputJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Hermes eUICC Manager - JSON-based eSIM Management CLI

Usage: %s [options] <command> [command-options]

Global Options:
  -device string
        Device path (e.g., /dev/cdc-wdm0, /dev/ttyUSB2)
  -driver string
        Driver type: qmi, mbim, at, ccid (auto-detect if not specified)
  -slot int
        SIM slot number (default 1)
  -timeout int
        HTTP timeout in seconds (default 30)
  -verbose
        Enable verbose logging

Commands:
  help                          Show this help message
  version                       Show version information
  eid                           Get EID
  info                          Get eUICC information (EID + EUICCInfo1 + EUICCInfo2)
  chip-info                     Get detailed chip information (parsed, includes memory/capabilities)
  list                          List all profiles
  enable <iccid>                Enable profile by ICCID
  disable <iccid>               Disable profile by ICCID
  delete <iccid>                Delete profile by ICCID
  nickname <iccid> <nickname>   Set profile nickname
  download                      Download profile (use --code, --imei, --confirmation-code, --confirm)
  discovery                     Discover profiles from SM-DS (use --server, --imei)
  discover-download             Discover and download first available profile (use --server, --imei)
  notifications                 List notifications
  notification-remove <seq>     Remove notification by sequence number
  notification-handle <seq>     Handle notification by sequence number
  auto-notification             Automatically process all pending notifications
  notification-process <seq...> Process specific notifications by sequence number(s)
  configured-addresses          Get configured SM-DP+/SM-DS addresses
  set-default-dp <address>      Set default SM-DP+ address
  challenge                     Get eUICC challenge
  memory-reset                  Reset eUICC memory

Examples:
  # Get EID
  %s eid

  # List profiles
  %s list

  # Download profile
  %s download --code "LPA:1$smdp.io$MATCHING-ID" --confirm

  # Enable profile
  %s enable 8944476500001224158

  # Discover profiles
  %s discovery --imei 356938035643809

All commands output JSON format.
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
