// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"testing"
)

func TestNewProxyDialer(t *testing.T) {
	proxyURL := "socks5://localhost:9050"

	dialer, err := NewProxyDialer(proxyURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if dialer == nil {
		t.Error("Expected dialer to be non-nil")
	}
}

func TestNewProxyDialerInvalidURL(t *testing.T) {
	proxyURL := "invalid-url"

	_, err := NewProxyDialer(proxyURL)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestNewProxyDialerEmptyURL(t *testing.T) {
	proxyURL := ""

	_, err := NewProxyDialer(proxyURL)
	if err == nil {
		t.Error("Expected error for empty URL, got nil")
	}
}
