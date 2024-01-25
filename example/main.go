package main

import (
	"fmt"

	"github.com/abakum/menu"
)

func main() {
	fmt.Println("Simple print menu with index start from 0 preselected 1")
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

	fmt.Println("\n\n\nPrint menu with index start from 1 preselected from items and not static prompt")
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

	fmt.Println("\n\n\nPrint menu selected by cyrillic letters and prerun `foo`")
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
