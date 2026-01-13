// Copyright (c) 2021 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package store

import (
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow/proto/waCompanionReg"
)

// BrowserConfig contains browser and OS information for device registration.
// This makes the WhatsApp connection appear as a genuine browser session.
type BrowserConfig struct {
	// OS name (e.g., "Mac OS", "Windows", "Ubuntu", "Linux")
	OS string
	// Browser name shown in linked devices (e.g., "Chrome", "Firefox", "Safari")
	Browser string
	// OS version (e.g., "14.4.1" for macOS, "10.0.22631" for Windows)
	OSVersion string
	// Platform type determines the browser icon shown in WhatsApp
	PlatformType waCompanionReg.DeviceProps_PlatformType
}

// Predefined browser configurations similar to Baileys
var Browsers = struct {
	// Ubuntu returns a browser config for Ubuntu Linux
	Ubuntu func(browser string) BrowserConfig
	// MacOS returns a browser config for macOS
	MacOS func(browser string) BrowserConfig
	// Windows returns a browser config for Windows
	Windows func(browser string) BrowserConfig
	// Linux returns a browser config for generic Linux
	Linux func(browser string) BrowserConfig
}{
	Ubuntu: func(browser string) BrowserConfig {
		return BrowserConfig{
			OS:           "Ubuntu",
			Browser:      browser,
			OSVersion:    "22.04.4",
			PlatformType: getPlatformType(browser),
		}
	},
	MacOS: func(browser string) BrowserConfig {
		return BrowserConfig{
			OS:           "Mac OS",
			Browser:      browser,
			OSVersion:    "14.4.1",
			PlatformType: getPlatformType(browser),
		}
	},
	Windows: func(browser string) BrowserConfig {
		return BrowserConfig{
			OS:           "Windows",
			Browser:      browser,
			OSVersion:    "10.0.22631",
			PlatformType: getPlatformType(browser),
		}
	},
	Linux: func(browser string) BrowserConfig {
		return BrowserConfig{
			OS:           "Linux",
			Browser:      browser,
			OSVersion:    "6.5.0",
			PlatformType: getPlatformType(browser),
		}
	},
}

// getPlatformType returns the appropriate platform type for a browser name
func getPlatformType(browser string) waCompanionReg.DeviceProps_PlatformType {
	switch browser {
	case "Chrome":
		return waCompanionReg.DeviceProps_CHROME
	case "Firefox":
		return waCompanionReg.DeviceProps_FIREFOX
	case "Safari":
		return waCompanionReg.DeviceProps_SAFARI
	case "Edge":
		return waCompanionReg.DeviceProps_EDGE
	case "Opera":
		return waCompanionReg.DeviceProps_OPERA
	case "IE":
		return waCompanionReg.DeviceProps_IE
	case "Desktop":
		return waCompanionReg.DeviceProps_DESKTOP
	default:
		return waCompanionReg.DeviceProps_CHROME // Default to Chrome
	}
}

// SetBrowserConfig configures the device to appear as a genuine browser.
// This should be called before connecting to WhatsApp.
//
// Example usage:
//
//	store.SetBrowserConfig(store.Browsers.Windows("Chrome"))
//	store.SetBrowserConfig(store.Browsers.MacOS("Safari"))
//	store.SetBrowserConfig(store.Browsers.Ubuntu("Firefox"))
//
// You can also create a custom config:
//
//	store.SetBrowserConfig(store.BrowserConfig{
//	    OS:           "Windows",
//	    Browser:      "Chrome",
//	    OSVersion:    "10.0.22631",
//	    PlatformType: waCompanionReg.DeviceProps_CHROME,
//	})
func SetBrowserConfig(config BrowserConfig) {
	// Update DeviceProps to appear as the specified browser
	DeviceProps.Os = proto.String(config.OS)
	DeviceProps.PlatformType = config.PlatformType.Enum()

	// Parse version from OSVersion string and set it
	version := parseVersion(config.OSVersion)
	DeviceProps.Version = &waCompanionReg.DeviceProps_AppVersion{
		Primary:   proto.Uint32(version[0]),
		Secondary: proto.Uint32(version[1]),
		Tertiary:  proto.Uint32(version[2]),
	}

	// Update BaseClientPayload with OS information
	BaseClientPayload.UserAgent.OsVersion = proto.String(config.OSVersion)
	BaseClientPayload.UserAgent.OsBuildNumber = proto.String(config.OSVersion)
}

// parseVersion parses a version string like "14.4.1" into [3]uint32
func parseVersion(version string) [3]uint32 {
	var v [3]uint32
	var part1, part2, part3 int
	n, _ := parseVersionParts(version, &part1, &part2, &part3)
	if n >= 1 {
		v[0] = uint32(part1)
	}
	if n >= 2 {
		v[1] = uint32(part2)
	}
	if n >= 3 {
		v[2] = uint32(part3)
	}
	return v
}

// parseVersionParts extracts up to 3 numeric parts from a version string
func parseVersionParts(version string, part1, part2, part3 *int) (int, error) {
	// Simple parsing without importing fmt
	parts := splitVersion(version)
	count := 0

	if len(parts) >= 1 {
		*part1 = atoi(parts[0])
		count = 1
	}
	if len(parts) >= 2 {
		*part2 = atoi(parts[1])
		count = 2
	}
	if len(parts) >= 3 {
		*part3 = atoi(parts[2])
		count = 3
	}
	return count, nil
}

// splitVersion splits a version string by dots
func splitVersion(s string) []string {
	var parts []string
	current := ""
	for _, c := range s {
		if c == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// atoi converts a string to int (simplified)
func atoi(s string) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			break // Stop at first non-digit
		}
	}
	return result
}
