//go:build linux && (amd64 || arm64)

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/damonto/euicc-go/driver"
	"github.com/damonto/euicc-go/driver/ccid"
)

// initCCIDDriver initializes CCID driver using pcscd (Linux)
func initCCIDDriver() (driver.SmartCardChannel, error) {
	ch := ccid.New()

	readers, err := ch.ListReaders()
	if err != nil {
		return nil, fmt.Errorf("failed to list readers: %w", err)
	}

	if len(readers) == 0 {
		return nil, fmt.Errorf("no CCID readers found (please connect a USB smart card reader)")
	}

	ch.SetReader(readers[0])
	return ch, nil
}

const ccidSupported = true
