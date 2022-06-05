// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ministo/ministo/pkg/types"
)

const maxPoolBody = 1 << 20

var poolEndpoint = "https://server.duinocoin.com/getPool"

// GetPoolInfo retrieves pool connection information from the Duinocoin server
// by making an HTTP GET request to the pool endpoint.
func GetPoolInfo() (*types.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, poolEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request pool info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pool info status %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxPoolBody+1))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if len(body) > maxPoolBody {
		return nil, fmt.Errorf("pool info response too large")
	}

	var result types.Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}
	if !result.Success {
		return nil, fmt.Errorf("pool info unsuccessful")
	}
	if result.IP == "" || result.Port <= 0 {
		return nil, fmt.Errorf("pool info missing address")
	}

	return &result, nil
}
