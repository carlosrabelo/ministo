// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
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

// FindHash performs SHA1 hashing to find the target hash
func FindHash(baseHash, targetHash string, difficulty int) HashResult {
	start := time.Now()

	for nonce := 0; nonce <= difficulty*100; nonce++ {
		h := sha1.New()
		h.Write([]byte(baseHash + strconv.Itoa(nonce)))
		hash := hex.EncodeToString(h.Sum(nil))

		if hash == targetHash {
			duration := time.Since(start)
			hashrate := int(float64(nonce) / duration.Seconds())

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
