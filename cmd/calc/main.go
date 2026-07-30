// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pointlander/calc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}
		value := strings.TrimSpace(scanner.Text())
		if value == "exit" {
			return
		}

		cal := &calc.Calculator[uint32]{Buffer: value}
		cal.Init()
		if err := cal.Parse(); err != nil {
			fmt.Println(err)
			continue
		}
		result := cal.Eval()
		if result.Matrix != nil {
			fmt.Printf("%s\n", result.Matrix.String())
		} else {
			fmt.Printf("%s\n", result.Expression.String())
		}
	}
}
