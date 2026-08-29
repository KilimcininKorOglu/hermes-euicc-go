//go:build linux

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/damonto/euicc-go/driver"
	"github.com/damonto/euicc-go/driver/at"
	"github.com/damonto/euicc-go/driver/mbim"
	"github.com/damonto/euicc-go/driver/qcom"
)

// newQMIDriver creates a new QMI driver instance (Linux only)
func newQMIDriver(device string, slot uint8) (driver.SmartCardChannel, error) {
	return qcom.NewQMI(qcom.WithAutoDetect(device), qcom.WithSlot(slot))
}

// newMBIMDriver creates a new MBIM driver instance (Linux only)
func newMBIMDriver(device string, slot uint8) (driver.SmartCardChannel, error) {
	return mbim.New(mbim.WithDirect(device), mbim.WithSlot(slot))
}

// newATDriver creates a new AT driver instance (Linux only)
func newATDriver(device string) (driver.SmartCardChannel, error) {
	return at.New(device)
}

// Driver availability flags
const (
	qmiSupported  = true
	mbimSupported = true
	atSupported   = true
)
