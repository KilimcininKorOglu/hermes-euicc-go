//go:build !linux || !(amd64 || arm64)

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
)

func initCCIDDriver() (apdu.SmartCardChannel, error) {
	return nil, fmt.Errorf("CCID driver not supported on this platform (requires amd64/arm64 + linux)")
}

const ccidSupported = false
