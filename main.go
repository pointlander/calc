package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the application"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	for {
		value := prompt.Input("> ", completer)
		if value == "exit" {
			return
		}

		calc := &Calculator{Buffer: value}
		calc.Init()
		if err := calc.Parse(); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(value + " = " + calc.Eval().String())
	}
}
