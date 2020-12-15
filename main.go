// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/big"

	"github.com/c-bata/go-prompt"
	complex "github.com/pointlander/c0mpl3x"
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
		result := calc.Eval()
		x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
		x.SetRat(result)
		fmt.Printf("%s\n", x.String())
	}
}
