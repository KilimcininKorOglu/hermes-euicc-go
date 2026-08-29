//go:build windows

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/damonto/euicc-go/driver"
	"github.com/damonto/euicc-go/driver/ccid"
)

// initCCIDDriver initializes CCID driver using winscard.dll (Windows)
// Windows has built-in smart card support via winscard.dll
func initCCIDDriver() (driver.SmartCardChannel, error) {
	ch := ccid.New()

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
