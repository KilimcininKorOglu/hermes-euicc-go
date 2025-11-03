//go:build !linux && !darwin && !windows

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/ccid"
)

// initCCIDDriver attempts CCID initialization on other Unix-like platforms (FreeBSD, etc)
// Uses pcsc-lite if available
func initCCIDDriver() (apdu.SmartCardChannel, error) {
	ch, err := ccid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CCID: %w (install pcsc-lite package)", err)
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
