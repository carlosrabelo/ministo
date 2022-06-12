# Ministo

Duino-Coin miner in Go (Ministo is Esperanto for “miner”). Built as a learning project; fully usable with a SOCKS5 proxy.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue.svg)](https://go.dev/)

[Português](README-PT.md)

## Highlights

- Mines Duino-Coin through a SOCKS5 proxy (Tor-friendly default)
- CLI with help by default and a `mine` command
- Flags for username, mining key, proxy, difficulty, and rig name
- Optional ESP32/ESP8266 identity emulation (`DUCOID` from MAC)
- Make-based build/install; binary under `bin/`

## Prerequisites

- **Go 1.18+** — to build from source; [download](https://go.dev/dl/)
- **SOCKS5 proxy** — default `socks5://localhost:9050` (e.g. Tor)

## Installation

### Build from source

```bash
make build
```

Install to `~/.local/bin` (default), or system-wide to `/usr/local/bin` (sudo only for the copy):

```bash
make install
make install-system
make uninstall
make uninstall-system
```

## Usage

No arguments (or `--help`) prints help:

```bash
./bin/ministo
./bin/ministo --help
./bin/ministo mine -h
```

### Mine

```bash
./bin/ministo mine -user YOUR_USERNAME
./bin/ministo mine -user YOUR_USERNAME -proxy socks5://127.0.0.1:9050 -rig pc -key None
```

| Flag | Default | Description |
|------|---------|-------------|
| `-user` | _(required)_ | Duino-Coin username |
| `-proxy` | `socks5://localhost:9050` | SOCKS5 proxy URL |
| `-key` | `None` | Mining key |
| `-rig` | `pc` | Rig identifier (`Auto` with `-emulate` derives from MAC) |
| `-difficulty` | _(empty)_ | Pool difficulty tier; set by `-emulate` if omitted |
| `-emulate` | _(off)_ | `esp32` or `esp8266` |
| `-mac` | _(random)_ | MAC for `DUCOID` when `-emulate` is set |

### Emulate ESP32 / ESP8266

Sets miner banner, difficulty tier, and share `DUCOID` like the official ESP firmwares:

```bash
./bin/ministo mine -user YOUR_USERNAME -emulate esp32
./bin/ministo mine -user YOUR_USERNAME -emulate esp8266
./bin/ministo mine -user YOUR_USERNAME -emulate esp32 -mac 24:0A:C4:12:34:56
```

Without `-mac`, a random Espressif-OUI MAC is generated for the session.

## Project Layout

```
ministo/cmd/ministo/     # Entry point
ministo/internal/cli/    # Commands and flags
ministo/internal/miner/  # Jobs, hashing, DUCOID / ESP emulate
ministo/internal/proxy/  # SOCKS5 client
ministo/pkg/types/       # Shared types
bin/                     # Build output (git-ignored)
.make/                   # Build and install scripts
```

## Development

```bash
make build             # Compile binary to bin/ministo
make test              # Run all tests
make quality           # Format, vet, and lint
make install           # Install binary to ~/.local/bin
make install-system    # Install binary to /usr/local/bin
make uninstall         # Remove from ~/.local/bin
make uninstall-system  # Remove from /usr/local/bin
```

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.
