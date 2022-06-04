// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"testing"
)

func TestResponseSerialization(t *testing.T) {
	response := Response{
		Client:  "test_client",
		IP:      "127.0.0.1",
		Name:    "test_pool",
		Port:    2812,
		Region:  "test_region",
		Server:  "test_server",
		Success: true,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.IP != response.IP {
		t.Errorf("Expected IP %s, got %s", response.IP, unmarshaled.IP)
	}

	if unmarshaled.Port != response.Port {
		t.Errorf("Expected port %d, got %d", response.Port, unmarshaled.Port)
	}

	if unmarshaled.Success != response.Success {
		t.Errorf("Expected success %t, got %t", response.Success, unmarshaled.Success)
	}
}

func TestConfigDefaults(t *testing.T) {
	config := Config{
		Username:      "testuser",
		Difficulty:    "MEDIUM",
		MiningKey:     "testkey",
		RigIdentifier: "testrig",
		MinerBanner:   "TestMiner",
		MinerVersion:  "1.0",
	}

	if config.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", config.Username)
	}

	if config.MinerVersion != "1.0" {
		t.Errorf("Expected version 1.0, got %s", config.MinerVersion)
	}
}
