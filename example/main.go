package main

import (
	"fmt"

	"github.com/abakum/menu"
)

func main() {
	fmt.Println("\nSimple print menu\n with index start from 0\n preselected 1\n static prompt\n exit on typo")
	items := []menu.ItemFunc{}
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("foo") // simple print choose
			return ""
		}
		return fmt.Sprintf("%d) %s", index, "foo") // print item of menu
	})
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("bar") // simple print choose
			return ""
		}
		return fmt.Sprintf("%d) %s", index, "bar") // print item of menu
	})
	menu.Menu(menu.Prompt, menu.MARK, '1', false, true, items...)

	fmt.Println("\n\n\nPrint menu\n with index start from 1\n preselected by items\n not static prompt\n not exit on typo\n exit on `bar`\n with custom `mark`")
	items = []menu.ItemFunc{}
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("foo") // simple print choose
			return ""
		}

		mark := ""
		if menu.IsAnsi() {
			mark = ">"
		}
		return fmt.Sprintf("%s%d) %s", mark, index+1, "foo") // print item of menu
	})
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("bar") // simple print choose
			return "exit"      // and exit if return not emty string
		}
		return fmt.Sprintf("%d) %s", index+1, "bar") // print item of menu
	})
	menu.Menu(func(index int, def rune) string {
		return fmt.Sprintf("Press %s for %s", string(def), items[index](index))
	}, '>', 0, false, false, items...)

	fmt.Println("\n\n\nPrint menu\n selected by cyrillic letters\n first run `foo`\n not exit on typo")
	items = []menu.ItemFunc{}
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("foo") // simple print choose
			return ""
		}
		return fmt.Sprintf("%c) %s", 'Ю'+index, "foo") // print item of menu
	})
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("bar") // simple print choose
			return ""
		}
		return fmt.Sprintf("%c) %s", 'Ю'+index, "bar") // print item of menu
	})
	menu.Menu(menu.Prompt, menu.MARK, 'Ю', true, false, items...)
}
