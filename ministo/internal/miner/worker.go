// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"ministo/ministo/pkg/types"
)

const ioTimeout = 30 * time.Second

// Worker represents a mining worker that connects to a pool
type Worker struct {
	conn    net.Conn
	reader  *bufio.Reader
	config  *types.Config
	minerID int
}

// NewWorker creates a new mining worker
func NewWorker(conn net.Conn, config *types.Config) *Worker {
	return &Worker{
		conn:    conn,
		reader:  bufio.NewReader(conn),
		config:  config,
		minerID: randomMinerID(),
	}
}

func randomMinerID() int {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return int(binary.BigEndian.Uint16(b[:])) % 2811
}

func (w *Worker) writeLine(s string) error {
	if err := w.conn.SetWriteDeadline(time.Now().Add(ioTimeout)); err != nil {
		return err
	}
	_, err := w.conn.Write([]byte(s))
	return err
}

func (w *Worker) readLine() (string, error) {
	if err := w.conn.SetReadDeadline(time.Now().Add(ioTimeout)); err != nil {
		return "", err
	}
	line, err := w.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// GetServerVersion reads the server version from the connection
func (w *Worker) GetServerVersion() (string, error) {
	version, err := w.readLine()
	if err != nil {
		return "", fmt.Errorf("read server version: %w", err)
	}

	log.Printf("Server version: %s", version)
	return version, nil
}

// RequestJob requests a new mining job from the server
func (w *Worker) RequestJob() (string, string, int, error) {
	jobRequest := fmt.Sprintf("%s,%s,%s,%s\n",
		"JOB", w.config.Username, w.config.Difficulty, w.config.MiningKey)

	if err := w.writeLine(jobRequest); err != nil {
		return "", "", 0, fmt.Errorf("send job request: %w", err)
	}

	response, err := w.readLine()
	if err != nil {
		return "", "", 0, fmt.Errorf("read job response: %w", err)
	}

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
	// Format matches Official ESP miners:
	// nonce,hashrate,banner ver,rig,DUCOID…,walletid
	submission := fmt.Sprintf("%d,%d,%s %s,%s,%s,%d\n",
		result.Result, result.Hashrate,
		w.config.MinerBanner, w.config.MinerVersion,
		w.config.RigIdentifier, w.config.DucoID, w.minerID)

	if err := w.writeLine(submission); err != nil {
		return "", fmt.Errorf("send result: %w", err)
	}

	feedback, err := w.readLine()
	if err != nil {
		return "", fmt.Errorf("read feedback: %w", err)
	}

	return feedback, nil
}

// Close closes the worker connection
func (w *Worker) Close() error {
	return w.conn.Close()
}
