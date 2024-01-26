package main

import (
	"fmt"
	"strings"

	"github.com/abakum/menu"
)

// helper for simple menu
func fooBarMenu(index int, pressed rune, pref rune, suf string, marked, exit bool) string {
	r := rune(int(pref) + index)
	switch pressed {
	case r:
		fooBar(suf)
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

func fooBar(suf string) {
	fmt.Printf("Option %q is run\n", suf)
}

// helper for menu
func fooBarAltMenu(index int, pressed rune, pref rune, suf string, marked, exit bool, validKey func(pressed rune) bool) string {
	r := rune(int(pref) + index)
	alt := rune('1' + index)
	switch {
	case pressed == menu.ITEM:
		return fmt.Sprintf("%c) %s", r, suf)
	case pressed == menu.MARKED:
		if marked {
			return menu.MARK
		}
	case strings.EqualFold(string(pressed), string(r)) || pressed == alt || validKey(pressed):
		fooBar(suf)
		if exit {
			return menu.EXIT
		}
		return string(r)
	}
	return ""
}

func main() {
	fmt.Println(`
Simple print menu
 with index start from 0
 preselected 1
 static prompt
 exit on typo`)
	items := []menu.MenuFunc{menu.Static("Choose").Prompt}
	items = append(items, func(index int, pressed rune) string {
		return fooBarMenu(index, pressed, '0', "foo", false, false)
	})
	items = append(items, func(index int, pressed rune) string {
		return fooBarMenu(index, pressed, '0', "bar", false, false)
	})
	menu.Menu('1', false, true, items...)

	fmt.Println(`

Print menu
 with index start from 1
 preselected by items
 not static prompt
 not exit on typo
 exit on "bar"`)
	items = []menu.MenuFunc{func(index int, pressed rune) string {
		if index == -1 || pressed == 0 {
			return menu.SELECT
		}
		return fmt.Sprintf("Press %c to select option %d", pressed, index)
	}}
	items = append(items, func(index int, pressed rune) string {
		return fooBarMenu(index, pressed, '1', "foo", menu.IsAnsi(), false)
	})
	items = append(items, func(index int, pressed rune) string {
		return fooBarMenu(index, pressed, '1', "bar", false, true)
	})
	menu.Menu(0, false, false, items...)

	fmt.Println(`

Print menu
 selected by cyrillic letters with typo tolerant
 first run "foo"
 not exit on typo`)
	items = []menu.MenuFunc{menu.Prompt}
	items = append(items, func(index int, pressed rune) string {
		return fooBarAltMenu(index, pressed, 'Ю', "foo", false, false, func(pressed rune) bool {
			return pressed == '>' || pressed == '.'
		})
	})
	items = append(items, func(index int, pressed rune) string {
		return fooBarAltMenu(index, pressed, 'Ю', "bar", false, false, func(pressed rune) bool {
			return strings.EqualFold(string(pressed), "z")
		})
	})
	menu.Menu('Ю', true, false, items...)
}
