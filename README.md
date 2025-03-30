# 🍾 linippet

linppet is tui tool which generate command from pre-registered your snipppet for bash/zsh.  
You can register snippet have dynamic arugments.


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

2. Choose your registered snippet. Then the snippet will be evaluated
```sh
lip
```

If you don't want the selected snippet to be evaluated, you can only output to A with the following command.

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

## LISENCE

[The MIT Licence](https://github.com/muleyuck/linippet/blob/main/LICENSE)
