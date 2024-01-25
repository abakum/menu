package main

import (
	"fmt"
	"strings"

	"github.com/abakum/menu"
)

func main() {
	fmt.Println("\nSimple print menu\n with index start from 0\n preselected 1\n static prompt\n exit on typo")
	items := []menu.MenuFunc{menu.Prompt}
	items = append(items, func(index int, pressed rune) string {
		r := rune('0' + index)
		switch pressed {
		case menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "foo")
		case r:
			fmt.Println("foo") // run
			return string(r)   //new def
		}
		return "" // not for me
	})
	items = append(items, func(index int, pressed rune) string {
		r := rune('0' + index)
		switch pressed {
		case menu.ITEM: // print item of menu
			return fmt.Sprintf("%c) %s", r, "bar")
		case r:
			fmt.Println("bar") // run
			return string(r)
		}
		return "" // not for me
	})
	menu.Menu('1', false, true, items...)

	fmt.Println("\n\n\nPrint menu\n with index start from 1\n preselected by items\n not static prompt\n not exit on typo\n exit on `bar`")
	items = []menu.MenuFunc{func(index int, pressed rune) string {
		return fmt.Sprintf("Press %c", pressed)
	}}
	items = append(items, func(index int, pressed rune) string {
		r := rune('1' + index)
		switch pressed {
		case menu.MARKED: // marked
			if menu.IsAnsi() {
				return menu.MARK
			}
		case menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "foo")
		case r:
			fmt.Println("foo") // run
			return string(r)
		}
		return "" // not for me
	})
	items = append(items, func(index int, pressed rune) string {
		r := rune('1' + index)
		switch pressed {
		case menu.ITEM: // print item of menu
			return fmt.Sprintf("%c) %s", r, "bar")
		case r:
			fmt.Println("bar") // run
			return menu.EXIT
		}
		return "" // not for me
	})
	menu.Menu(0, false, true, items...)

	fmt.Println("\n\n\nPrint menu\n selected by cyrillic letters with typo tolerant\n first run `foo`\n not exit on typo")
	items = []menu.MenuFunc{menu.Prompt}
	items = append(items, func(index int, pressed rune) string {
		r := rune('Ю' + index)
		alt := rune('1' + index)
		switch {
		case pressed == menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "foo")
		case strings.EqualFold(string(r), string(pressed)) || pressed == '>' || pressed == '.' || pressed == alt:
			fmt.Println("foo") // run
			return string(r)
		}
		return "" // not for me
	})
	items = append(items, func(index int, pressed rune) string {
		r := rune('Ю' + index)
		alt := rune('1' + index)
		switch {
		case pressed == menu.ITEM: // item of menu
			return fmt.Sprintf("%c) %s", r, "bar")
		case strings.EqualFold(string(r), string(pressed)) || strings.EqualFold("z", string(pressed)) || pressed == alt:
			fmt.Println("bar") // run
			return string(r)
		}
		return "" // not for me
	})
	menu.Menu('Ю', true, false, items...)
}
