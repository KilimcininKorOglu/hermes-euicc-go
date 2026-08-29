//go:build linux && !amd64 && !arm64

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/damonto/euicc-go/driver"
)

// initCCIDDriver returns error on unsupported platforms (MIPS, 32-bit, etc.)
// CCID (PC/SC) requires purego which doesn't support these architectures
func initCCIDDriver() (driver.SmartCardChannel, error) {
	return nil, fmt.Errorf("CCID driver not supported on this platform (use QMI, MBIM, or AT driver)")
}

const ccidSupported = false
