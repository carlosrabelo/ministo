// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"strings"
	"testing"
)

func TestParseMAC(t *testing.T) {
	tests := []struct {
		in      string
		want    [6]byte
		wantErr bool
	}{
		{in: "24:0A:C4:12:34:56", want: [6]byte{0x24, 0x0A, 0xC4, 0x12, 0x34, 0x56}},
		{in: "240ac4123456", want: [6]byte{0x24, 0x0A, 0xC4, 0x12, 0x34, 0x56}},
		{in: "24-0a-c4-12-34-56", want: [6]byte{0x24, 0x0A, 0xC4, 0x12, 0x34, 0x56}},
		{in: "short", wantErr: true},
		{in: "GG:GG:GG:GG:GG:GG", wantErr: true},
	}
	for _, tt := range tests {
		got, err := ParseMAC(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("ParseMAC(%q): expected error", tt.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseMAC(%q): %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("ParseMAC(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestESP32DucoID(t *testing.T) {
	mac := [6]byte{0x24, 0x0A, 0xC4, 0x12, 0x34, 0x56}
	id := ESP32DucoID(mac)
	if !strings.HasPrefix(id, "DUCOID") {
		t.Fatalf("missing DUCOID prefix: %s", id)
	}
	if len(id) != len("DUCOID")+12 {
		t.Fatalf("unexpected length: %s", id)
	}
	if id != strings.ToUpper(id) {
		t.Fatalf("ESP32 DUCOID should be uppercase: %s", id)
	}
}

func TestESP8266DucoID(t *testing.T) {
	mac := [6]byte{0x24, 0x0A, 0xC4, 0xAB, 0xCD, 0xEF}
	id := ESP8266DucoID(mac)
	if id != "DUCOIDabcdef" {
		t.Fatalf("got %s, want DUCOIDabcdef", id)
	}
}

func TestRandomEspressifMAC(t *testing.T) {
	mac, err := RandomEspressifMAC()
	if err != nil {
		t.Fatal(err)
	}
	if mac[0] != 0x24 || mac[1] != 0x0A || mac[2] != 0xC4 {
		t.Fatalf("unexpected OUI: %s", FormatMAC(mac))
	}
}
