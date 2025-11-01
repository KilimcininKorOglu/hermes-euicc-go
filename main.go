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
	"net/url"
	"os"
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
	SequenceNumber int    `json:"sequence_number"`
	ICCID          string `json:"iccid"`
	Operation      int    `json:"operation"`
}

type FailedNotification struct {
	SequenceNumber int    `json:"sequence_number"`
	ICCID          string `json:"iccid"`
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
		"list":                  true,
		"enable":                true,
		"disable":               true,
		"delete":                true,
		"nickname":              true,
		"download":              true,
		"discovery":             true,
		"notifications":         true,
		"notification-remove":   true,
		"notification-handle":   true,
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
	case "notifications":
		handleNotifications(client)
	case "notification-remove":
		handleNotificationRemove(client)
	case "notification-handle":
		handleNotificationHandle(client)
	case "auto-notification":
		handleAutoNotification(client)
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
		server = flag.String("server", "", "SM-DS server URL (default: lpa.ds.gsma.com)")
		imei   = flag.String("imei", "", "IMEI")
	)

	discoveryFlags := flag.NewFlagSet("discovery", flag.ExitOnError)
	discoveryFlags.StringVar(server, "server", "", "SM-DS server")
	discoveryFlags.StringVar(imei, "imei", "", "IMEI")
	discoveryFlags.Parse(flag.Args()[1:])

	var imeiBytes sgp22.IMEI
	if *imei != "" {
		var err error
		imeiBytes, err = sgp22.NewIMEI(*imei)
		if err != nil {
			outputError(fmt.Errorf("invalid IMEI: %w", err))
			os.Exit(1)
		}
	}

	servers := []string{"lpa.ds.gsma.com", "lpa.live.esimdiscovery.com"}
	if *server != "" {
		servers = []string{*server}
	}

	allEntries := make([]DiscoveryResponse, 0, len(servers)*2)
	for _, srv := range servers {
		address := &url.URL{Scheme: "https", Host: srv}
		entries, err := client.Discovery(address, imeiBytes)
		if err != nil {
			if *verbose {
				log.Printf("Discovery failed for %s: %v\n", srv, err)
			}
			continue
		}

		for _, entry := range entries {
			allEntries = append(allEntries, DiscoveryResponse{
				EventID: entry.EventID,
				Address: entry.Address,
			})
		}
	}

	outputSuccess(allEntries)
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
	// Retrieve all pending notifications metadata
	notificationList, err := client.ListNotification()
	if err != nil {
		outputError(err)
		os.Exit(1)
	}

	if len(notificationList) == 0 {
		outputSuccess(AutoNotificationResponse{
			Message:       "no pending notifications",
			Total:         0,
			Processed:     0,
			Failed:        0,
			ProcessedList: []ProcessedNotification{},
			FailedList:    []FailedNotification{},
		})
		return
	}

	// Process all notifications
	processed := make([]ProcessedNotification, 0, len(notificationList))
	failed := make([]FailedNotification, 0)

	for _, metadata := range notificationList {
		// Retrieve the actual notification using sequence number
		notifications, err := client.RetrieveNotificationList(metadata.SequenceNumber)
		if err != nil {
			failed = append(failed, FailedNotification{
				SequenceNumber: int(metadata.SequenceNumber),
				ICCID:          metadata.ICCID.String(),
				Error:          err.Error(),
			})
			continue
		}

		if len(notifications) == 0 {
			failed = append(failed, FailedNotification{
				SequenceNumber: int(metadata.SequenceNumber),
				ICCID:          metadata.ICCID.String(),
				Error:          "notification not found",
			})
			continue
		}

		// Handle the notification
		err = client.HandleNotification(notifications[0])
		if err != nil {
			failed = append(failed, FailedNotification{
				SequenceNumber: int(metadata.SequenceNumber),
				ICCID:          metadata.ICCID.String(),
				Error:          err.Error(),
			})
		} else {
			processed = append(processed, ProcessedNotification{
				SequenceNumber: int(metadata.SequenceNumber),
				ICCID:          metadata.ICCID.String(),
				Operation:      int(metadata.ProfileManagementOperation),
			})
		}
	}

	outputSuccess(AutoNotificationResponse{
		Message:       "auto notification processing completed",
		Total:         len(notificationList),
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
  list                          List all profiles
  enable <iccid>                Enable profile by ICCID
  disable <iccid>               Disable profile by ICCID
  delete <iccid>                Delete profile by ICCID
  nickname <iccid> <nickname>   Set profile nickname
  download                      Download profile (use --code, --imei, --confirmation-code, --confirm)
  discovery                     Discover profiles from SM-DS (use --server, --imei)
  notifications                 List notifications
  notification-remove <seq>     Remove notification by sequence number
  notification-handle <seq>     Handle notification by sequence number
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
