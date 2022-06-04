// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
)

// HashResult contains the result of a mining operation
type HashResult struct {
	Result   int
	Hashrate int
	Duration time.Duration
	Found    bool
}

// FindHash performs SHA1 hashing to find the target hash by iterating through nonces
// until a match is found or the maximum attempts based on difficulty are reached.
func FindHash(baseHash, targetHash string, difficulty int) HashResult {
	return FindHashWithContext(context.Background(), baseHash, targetHash, difficulty)
}

// FindHashWithContext performs SHA1 hashing with context support for cancellation.
// Duino-Coin jobs are SHA-1(baseHash + decimalNonce), matching the official miners.
func FindHashWithContext(ctx context.Context, baseHash, targetHash string, difficulty int) HashResult {
	if difficulty < 0 {
		return HashResult{Found: false}
	}

	target, err := hex.DecodeString(targetHash)
	if err != nil || len(target) != sha1.Size {
		return HashResult{Found: false}
	}

	start := time.Now()
	prefix := []byte(baseHash)
	h := sha1.New()
	sum := make([]byte, 0, sha1.Size)
	nonceBuf := make([]byte, 0, 20)
	maxNonce := difficulty * 100

	for nonce := 0; nonce <= maxNonce; nonce++ {
		if nonce%64 == 0 {
			select {
			case <-ctx.Done():
				return HashResult{Found: false}
			default:
			}
		}

		h.Reset()
		h.Write(prefix)
		nonceBuf = strconv.AppendInt(nonceBuf[:0], int64(nonce), 10)
		h.Write(nonceBuf)
		sum = h.Sum(sum[:0])

		if bytes.Equal(sum, target) {
			duration := time.Since(start)
			hashrate := 0
			if secs := duration.Seconds(); secs > 0 {
				hashrate = int(float64(nonce) / secs)
			}
			return HashResult{
				Result:   nonce,
				Hashrate: hashrate,
				Duration: duration,
				Found:    true,
			}
		}
	}

	return HashResult{Found: false}
}
