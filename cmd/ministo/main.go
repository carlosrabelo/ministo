// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/carlosrabelo/ministo/internal/miner"
	"github.com/carlosrabelo/ministo/internal/proxy"
	"github.com/carlosrabelo/ministo/pkg/types"
)

var defaultConfig = &types.Config{
	Username:      "solracolebar",
	Difficulty:    "",
	MiningKey:     "None",
	RigIdentifier: "pc",
	MinerBanner:   "Ministo",
	MinerVersion:  "0.1",
}

func main() {
	// Get pool information
	poolInfo, err := miner.GetPoolInfo()
	if err != nil {
		log.Fatalf("Failed to get pool info: %v", err)
	}

	// Create proxy dialer
	proxyDialer, err := proxy.NewProxyDialer("socks5://localhost:9050")
	if err != nil {
		log.Fatalf("Failed to create proxy dialer: %v", err)
	}

	// Connect to mining server
	conn, err := proxy.ConnectToServer(proxyDialer, poolInfo.IP, poolInfo.Port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create worker
	worker := miner.NewWorker(conn, defaultConfig)

	// Get server version
	_, err = worker.GetServerVersion()
	if err != nil {
		log.Fatalf("Failed to get server version: %v", err)
	}

	// Mining loop
	for {
		// Request job
		hash, target, difficulty, err := worker.RequestJob()
		if err != nil {
			log.Fatalf("Failed to request job: %v", err)
		}

		// Perform mining
		result := miner.FindHash(hash, target, difficulty)
		if !result.Found {
			log.Printf("No solution found for difficulty %d", difficulty)
			continue
		}

		// Submit result
		feedback, err := worker.SubmitResult(result)
		if err != nil {
			log.Fatalf("Failed to submit result: %v", err)
		}

		// Process feedback
		khashrate := result.Hashrate / 1000
		switch feedback {
		case "GOOD":
			log.Printf("Accepted share %d Hashrate %d kH/s Difficulty %d",
				result.Result, khashrate, difficulty)
		case "BAD":
			log.Printf("Rejected share %d Hashrate %d kH/s Difficulty %d",
				result.Result, khashrate, difficulty)
		default:
			log.Fatalf("Unknown feedback: %s", feedback)
		}
	}
}
