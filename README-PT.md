# Ministo

Minerador Duino-Coin em Go (Ministo é a palavra em esperanto para “minerador”). Feito como projeto de aprendizado; usável com um proxy SOCKS5.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue.svg)](https://go.dev/)

[English](README.md)

## Destaques

- Mina Duino-Coin por proxy SOCKS5 (padrão amigável ao Tor)
- CLI com help por padrão e comando `mine`
- Flags para username, mining key, proxy, dificuldade e nome do rig
- Emulação opcional de identidade ESP32/ESP8266 (`DUCOID` a partir do MAC)
- Build/install via Make; binário em `bin/`

## Pré-requisitos

- **Go 1.18+** — para compilar a partir do código; [download](https://go.dev/dl/)
- **Proxy SOCKS5** — padrão `socks5://localhost:9050` (ex.: Tor)

## Instalação

### Compilar a partir do código

```bash
make build
```

Instale em `~/.local/bin` (padrão), ou em `/usr/local/bin` no sistema (sudo só para a cópia):

```bash
make install
make install-system
make uninstall
make uninstall-system
```

## Uso

Sem argumentos (ou `--help`) mostra o help:

```bash
./bin/ministo
./bin/ministo --help
./bin/ministo mine -h
```

### Minerar

```bash
./bin/ministo mine -user SEU_USERNAME
./bin/ministo mine -user SEU_USERNAME -proxy socks5://127.0.0.1:9050 -rig pc -key None
```

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `-user` | _(obrigatório)_ | Username Duino-Coin |
| `-proxy` | `socks5://localhost:9050` | URL do proxy SOCKS5 |
| `-key` | `None` | Mining key |
| `-rig` | `pc` | Identificador do rig (`Auto` com `-emulate` deriva do MAC) |
| `-difficulty` | _(vazio)_ | Tier de dificuldade do pool; definido por `-emulate` se omitido |
| `-emulate` | _(off)_ | `esp32` ou `esp8266` |
| `-mac` | _(aleatório)_ | MAC para `DUCOID` quando `-emulate` está ativo |

### Emular ESP32 / ESP8266

Define banner do minerador, tier de dificuldade e `DUCOID` do share como no firmware oficial ESP:

```bash
./bin/ministo mine -user SEU_USERNAME -emulate esp32
./bin/ministo mine -user SEU_USERNAME -emulate esp8266
./bin/ministo mine -user SEU_USERNAME -emulate esp32 -mac 24:0A:C4:12:34:56
```

Sem `-mac`, um MAC aleatório com OUI Espressif é gerado para a sessão.

## Estrutura do Projeto

```
ministo/cmd/ministo/     # Entry point
ministo/internal/cli/    # Commands and flags
ministo/internal/miner/  # Jobs, hashing, DUCOID / ESP emulate
ministo/internal/proxy/  # SOCKS5 client
ministo/pkg/types/       # Shared types
bin/                     # Build output (git-ignored)
.make/                   # Build and install scripts
```

## Desenvolvimento

```bash
make build             # Compile binary to bin/ministo
make test              # Run all tests
make quality           # Format, vet, and lint
make install           # Install binary to ~/.local/bin
make install-system    # Install binary to /usr/local/bin
make uninstall         # Remove from ~/.local/bin
make uninstall-system  # Remove from /usr/local/bin
```

## Licença

Este projeto está sob a licença MIT — veja [LICENSE](LICENSE) para detalhes.
