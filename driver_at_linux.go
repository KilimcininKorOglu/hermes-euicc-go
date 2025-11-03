//go:build linux

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

// defaultATDevices returns common AT modem device paths on Linux
func defaultATDevices() []string {
	return []string{
		"/dev/ttyUSB2",
		"/dev/ttyUSB3",
		"/dev/ttyUSB1",
		"/dev/ttyUSB0",
		"/dev/ttyACM0",
		"/dev/ttyACM1",
		"/dev/ttyACM2",
	}
}
