package scripts

import _ "embed"

//go:embed initializer.bash
var InitializeBashScript string

//go:embed initializer.zsh
var InitializeZShellScript string
