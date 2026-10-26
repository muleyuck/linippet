package scripts

import (
	"os"
	"os/exec"
	"testing"
)

func TestLipForwardsArguments(t *testing.T) {
	for _, shell := range []string{"bash", "zsh"} {
		t.Run(shell, func(t *testing.T) {
			if _, err := exec.LookPath(shell); err != nil {
				t.Skipf("%s is not installed", shell)
			}
			for _, args := range [][]string{nil, {"--", "echo テスト🌱; $(false) *"}} {
				script := `
linippet() {
    for arg in "$@"; do printf '%s\000' "$arg" >&2; done
    printf ':'
}
source "$1"
shift
linippet_apply "$@"
`
				cmd := exec.Command(shell, append([]string{"-c", script, "test", "./initializer." + shell}, args...)...)
				cmd.Env = append(os.Environ(), "LINIPPET_TRIGGER_BIND_KEY=")
				got, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("args %q: %v\n%s", args, err, got)
				}
				want := ""
				for _, arg := range args {
					want += arg + "\x00"
				}
				if string(got) != want {
					t.Errorf("args %q: forwarded = %q, want %q", args, got, want)
				}
			}
		})
	}
}
