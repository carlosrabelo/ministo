// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/carlosrabelo/ministo/pkg/types"
)

const poolEndpoint = "https://server.duinocoin.com/getPool"

// GetPoolInfo retrieves pool connection information from the server
func GetPoolInfo() (*types.Response, error) {
	resp, err := http.Get(poolEndpoint)
	if err != nil {
		return nil, fmt.Errorf("request pool info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var result types.Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return &result, nil
}
