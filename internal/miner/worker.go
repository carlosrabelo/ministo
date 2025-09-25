// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/carlosrabelo/ministo/pkg/types"
)

// Worker represents a mining worker that connects to a pool
type Worker struct {
	conn    net.Conn
	config  *types.Config
	minerID int
}

// NewWorker creates a new mining worker
func NewWorker(conn net.Conn, config *types.Config) *Worker {
	return &Worker{
		conn:    conn,
		config:  config,
		minerID: rand.Intn(2811),
	}
}

// GetServerVersion reads the server version from the connection
func (w *Worker) GetServerVersion() (string, error) {
	reader := bufio.NewReader(w.conn)
	version, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read server version: %w", err)
	}

	version = strings.TrimSpace(version)
	log.Printf("Server version: %s", version)
	return version, nil
}

// RequestJob requests a new mining job from the server
func (w *Worker) RequestJob() (string, string, int, error) {
	jobRequest := fmt.Sprintf("%s,%s,%s,%s\n",
		"JOB", w.config.Username, w.config.Difficulty, w.config.MiningKey)

	_, err := w.conn.Write([]byte(jobRequest))
	if err != nil {
		return "", "", 0, fmt.Errorf("send job request: %w", err)
	}

	reader := bufio.NewReader(w.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", "", 0, fmt.Errorf("read job response: %w", err)
	}

	response = strings.TrimSpace(response)
	parts := strings.Split(response, ",")
	if len(parts) < 3 {
		return "", "", 0, fmt.Errorf("invalid job response: %s", response)
	}

	hash := parts[0]
	target := parts[1]
	difficulty, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", 0, fmt.Errorf("parse difficulty: %w", err)
	}

	return hash, target, difficulty, nil
}

// SubmitResult submits a mining result to the server
func (w *Worker) SubmitResult(result HashResult) (string, error) {
	submission := fmt.Sprintf("%d,%d,%s %s,%s,,%d\n",
		result.Result, result.Hashrate,
		w.config.MinerBanner, w.config.MinerVersion,
		w.config.RigIdentifier, w.minerID)

	_, err := w.conn.Write([]byte(submission))
	if err != nil {
		return "", fmt.Errorf("send result: %w", err)
	}

	reader := bufio.NewReader(w.conn)
	feedback, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read feedback: %w", err)
	}

	feedback = strings.TrimSpace(feedback)
	return feedback, nil
}

// Close closes the worker connection
func (w *Worker) Close() error {
	return w.conn.Close()
}
