// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

// Local reimplementation of the fork's SM-DS discovery convenience API on the
// upstream damonto/euicc-go lpa.Client (which exposes only low-level Discovery).

package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/damonto/euicc-go/lpa"
)

// defaultSMDSAddress is the default GSMA SM-DS server address.
const defaultSMDSAddress = "lpa.ds.gsma.com"

// discoveredProfile is a profile discovered from an SM-DS server.
type discoveredProfile struct {
	EventID     string
	SMDPAddress string
}

// discoverProfiles queries an SM-DS server for pending profile downloads.
func discoverProfiles(client *lpa.Client, smds *string, imei []byte) ([]*discoveredProfile, error) {
	address := defaultSMDSAddress
	if smds != nil && *smds != "" {
		address = *smds
	}
	smdsURL := &url.URL{Scheme: "https", Host: address}

	entries, err := client.Discovery(smdsURL, imei)
	if err != nil {
		return nil, fmt.Errorf("SM-DS discovery failed: %w", err)
	}

	profiles := make([]*discoveredProfile, len(entries))
	for i, entry := range entries {
		profiles[i] = &discoveredProfile{EventID: entry.EventID, SMDPAddress: entry.Address}
	}
	return profiles, nil
}

// discoverAndDownload discovers profiles then downloads the first one found.
// Returns (false, nil) when no profiles are available.
func discoverAndDownload(client *lpa.Client, smds *string, imei []byte) (bool, error) {
	profiles, err := discoverProfiles(client, smds, imei)
	if err != nil {
		return false, err
	}
	if len(profiles) == 0 {
		return false, nil
	}

	first := profiles[0]
	ac := &lpa.ActivationCode{
		SMDP: &url.URL{Scheme: "https", Host: first.SMDPAddress},
		IMEI: string(imei),
	}
	if _, err := client.DownloadProfile(context.Background(), ac, nil); err != nil {
		return false, fmt.Errorf("download from %s failed: %w", first.SMDPAddress, err)
	}
	return true, nil
}
