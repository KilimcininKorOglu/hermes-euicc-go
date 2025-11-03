//go:build linux

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/KilimcininKorOglu/euicc-go/apdu"
	"github.com/KilimcininKorOglu/euicc-go/driver/at"
	"github.com/KilimcininKorOglu/euicc-go/driver/mbim"
	"github.com/KilimcininKorOglu/euicc-go/driver/qmi"
)

// newQMIDriver creates a new QMI driver instance (Linux only)
func newQMIDriver(device string, slot uint8) (apdu.SmartCardChannel, error) {
	return qmi.New(device, slot)
}

// newMBIMDriver creates a new MBIM driver instance (Linux only)
func newMBIMDriver(device string, slot uint8) (apdu.SmartCardChannel, error) {
	return mbim.New(device, slot)
}

// newATDriver creates a new AT driver instance (Linux only)
func newATDriver(device string) (apdu.SmartCardChannel, error) {
	return at.New(device)
}

// Driver availability flags
const (
	qmiSupported  = true
	mbimSupported = true
	atSupported   = true
)
