//go:build windows

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/at"
)

// newATDriver creates a new AT driver instance for Windows
// Uses Win32 COM port API
func newATDriver(device string) (apdu.SmartCardChannel, error) {
	return at.New(device)
}

// atSupported indicates AT driver is available on Windows
const atSupported = true

// defaultATDevices returns common AT modem COM port paths on Windows
func defaultATDevices() []string {
	return []string{
		"COM1",
		"COM2",
		"COM3",
		"COM4",
		"COM5",
		"COM6",
		"COM7",
		"COM8",
		"COM9",
		"COM10",
	}
}
