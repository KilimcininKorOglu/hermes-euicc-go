//go:build openwrt

// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

package main

// readConfigFileCustom is a stub for OpenWRT builds
// On OpenWRT, UCI is always used for configuration, config files are not supported
func readConfigFileCustom(configPath string) (*UCIConfig, error) {
	// Return empty config, this will be ignored
	return &UCIConfig{}, nil
}
