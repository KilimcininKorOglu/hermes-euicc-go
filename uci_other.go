//go:build !openwrt

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

// UCI Configuration structure
type UCIConfig struct {
	Driver  string
	Device  string
	Slot    int
	Timeout int
}

// readUCIConfig returns default configuration (UCI not available on non-OpenWRT systems)
func readUCIConfig() *UCIConfig {
	return &UCIConfig{
		Driver:  "auto",
		Device:  "",
		Slot:    1,
		Timeout: 30,
	}
}
