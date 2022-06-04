// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package miner

import (
	"context"
	"testing"
)

func TestFindHash(t *testing.T) {
	baseHash := "test"
	targetHash := "b444ac06613fc8d63795be9ad0beaf55011936ac" // SHA1 of "test1"
	difficulty := 5

	result := FindHash(baseHash, targetHash, difficulty)

	if !result.Found {
		t.Error("Expected to find hash but didn't")
	}

	if result.Result != 1 {
		t.Errorf("Expected result 1, got %d", result.Result)
	}
}

func TestFindHashWithContext(t *testing.T) {
	baseHash := "test"
	targetHash := "b444ac06613fc8d63795be9ad0beaf55011936ac" // SHA1 of "test1"
	difficulty := 5

	result := FindHashWithContext(context.Background(), baseHash, targetHash, difficulty)

	if !result.Found {
		t.Error("Expected to find hash but didn't")
	}

	if result.Result != 1 {
		t.Errorf("Expected result 1, got %d", result.Result)
	}
}

func TestFindHashWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := FindHashWithContext(ctx, "test", "0000000000000000000000000000000000000000", 100000)
	if result.Found {
		t.Error("Expected hash not to be found due to context cancellation")
	}
}

func TestFindHashNotFound(t *testing.T) {
	result := FindHash("test", "0000000000000000000000000000000000000000", 1)
	if result.Found {
		t.Error("Expected hash not to be found")
	}
}

func TestFindHashInvalidTarget(t *testing.T) {
	result := FindHash("test", "not-hex", 5)
	if result.Found {
		t.Error("Expected invalid target not to match")
	}
}
