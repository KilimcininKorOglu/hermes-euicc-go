//go:build openwrt

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"os/exec"
	"strconv"
	"strings"
)

// UCI Configuration structure
type UCIConfig struct {
	Driver  string
	Device  string
	Slot    int
	Timeout int
}

// readUCIConfig reads configuration from OpenWRT UCI system
func readUCIConfig() *UCIConfig {
	config := &UCIConfig{
		Driver:  "auto",
		Device:  "",
		Slot:    1,
		Timeout: 30,
	}

	// Check if uci command exists (OpenWRT only)
	if _, err := exec.LookPath("uci"); err != nil {
		return config // Return defaults if UCI not available
	}

	// Read driver setting
	if out, err := exec.Command("uci", "get", "hermes-euicc.hermes-euicc.driver").Output(); err == nil {
		driver := strings.TrimSpace(string(out))
		if driver != "" && driver != "auto" {
			config.Driver = driver
		}
	}

	// Read device setting
	if out, err := exec.Command("uci", "get", "hermes-euicc.hermes-euicc.device").Output(); err == nil {
		device := strings.TrimSpace(string(out))
		if device != "" {
			config.Device = device
		}
	}

	// Read slot setting
	if out, err := exec.Command("uci", "get", "hermes-euicc.hermes-euicc.slot").Output(); err == nil {
		if slot, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && slot > 0 {
			config.Slot = slot
		}
	}

	// Read timeout setting
	if out, err := exec.Command("uci", "get", "hermes-euicc.hermes-euicc.timeout").Output(); err == nil {
		if timeout, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && timeout > 0 {
			config.Timeout = timeout
		}
	}

	return config
}
