/*
git clone https://github.com/abakum/menu
go mod init github.com/abakum/menu

go get github.com/eiannone/keyboard@latest
go get github.com/mitchellh/go-ps@latest

go mod tidy
*/

package menu

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/mitchellh/go-ps"
)

const (
	MENU    = "Select"
	MARK    = '(' // default option selected rune
	BUG     = "Ж"
	GT      = ">"
	RunFunc = -1
)

type (
	PromptFunc func(int, rune) string
	ItemFunc   func(int) string
)

var (
	Bug = BUG
	Gt  = GT
)

// helper for static prompt
func Prompt(index int, def rune) string {
	return MENU
}

// Console menu
func Menu(prompt PromptFunc,
	mark, def rune,
	keyEnter, exitOnTypo bool,
	items ...ItemFunc) {
	const (
		ansiReset     = "\u001B[0m"
		ansiRedBGBold = "\u001B[41m\u001B[1m"
		ansiGreenFG   = "\u001B[32m\u001B[1m"
	)
	var (
		key   keyboard.Key
		err   error
		r     rune
		index = -1
	)
	if os.Getenv("NO_COLOR") == "" && IsAnsi() {
		Bug = ansiRedBGBold + BUG + ansiReset
		Gt = ansiGreenFG + GT + ansiReset
	}

	for {
		newD := false
		// search mark or set GT by index
		for i, item := range items {
			rs := []rune(item(i)) //get menu item. Item my use i for assign key
			if len(rs) < 1 {
				continue
			}
			newD = rs[0] == mark
			if newD {
				if len(rs) < 2 {
					continue
				}
				rs = rs[1:]
			}
			if index > -1 { //set GT by index. Used for arrow key navigation
				if index == i {
					def = rs[0]
				}
			} else { //set GT by item
				if newD {
					def = rs[0]
				}
			}
		}

		//print menu
		fmt.Println()
		index = -1
		for i, item := range items {
			rs := []rune(item(i)) //get menu item
			if len(rs) < 1 {
				continue
			}
			m := " "
			if rs[0] == mark { // new def
				if len(rs) < 2 {
					continue
				}
				m = string(mark)
				rs = rs[1:]
			}
			if def == rs[0] {
				m = Gt
				index = i
			}
			fmt.Printf("%s%s\n", m, string(rs))
		}
		fmt.Print(prompt(index, def), Gt)
		if keyEnter {
			r = def
		} else {
			r, key, err = keyboard.GetSingleKey()
			if err != nil {
				fmt.Println(Bug)
				return
			}
			if key == keyboard.KeyEnter {
				r = def
			}
		}
		keyEnter = false
		def = r
		if r == 0 {
			fmt.Printf("0x%X\n", key)
			switch key {
			case keyboard.KeyHome:
				index = 0
				continue
			case keyboard.KeyArrowUp:
				if index == 0 {
					index = len(items) - 1
				} else {
					index--
				}
				continue
			case keyboard.KeyEnd:
				index = len(items) - 1
				continue
			case keyboard.KeyArrowDown:
				if index == len(items)-1 {
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
	doit:
		for i, item := range items {
			rs := []rune(item(i)) //get menu item
			if len(rs) < 1 {
				continue
			}
			if rs[0] == mark { //ignore mark from item
				if len(rs) < 2 {
					continue
				}
				rs = rs[1:]
			}
			ok = def == rs[0]
			if ok {
				if len(item(RunFunc)) > 0 { // run func of menu item
					return // for once selected menu
				}
				break doit
			}
		}
		if exitOnTypo && !ok {
			return
		}
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
	for _, exe := range []string{
		"powershell.exe",
		"ansicon.exe",
		"conemuc.exe"} {
		ok = strings.ToLower(parent.Executable()) == exe
		if ok {
			break
		}
	}
	return
}
