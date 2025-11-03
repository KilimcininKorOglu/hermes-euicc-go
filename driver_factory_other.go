//go:build !linux

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/KilimcininKorOglu/euicc-go/apdu"
)

// newQMIDriver returns an error on non-Linux platforms (QMI is Linux-only)
func newQMIDriver(device string, slot uint8) (apdu.SmartCardChannel, error) {
	return nil, fmt.Errorf("QMI driver not supported on this platform (Linux only)")
}

// newMBIMDriver returns an error on non-Linux platforms (MBIM is Linux-only)
func newMBIMDriver(device string, slot uint8) (apdu.SmartCardChannel, error) {
	return nil, fmt.Errorf("MBIM driver not supported on this platform (Linux only)")
}

// newATDriver is implemented in platform-specific files (driver_at_*.go)
// No stub needed here - build tags ensure correct file is compiled

// Driver availability flags
const (
	qmiSupported  = false
	mbimSupported = false
	// atSupported is defined in platform-specific driver_at_*.go files
)
