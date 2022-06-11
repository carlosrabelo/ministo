// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ministo/ministo/internal/miner"
	"ministo/ministo/pkg/types"
)

const usage = `Ministo — Duino-Coin miner

Usage:
  ministo <command> [flags]

Commands:
  mine    Start mining
  help    Show this help

Examples:
  ministo mine -user alice
  ministo mine -user alice -emulate esp32
  ministo mine -user alice -emulate esp8266 -mac 24:0A:C4:12:34:56

Run 'ministo mine -h' for mine flags.
`

// Run parses args and executes the requested command.
// With no args (or help), it prints usage and returns nil.
func Run(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		fmt.Fprint(os.Stdout, usage)
		return nil
	}

	switch args[0] {
	case "mine":
		return runMine(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		fmt.Fprint(os.Stderr, usage)
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runMine(args []string) error {
	fs := flag.NewFlagSet("mine", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	user := fs.String("user", "", "Duino-Coin username (required)")
	proxyURL := fs.String("proxy", "socks5://localhost:9050", "SOCKS5 proxy URL")
	key := fs.String("key", "None", "mining key")
	rig := fs.String("rig", "pc", "rig identifier (Auto = derive from emulated MAC)")
	difficulty := fs.String("difficulty", "", "difficulty (empty = pool default or emulate default)")
	emulate := fs.String("emulate", "", "emulate board identity: esp32 or esp8266")
	macStr := fs.String("mac", "", "MAC for DUCOID (AA:BB:CC:DD:EE:FF); random Espressif OUI if empty")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ministo mine [flags]")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *user == "" {
		fs.Usage()
		return fmt.Errorf("-user is required")
	}

	cfg := &types.Config{
		Username:      *user,
		Difficulty:    *difficulty,
		MiningKey:     *key,
		RigIdentifier: *rig,
		MinerBanner:   "Ministo",
		MinerVersion:  "0.1",
	}

	if err := applyEmulate(cfg, *emulate, *macStr); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := miner.RunContext(ctx, cfg, *proxyURL)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func applyEmulate(cfg *types.Config, emulate, macStr string) error {
	emulate = strings.ToLower(strings.TrimSpace(emulate))
	if emulate == "" {
		if macStr != "" {
			return fmt.Errorf("-mac requires -emulate esp32 or esp8266")
		}
		return nil
	}

	var mac [6]byte
	var err error
	if macStr != "" {
		mac, err = miner.ParseMAC(macStr)
	} else {
		mac, err = miner.RandomEspressifMAC()
	}
	if err != nil {
		return err
	}

	switch emulate {
	case "esp32":
		cfg.DucoID = miner.ESP32DucoID(mac)
		if cfg.Difficulty == "" {
			cfg.Difficulty = "ESP32"
		}
		cfg.MinerBanner = "Official ESP32 Miner"
		cfg.MinerVersion = "3.5"
		if cfg.RigIdentifier == "" || cfg.RigIdentifier == "pc" || strings.EqualFold(cfg.RigIdentifier, "Auto") {
			cfg.RigIdentifier = "ESP32-" + miner.ESP32ChipHex(mac)
		}
	case "esp8266":
		cfg.DucoID = miner.ESP8266DucoID(mac)
		if cfg.Difficulty == "" {
			cfg.Difficulty = "ESP8266N"
		}
		cfg.MinerBanner = "Official ESP8266 Miner"
		cfg.MinerVersion = "3.5"
		if cfg.RigIdentifier == "" || cfg.RigIdentifier == "pc" || strings.EqualFold(cfg.RigIdentifier, "Auto") {
			cfg.RigIdentifier = "ESP8266-" + strings.ToUpper(miner.ESP8266ChipHex(mac))
		}
	default:
		return fmt.Errorf("unknown -emulate %q (want esp32 or esp8266)", emulate)
	}

	log.Printf("Emulating %s MAC %s → %s (rig %s, diff %s)",
		emulate, miner.FormatMAC(mac), cfg.DucoID, cfg.RigIdentifier, cfg.Difficulty)
	return nil
}

func isHelp(s string) bool {
	return s == "help" || s == "-h" || s == "--help" || s == "-help"
}
