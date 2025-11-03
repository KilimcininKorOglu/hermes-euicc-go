//go:build !linux && !darwin && !windows

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/at"
)

// newATDriver creates a new AT driver instance for FreeBSD and other Unix-like systems
// Uses standard serial port API
func newATDriver(device string) (apdu.SmartCardChannel, error) {
	return at.New(device)
}

// atSupported indicates AT driver is available on FreeBSD and other Unix systems
const atSupported = true

// defaultATDevices returns common AT modem device paths on FreeBSD/Unix
func defaultATDevices() []string {
	return []string{
		"/dev/cuaU0",
		"/dev/cuaU1",
		"/dev/cuaU2",
		"/dev/cuaU3",
		"/dev/ttyU0",
		"/dev/ttyU1",
		"/dev/ttyU2",
		"/dev/ttyU3",
		"/dev/cuau0",
		"/dev/cuau1",
	}
}
