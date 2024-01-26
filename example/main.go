package main

import (
	"fmt"
	"strings"

	"github.com/abakum/menu"
)

// helper for simple menu
func fooBar(index int, pressed rune, pref rune, suf string, marked, exit bool) string {
	r := rune(int(pref) + index)
	switch pressed {
	case r:
		fmt.Printf("Option %q is run\n", suf)
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

func main() {
	fmt.Println("\nSimple print menu\n with index start from 0\n preselected 1\n static prompt\n exit on typo")
	items := []menu.MenuFunc{menu.Static("Choose").Prompt}
	items = append(items, func(index int, pressed rune) string {
		return fooBar(index, pressed, '0', "foo", false, false)
	})
	items = append(items, func(index int, pressed rune) string {
		return fooBar(index, pressed, '0', "bar", false, false)
	})
	menu.Menu('1', false, true, items...)

	fmt.Println("\n\n\nPrint menu\n with index start from 1\n preselected by items\n not static prompt\n not exit on typo\n exit on `bar`")
	items = []menu.MenuFunc{func(index int, pressed rune) string {
		return fmt.Sprintf("Press %c to select option %d", pressed, index)
	}}
	items = append(items, func(index int, pressed rune) string {
		return fooBar(index, pressed, '1', "foo", menu.IsAnsi(), false)
	})
	items = append(items, func(index int, pressed rune) string {
		return fooBar(index, pressed, '1', "foo", menu.IsAnsi(), true)
	})
	menu.Menu(0, false, true, items...)

	fmt.Println("\n\n\nPrint menu\n selected by cyrillic letters with typo tolerant\n first run `foo`\n not exit on typo")
	items = []menu.MenuFunc{menu.Prompt}
	items = append(items, func(index int, pressed rune) string {
		r := rune('Ю' + index)
		alt := rune('1' + index)
		switch {
		case strings.EqualFold(string(r), string(pressed)) || pressed == '>' || pressed == '.' || pressed == alt:
			fmt.Printf("Option %q is run\n", "foo")
			return string(r)
		case pressed == menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "foo")
		}
		return "" // not for me
	})
	items = append(items, func(index int, pressed rune) string {
		r := rune('Ю' + index)
		alt := rune('1' + index)
		switch {
		case strings.EqualFold(string(r), string(pressed)) || strings.EqualFold("z", string(pressed)) || pressed == alt:
			fmt.Printf("Option %q is run\n", "bar")
			return string(r)
		case pressed == menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "bar")
		}
		return "" // not for me
	})
	menu.Menu('Ю', true, false, items...)
}
