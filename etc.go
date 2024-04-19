//go:build !windows
// +build !windows

package menu

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsAnsi() (ok bool) {
	return isatty.IsTerminal(os.Stdout.Fd())
}
