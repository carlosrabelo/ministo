// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// Espressif OUI used for generated MACs (common on ESP modules).
var espressifOUI = [3]byte{0x24, 0x0A, 0xC4}

// RandomEspressifMAC returns a random 48-bit MAC with an Espressif OUI.
func RandomEspressifMAC() ([6]byte, error) {
	var mac [6]byte
	copy(mac[:3], espressifOUI[:])
	if _, err := rand.Read(mac[3:]); err != nil {
		return mac, fmt.Errorf("generate mac: %w", err)
	}
	return mac, nil
}

// ParseMAC parses AA:BB:CC:DD:EE:FF or AABBCCDDEEFF into 6 bytes.
func ParseMAC(s string) ([6]byte, error) {
	var mac [6]byte
	cleaned := strings.NewReplacer(":", "", "-", "", ".", "").Replace(strings.TrimSpace(s))
	if len(cleaned) != 12 {
		return mac, fmt.Errorf("invalid mac %q: need 12 hex digits", s)
	}
	b, err := hex.DecodeString(cleaned)
	if err != nil {
		return mac, fmt.Errorf("invalid mac %q: %w", s, err)
	}
	copy(mac[:], b)
	return mac, nil
}

// FormatMAC returns AA:BB:CC:DD:EE:FF (uppercase).
func FormatMAC(mac [6]byte) string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

// ESP32DucoID builds DUCOID like Official ESP32 Miner from a 48-bit MAC
// (ESP.getEfuseMac → "%04X%08X").
func ESP32DucoID(mac [6]byte) string {
	chipid := uint64(mac[0]) |
		uint64(mac[1])<<8 |
		uint64(mac[2])<<16 |
		uint64(mac[3])<<24 |
		uint64(mac[4])<<32 |
		uint64(mac[5])<<40
	chip := uint16(chipid >> 32)
	low := uint32(chipid)
	return fmt.Sprintf("DUCOID%04X%08X", chip, low)
}

// ESP8266DucoID builds DUCOID like Official ESP8266 Miner from chip id
// (last 3 MAC octets as lowercase hex, Arduino String(..., HEX)).
func ESP8266DucoID(mac [6]byte) string {
	chipID := uint32(mac[3])<<16 | uint32(mac[4])<<8 | uint32(mac[5])
	return fmt.Sprintf("DUCOID%x", chipID)
}

// ESP32ChipHex returns the 12-char chip id used in Auto rig names.
func ESP32ChipHex(mac [6]byte) string {
	return strings.TrimPrefix(ESP32DucoID(mac), "DUCOID")
}

// ESP8266ChipHex returns the chip id hex used in Auto rig names.
func ESP8266ChipHex(mac [6]byte) string {
	return strings.TrimPrefix(ESP8266DucoID(mac), "DUCOID")
}
