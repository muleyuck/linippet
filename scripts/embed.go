package scripts

import _ "embed"

//go:embed initializer.bash
var InitializeBashScript string

//go:embed initializer.zsh
var InitializeZShellScript string

//go:embed app_version
var AppVersion string
