//go:build linux && (amd64 || arm64)

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/ccid"
)

func initCCIDDriver() (apdu.SmartCardChannel, error) {
	ch, err := ccid.New()
	if err != nil {
		return nil, err
	}

	readers, err := ch.ListReaders()
	if err != nil {
		return nil, err
	}

	if len(readers) == 0 {
		return nil, fmt.Errorf("no CCID readers found")
	}

	ch.SetReader(readers[0])
	return ch, nil
}

const ccidSupported = true
