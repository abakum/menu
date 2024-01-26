/*
git clone https://github.com/abakum/menu
go mod init github.com/abakum/menu

go get github.com/eiannone/keyboard@latest
go get github.com/mitchellh/go-ps@latest

go mod tidy

// usage
func main() {
 	items := []menu.MenuFunc{menu.Prompt}
	items = append(items, func(index int, pressed rune) string {
		r := rune('1' + index) // menu starts with 1)
		switch pressed {
		case r:
			foo() // run
			if exit {
				return menu.EXIT
			}
			return string(r) //new def
		case menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "foo")
		}
		return "" // not for me
	})
	menu.Menu('1', false, true, items...)
}
*/

package menu

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
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
	bug = BUG
	gt  = GT
)

// helper for prompt `Select`
func Prompt(int, rune) string {
	return SELECT
}

// helper for static prompt
func (s Static) Prompt(int, rune) string {
	return string(s)
}

// Console menu
func Menu(def rune, // preselected item of menu
	keyEnter, // first run preselected menu item
	exitOnTypo bool, // exit from menu on typo
	items ...MenuFunc, // first item must be `Prompt` like
) {
	const (
		ansiReset     = "\u001B[0m"
		ansiRedBGBold = "\u001B[41m\u001B[1m"
		ansiGreenFG   = "\u001B[32m\u001B[1m"
	)
	var (
		key     keyboard.Key
		err     error
		pressed rune
		index   = -1
		mark    string
	)
	if IsColor() {
		bug = ansiRedBGBold + BUG + ansiReset
		gt = ansiGreenFG + GT + ansiReset
	}
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
		fmt.Println()
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

			fmt.Printf("%s%s\n", mark, string(rs))
		}
		fmt.Print(items[0](index, def), gt)
		if keyEnter {
			pressed = def
		} else {
			pressed, key, err = keyboard.GetSingleKey()
			if err != nil {
				fmt.Println(bug)
				return
			}
			if key == keyboard.KeyEnter {
				pressed = def
			}
		}
		keyEnter = false
		def = pressed
		if pressed == 0 {
			fmt.Printf("0x%X\n", key)
			switch key {
			case keyboard.KeyEsc:
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
			fmt.Printf("%c\n", def)
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

// for old windows `choco install ansicon`
func IsAnsi() (ok bool) {
	if runtime.GOOS != "windows" {
		return true
	}
	if os.Getenv("ANSICON") != "" {
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

// is console color enable
func IsColor() bool {
	return os.Getenv("NO_COLOR") == "" && IsAnsi()
}
