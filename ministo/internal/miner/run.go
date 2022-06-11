// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"context"
	"fmt"
	"log"
	"strings"

	"ministo/ministo/internal/proxy"
	"ministo/ministo/pkg/types"
)

// Run connects to the pool through the given SOCKS5 proxy and mines until an error occurs.
func Run(cfg *types.Config, proxyURL string) error {
	return RunContext(context.Background(), cfg, proxyURL)
}

// RunContext is like Run but stops when ctx is cancelled.
func RunContext(ctx context.Context, cfg *types.Config, proxyURL string) error {
	poolInfo, err := GetPoolInfo()
	if err != nil {
		return fmt.Errorf("get pool info: %w", err)
	}

	proxyDialer, err := proxy.NewProxyDialer(proxyURL)
	if err != nil {
		return fmt.Errorf("create proxy dialer: %w", err)
	}

	conn, err := proxy.ConnectToServerWithContext(ctx, proxyDialer, poolInfo.IP, poolInfo.Port)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}
	defer conn.Close()

	worker := NewWorker(conn, cfg)

	if _, err := worker.GetServerVersion(); err != nil {
		return fmt.Errorf("get server version: %w", err)
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		hash, target, difficulty, err := worker.RequestJob()
		if err != nil {
			return fmt.Errorf("request job: %w", err)
		}

		result := FindHashWithContext(ctx, hash, target, difficulty)
		if !result.Found {
			if err := ctx.Err(); err != nil {
				return err
			}
			log.Printf("No solution found for difficulty %d", difficulty)
			continue
		}

		feedback, err := worker.SubmitResult(result)
		if err != nil {
			return fmt.Errorf("submit result: %w", err)
		}

		khashrate := result.Hashrate / 1000
		switch {
		case strings.HasPrefix(feedback, "GOOD"), strings.HasPrefix(feedback, "BLOCK"):
			log.Printf("Accepted share %d Hashrate %d kH/s Difficulty %d",
				result.Result, khashrate, difficulty)
		case strings.HasPrefix(feedback, "BAD"):
			log.Printf("Rejected share %d Hashrate %d kH/s Difficulty %d",
				result.Result, khashrate, difficulty)
		default:
			return fmt.Errorf("unknown feedback: %s", feedback)
		}
	}
}
