// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exp", Description: "The natural number raised to a value"},
		{Text: "pi", Description: "The constant PI"},
		{Text: "prec", Description: "Sets the precision for calculations"},
		{Text: "derivative", Description: "Computes the symbolic derivative of the expression"},
		{Text: "log", Description: "The natural logarithm of the input"},
		{Text: "sqrt", Description: "The square root of the value"},
		{Text: "cos", Description: "The cosine of the value"},
		{Text: "sin", Description: "The sine of the value"},
		{Text: "tan", Description: "The tangent of the value"},
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
		result := calc.Eval()
		fmt.Printf("%s\n", result.String())
	}
}
