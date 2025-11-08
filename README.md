[![unit-test](https://github.com/muleyuck/linippet/actions/workflows/unit-test.yml/badge.svg)](https://github.com/muleyuck/linippet/actions/workflows/unit-test.yml)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![Release](https://img.shields.io/github/release/muleyuck/linippet.svg)](https://github.com/muleyuck/linippet/releases/latest)
[![GoDoc](https://godoc.org/github.com/muleyuck/linippet?status.svg)](https://godoc.org/github.com/muleyuck/linippet) 


# 🍾 linippet

linppet is tui tool which generate command from pre-registered your snipppet for bash/zsh.  
You can register snippet have dynamic arguments.

![demo](https://github.com/user-attachments/assets/a65dc0b4-a436-4fe5-b604-b85f0dd35375)


## Installation

### 1. Install binary  
 Go
```sh
go install github.com/muleyuck/linippet@latest
```
Homebrew
```sh
brew install muleyuck/tap/linippet
```
Or, download binary from [Releases](https://github.com/muleyuck/linippet/releases)

### 2. Setup you shell
zsh
```sh
eval "$(linippet init zsh)"
```
bash
```sh
eval "$(linippet init bash)"
```

## Get Started

### Quick Start
1. Please enter a snippet which you want to use in Modal Tui
```sh
linippet create
```

2. Choose your registered snippet. Then the snippet will be evaluated in your shell.
```sh
lip
```

If you don't want the selected snippet to be evaluated, you can only output to stdout with the following command.

```sh
linippet
```

### Triggered by a bind key

export `LINIPPET_TRIGGER_BIND_KEY` environment value. For example: `^o`
 ```sh
export LINIPPET_TRIGGER_BIND_KEY="^o"
```
linippet will be triggered when press `^o` key. And your selected snippet will set current readline on shell.

### CRUD snippet

```sh
linippet [create|edit|remove]
```

## Inspired by

linippet inspired by [Warp workflow](https://docs.warp.dev/features/warp-drive/workflows)

## LICENCE

[The MIT Licence](https://github.com/muleyuck/linippet/blob/main/LICENSE)
