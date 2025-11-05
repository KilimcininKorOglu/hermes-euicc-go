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
	if out, err := exec.Command("uci", "get", "hermes_euicc.config.driver").Output(); err == nil {
		driver := strings.TrimSpace(string(out))
		if driver != "" && driver != "auto" {
			// Accept 'uqmi' as alias for 'qmi' driver
			if driver == "uqmi" {
				driver = "qmi"
			}
			config.Driver = driver
		}
	}

	// Read device setting - try general device first, then driver-specific
	if out, err := exec.Command("uci", "get", "hermes_euicc.config.qmi_device").Output(); err == nil {
		device := strings.TrimSpace(string(out))
		if device != "" {
			config.Device = device
		}
	}
	// Fallback to mbim_device if qmi_device not found
	if config.Device == "" {
		if out, err := exec.Command("uci", "get", "hermes_euicc.config.mbim_device").Output(); err == nil {
			device := strings.TrimSpace(string(out))
			if device != "" {
				config.Device = device
			}
		}
	}
	// Fallback to at_device if still not found
	if config.Device == "" {
		if out, err := exec.Command("uci", "get", "hermes_euicc.config.at_device").Output(); err == nil {
			device := strings.TrimSpace(string(out))
			if device != "" {
				config.Device = device
			}
		}
	}

	// Read slot setting - try qmi_sim_slot first, then general slot
	if out, err := exec.Command("uci", "get", "hermes_euicc.config.qmi_sim_slot").Output(); err == nil {
		if slot, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && slot > 0 {
			config.Slot = slot
		}
	}
	// Fallback to general slot if qmi_sim_slot not found
	if config.Slot == 1 {
		if out, err := exec.Command("uci", "get", "hermes_euicc.config.slot").Output(); err == nil {
			if slot, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && slot > 0 {
				config.Slot = slot
			}
		}
	}

	// Read timeout setting
	if out, err := exec.Command("uci", "get", "hermes_euicc.config.timeout").Output(); err == nil {
		if timeout, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && timeout > 0 {
			config.Timeout = timeout
		}
	}

	return config
}
