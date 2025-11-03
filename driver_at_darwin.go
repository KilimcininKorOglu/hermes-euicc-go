//go:build darwin

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/at"
)

// newATDriver creates a new AT driver instance for macOS
// Uses serial port API for /dev/cu.usbserial* devices
func newATDriver(device string) (apdu.SmartCardChannel, error) {
	return at.New(device)
}

// atSupported indicates AT driver is available on macOS
const atSupported = true

// defaultATDevices returns common AT modem device paths on macOS
func defaultATDevices() []string {
	return []string{
		"/dev/cu.usbserial",
		"/dev/cu.usbserial-0",
		"/dev/cu.usbserial-1",
		"/dev/cu.usbserial-2",
		"/dev/cu.usbmodem1",
		"/dev/cu.usbmodem2",
		"/dev/cu.usbmodem3",
	}
}
