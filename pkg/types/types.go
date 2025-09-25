// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

// Response represents the pool server response for connection details
type Response struct {
	Client  string `json:"client"`
	IP      string `json:"ip"`
	Name    string `json:"name"`
	Port    int    `json:"port"`
	Region  string `json:"region"`
	Server  string `json:"server"`
	Success bool   `json:"success"`
}

// Config holds the mining configuration
type Config struct {
	Username      string
	Difficulty    string
	MiningKey     string
	RigIdentifier string
	MinerBanner   string
	MinerVersion  string
}
