//go:build linux && (amd64 || arm64)

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/ccid"
)

// initCCIDDriver initializes CCID driver using pcscd (Linux)
func initCCIDDriver() (apdu.SmartCardChannel, error) {
	ch, err := ccid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CCID: %w (is pcscd running?)", err)
	}

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
