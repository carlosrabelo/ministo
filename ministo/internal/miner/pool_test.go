// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func withPoolEndpoint(t *testing.T, url string) {
	t.Helper()
	orig := poolEndpoint
	poolEndpoint = url
	t.Cleanup(func() { poolEndpoint = orig })
}

func TestGetPoolInfo(t *testing.T) {
	mockResponse := `{
		"client": "test_client",
		"ip": "127.0.0.1",
		"name": "test_pool",
		"port": 2812,
		"region": "test_region",
		"server": "test_server",
		"success": true
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()
	withPoolEndpoint(t, server.URL)

	result, err := GetPoolInfo()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.IP != "127.0.0.1" {
		t.Errorf("Expected IP 127.0.0.1, got %s", result.IP)
	}

	if result.Port != 2812 {
		t.Errorf("Expected port 2812, got %d", result.Port)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}
}

func TestGetPoolInfoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()
	withPoolEndpoint(t, server.URL)

	_, err := GetPoolInfo()
	if err == nil {
		t.Fatal("Expected error for unreachable pool, got nil")
	}
}

func TestGetPoolInfoHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()
	withPoolEndpoint(t, server.URL)

	_, err := GetPoolInfo()
	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}
}

func TestGetPoolInfoUnsuccessful(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": false, "ip": "127.0.0.1", "port": 2812}`))
	}))
	defer server.Close()
	withPoolEndpoint(t, server.URL)

	_, err := GetPoolInfo()
	if err == nil {
		t.Fatal("Expected error for unsuccessful pool response, got nil")
	}
}
