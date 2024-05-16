package menu

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
)

// for old windows `choco install ansicon`
func IsAnsi() (ok bool) {
	if os.Getenv("ANSICON") != "" || isatty.IsCygwinTerminal(Std.Fd()) {
		return true
	}
	parent, err := ps.FindProcess(os.Getppid())
	if err != nil {
		fmt.Println(BUG, err)
		return
	}
	ma, mi, _ := windows.RtlGetNtVersionNumbers()
	ae := []string{"ansicon.exe", "conemuc.exe"}
	if ma*10+mi > 61 { // after win7
		ae = append(ae, "powershell.exe")
	}
	for _, exe := range ae {
		ok = strings.EqualFold(parent.Executable(), exe)
		if ok {
			break
		}
	}
	return
}
