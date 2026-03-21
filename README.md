[![unit-test](https://github.com/muleyuck/linippet/actions/workflows/unit-test.yml/badge.svg)](https://github.com/muleyuck/linippet/actions/workflows/unit-test.yml)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![Release](https://img.shields.io/github/release/muleyuck/linippet.svg)](https://github.com/muleyuck/linippet/releases/latest)
[![GoDoc](https://godoc.org/github.com/muleyuck/linippet?status.svg)](https://godoc.org/github.com/muleyuck/linippet)
[![GitHub stars](https://img.shields.io/github/stars/muleyuck/linippet?style=flat-square)](https://github.com/muleyuck/linippet/stargazers)


# 🍾 linippet

Never forget your one-liner commands again.
**linippet** is a TUI snippet manager for bash/zsh — store, fuzzy-search, and execute shell commands with dynamic arguments (`${{arg_name}}`).

![demo](https://github.com/user-attachments/assets/a65dc0b4-a436-4fe5-b604-b85f0dd35375)


## Installation

### 1. Install binary
Shell script (Linux/macOS)
```sh
curl -sSfL https://raw.githubusercontent.com/muleyuck/linippet/main/install.sh | sh
```
Go
```sh
go install github.com/muleyuck/linippet@latest
```
Homebrew
```sh
brew install muleyuck/tap/linippet
```
Or, download binary from [Releases](https://github.com/muleyuck/linippet/releases)

### 2. Setup your shell
zsh
```sh
eval "$(linippet init zsh)"
```
bash
```sh
eval "$(linippet init bash)"
```

## Features

- **Fuzzy search** — quickly find snippets from your list
- **Dynamic arguments** — use `${{arg_name}}` placeholders, filled interactively at run time
- **Default values** — use `${{arg_name:default}}` to pre-fill arguments
- **Keybind trigger** — invoke linippet from anywhere in your shell with a single key chord
- **Vim / Emacs navigation** — familiar key bindings in the TUI

## Get Started

### Quick Start

1. Register a snippet:
```sh
linippet create
```

2. Search and execute:
```sh
lip
```

To output the selected snippet to stdout without executing it:
```sh
linippet
```

### Triggered by a bind key

Set `LINIPPET_TRIGGER_BIND_KEY` to invoke linippet mid-command. For example:
```sh
export LINIPPET_TRIGGER_BIND_KEY="^o"
```
Pressing `Ctrl+o` will open the TUI and paste the selected snippet into your current readline.

### CRUD snippets

```sh
linippet [create|edit|remove]
```

## Inspired by

linippet inspired by [Warp workflow](https://docs.warp.dev/features/warp-drive/workflows)

## License

[MIT](https://github.com/muleyuck/linippet/blob/main/LICENSE)
