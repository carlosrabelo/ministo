// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"log"
	"net"
	"net/url"

	"golang.org/x/net/proxy"
)

// NewProxyDialer creates a new SOCKS5 proxy dialer
func NewProxyDialer(proxyURL string) (proxy.Dialer, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("parse proxy URL: %w", err)
	}

	d, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("create proxy dialer: %w", err)
	}

	return d, nil
}

// ConnectToServer establishes a connection through the proxy to the mining server
func ConnectToServer(dialer proxy.Dialer, serverIP string, port int) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", serverIP, port)

	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial server: %w", err)
	}

	log.Printf("Connected to mining server: %s", address)
	return conn, nil
}
