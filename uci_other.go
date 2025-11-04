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

// readUCIConfig reads configuration from config file (non-OpenWRT systems)
// Priority: 1) -config flag, 2) ./hermes-euicc.conf, 3) ~/.config/hermes-euicc/config, 4) %APPDATA%\hermes-euicc\config
func readUCIConfig() *UCIConfig {
	defaults := &UCIConfig{
		Driver:  "auto",
		Device:  "",
		Slot:    1,
		Timeout: 30,
	}

	var configPath string

	// Check if -config flag will be provided (we need to parse flags first)
	// This is called before flag.Parse(), so we can't use *configFile yet
	// We'll handle this in main() after flag.Parse()

	// Try to find config file in standard locations
	configPath = findConfigFile()
	if configPath == "" {
		return defaults
	}

	// Read and parse config file
	config, err := readConfigFile(configPath)
	if err != nil {
		return defaults
	}

	return config
}
