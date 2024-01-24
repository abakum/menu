package main

import (
	"fmt"

	"github.com/abakum/menu"
)

func main() {
	fmt.Println("simple print menu with index start from 0 preselected 1")
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

	fmt.Println("simple print menu with index start from 1 preselected from items")
	items = []menu.ItemFunc{}
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("foo") // simple print choose
			return ""
		}
		mark := ""
		if menu.IsAnsi() {
			mark = string([]rune(menu.GT)[0])
		}

		return fmt.Sprintf("%s%d) %s", mark, index+1, "foo") // print item of menu
	})
	items = append(items, func(index int) string {
		if index == menu.RunFunc {
			fmt.Println("bar") // simple print choose
			return ""
		}
		return fmt.Sprintf("%d) %s", index+1, "bar") // print item of menu
	})
	menu.Menu(menu.Prompt, []rune(menu.GT)[0], 0, false, true, items...)

}
