// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// NewProxyDialer creates a new SOCKS5 proxy dialer from the provided URL.
// It parses the proxy URL and returns a dialer that can be used for connections.
func NewProxyDialer(proxyURL string) (proxy.Dialer, error) {
	if strings.TrimSpace(proxyURL) == "" {
		return nil, fmt.Errorf("proxy URL is empty")
	}

	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("parse proxy URL: %w", err)
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("proxy URL missing scheme")
	}

	d, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("create proxy dialer: %w", err)
	}

	return d, nil
}

// ConnectToServer establishes a connection through the proxy to the mining server
// using the provided dialer and server address information.
func ConnectToServer(dialer proxy.Dialer, serverIP string, port int) (net.Conn, error) {
	return ConnectToServerWithContext(context.Background(), dialer, serverIP, port)
}

// ConnectToServerWithContext establishes a connection with context support for timeout and cancellation.
func ConnectToServerWithContext(ctx context.Context, dialer proxy.Dialer, serverIP string, port int) (net.Conn, error) {
	address := net.JoinHostPort(serverIP, strconv.Itoa(port))

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if contextDialer, ok := dialer.(proxy.ContextDialer); ok {
		conn, err := contextDialer.DialContext(ctx, "tcp", address)
		if err != nil {
			return nil, fmt.Errorf("dial server: %w", err)
		}
		log.Printf("Connected to mining server: %s", address)
		return conn, nil
	}

	type dialResult struct {
		conn net.Conn
		err  error
	}
	ch := make(chan dialResult, 1)
	go func() {
		conn, err := dialer.Dial("tcp", address)
		ch <- dialResult{conn, err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("dial server: %w", ctx.Err())
	case r := <-ch:
		if r.err != nil {
			return nil, fmt.Errorf("dial server: %w", r.err)
		}
		log.Printf("Connected to mining server: %s", address)
		return r.conn, nil
	}
}
