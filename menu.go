/*
git clone https://github.com/abakum/menu
go mod init github.com/abakum/menu

go get github.com/eiannone/keyboard@latest
go get github.com/mitchellh/go-ps@latest
go get github.com/mattn/go-colorable@latest

go mod tidy
*/

package menu

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-ps"
)

const (
	SELECT = "Select" // for Prompt()
	MARK   = "("      // default option selected rune
	BUG    = "Ж"
	GT     = ">"
	MARKED = -1
	ITEM   = -2
	EXIT   = "\x00"
)

type (
	MenuFunc func(int, rune) string
	Static   string
)

var (
	Std = os.Stderr
)

// helper for prompt `Select`
func Prompt(int, rune) string {
	return SELECT
}

// helper for static prompt
func (s Static) Prompt(int, rune) string {
	return string(s)
}

/*
// template of helper for fooBar(anys...)
func fooBarMenu(index int, pressed rune, pref rune, suf string, marked, exit bool, anys ...any) string {
	r := rune(int(pref) + index)
	switch pressed {
	case r:
		fooBar(anys...)
		if exit {
			return menu.EXIT
		}
		return string(r)
	case menu.ITEM:
		return fmt.Sprintf("%c) %s", r, suf)
	case menu.MARKED:
		if marked {
			return menu.MARK
		}
	}
	return ""
}
*/

// Console menu
func Menu(def rune, // preselected item of menu
	keyEnter, // first run preselected menu item
	exitOnTypo bool, // exit from menu on typo
	items ...MenuFunc, // first item must be `Prompt` like
) {
	var (
		key     keyboard.Key
		err     error
		pressed rune
		index   = -1
		mark    string
	)
	bug, gt, out := BugGtOut()
exit:
	for {
		// set def by index. Used for arrow key navigation
		if index > -1 {
			def = 0
			if index < len(items) {
				rs := []rune(items[index+1](index, ITEM))
				if len(rs) > 0 {
					def = rs[0]
				}
			}
		}
		if def == 0 {
			for i, item := range items[1:] {
				s := item(i, MARKED) // is menu item marked?
				if s == "" {
					continue
				}
				rs := []rune(item(i, ITEM))
				if len(rs) < 1 {
					continue
				}
				def = rs[0]
			}
		}

		//print menu
		fmt.Fprintln(out)
		index = -1
		for i, item := range items[1:] {
			rs := []rune(item(i, ITEM)) //get menu item
			if len(rs) < 1 {
				continue
			}
			if def == 0 { //if def empty then select first item of menu
				def = rs[0]
			}
			if def == rs[0] {
				mark = gt
				index = i
			} else {
				mark = item(i, MARKED)
			}
			if mark == "" {
				mark = " "
			}

			fmt.Fprintf(out, "%s%s\n", mark, string(rs))
		}
		fmt.Fprint(out, items[0](index, def), gt)
		if keyEnter {
			pressed = def
		} else {
			pressed, key, err = keyboard.GetSingleKey()
			if err != nil {
				fmt.Fprintln(out, bug)
				return
			}
			if key == keyboard.KeyEnter {
				pressed = def
			}
		}
		keyEnter = false
		def = pressed
		if pressed == 0 {
			fmt.Fprintf(out, "0x%X\n", key)
			switch key {
			case keyboard.KeyEsc: // KeyEsc not typo
				break exit
			case keyboard.KeyHome:
				index = 0
				continue
			case keyboard.KeyArrowUp:
				if index == 0 {
					index = len(items) - 2
				} else {
					index--
				}
				continue
			case keyboard.KeyEnd:
				index = len(items) - 2
				continue
			case keyboard.KeyArrowDown:
				if index == len(items)-2 {
					index = 0
				} else {
					index++
				}
				continue
			}
		} else {
			fmt.Fprintf(out, "%c\n", def)
		}
		index = -1
		ok := false
	run:
		for i, item := range items[1:] {
			s := item(i, def)
			switch s {
			case "":
				continue
			case EXIT:
				break exit
			}
			def = []rune(s)[0] // allow item channge next def
			ok = true
			break run
		}
		if exitOnTypo && !ok {
			break exit
		}
		// on exit
	}
}

// no color need
func NoColor() bool {
	return os.Getenv("NO_COLOR") != "" ||
		os.Getenv("TERM") == "dumb" ||
		!isatty.IsTerminal(Std.Fd())
}

// get color dependents
func BugGtOut() (bug string, gt string, std io.Writer) {
	bug, gt, std = BUG, GT, Std
	if NoColor() {
		return
	}
	bug = fmt.Sprintf("\033[%d;%dm%s\033[%dm", color.Bold, color.BgRed, BUG, color.Reset)
	gt = fmt.Sprintf("\033[%d;%dm%s\033[%dm", color.Bold, color.FgGreen, GT, color.Reset)

	if IsAnsi() {
		return
	}
	std = colorable.NewColorable(Std)
	return
}

func PressAnyKey(s string, d time.Duration) {
	parent, err := ps.FindProcess(os.Getppid())
	if err == nil {
		for _, exe := range []string{"powershell.exe", "conemuc.exe", "cmd.exe"} {
			if strings.EqualFold(parent.Executable(), exe) {
				return
			}
		}
	}
	if d > 0 {
		time.AfterFunc(d, func() {
			keyboard.Close()
		})
	}
	bug, gt, out := BugGtOut()
	fmt.Fprint(out, s, gt)
	_, _, err = keyboard.GetSingleKey()
	if err != nil {
		fmt.Fprint(out, bug)
	}
	fmt.Fprintln(out)
}
