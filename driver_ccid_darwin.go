//go:build darwin

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/ccid"
)

// initCCIDDriver initializes CCID driver using PC/SC framework (macOS)
// macOS has built-in PC/SC support via CryptoTokenKit framework
func initCCIDDriver() (apdu.SmartCardChannel, error) {
	ch, err := ccid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CCID: %w (PC/SC framework may not be available)", err)
	}

	readers, err := ch.ListReaders()
	if err != nil {
		return nil, fmt.Errorf("failed to list readers: %w", err)
	}

	if len(readers) == 0 {
		return nil, fmt.Errorf("no CCID readers found (please connect a USB smart card reader)")
	}

	// Use first available reader
	ch.SetReader(readers[0])
	return ch, nil
}

const ccidSupported = true
