//go:build !openwrt

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readConfigFile reads configuration from a key=value format config file
func readConfigFile(configPath string) (*UCIConfig, error) {
	config := &UCIConfig{
		Driver:  "auto",
		Device:  "",
		Slot:    1,
		Timeout: 30,
	}

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		// Set config values
		switch key {
		case "driver":
			if value != "" {
				config.Driver = value
			}
		case "device":
			if value != "" {
				config.Device = value
			}
		case "slot":
			if slot, err := strconv.Atoi(value); err == nil && slot > 0 {
				config.Slot = slot
			}
		case "timeout":
			if timeout, err := strconv.Atoi(value); err == nil && timeout > 0 {
				config.Timeout = timeout
			}
		}
	}

	return config, scanner.Err()
}

// readConfigFileCustom is a wrapper for readConfigFile that can be called from main
// It's the same as readConfigFile but named differently to avoid confusion
func readConfigFileCustom(configPath string) (*UCIConfig, error) {
	return readConfigFile(configPath)
}

// findConfigFile searches for config file in standard locations
func findConfigFile() string {
	// Priority order:
	// 1. ./hermes-euicc.conf (current directory)
	// 2. $HOME/.config/hermes-euicc/config (Linux/macOS/FreeBSD)
	// 3. %APPDATA%\hermes-euicc\config (Windows)

	// Check current directory
	if _, err := os.Stat("hermes-euicc.conf"); err == nil {
		return "hermes-euicc.conf"
	}

	// Check user config directory
	if home, err := os.UserHomeDir(); err == nil {
		// Unix-like systems (Linux, macOS, FreeBSD)
		configPath := filepath.Join(home, ".config", "hermes-euicc", "config")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// Windows
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			configPath = filepath.Join(appdata, "hermes-euicc", "config")
			if _, err := os.Stat(configPath); err == nil {
				return configPath
			}
		}
	}

	return ""
}
