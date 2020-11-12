// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/big"
	"strings"

	"github.com/ALTree/bigfloat"
)

func (c *Calculator) Eval() *big.Rat {
	return c.Rulee(c.AST())
}

func (c *Calculator) Rulee(node *node32) *big.Rat {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulee1:
			return c.Rulee1(node)
		}
		node = node.next
	}
	return nil
}

func (c *Calculator) Rulee1(node *node32) *big.Rat {
	node = node.up
	var a *big.Rat
	for node != nil {
		switch node.pegRule {
		case rulee2:
			a = c.Rulee2(node)
		case ruleadd:
			node = node.next
			b := c.Rulee2(node)
			a.Add(a, b)
		case ruleminus:
			node = node.next
			b := c.Rulee2(node)
			a.Sub(a, b)
		}
		node = node.next
	}
	return a
}

func (c *Calculator) Rulee2(node *node32) *big.Rat {
	node = node.up
	var a *big.Rat
	for node != nil {
		switch node.pegRule {
		case rulee3:
			a = c.Rulee3(node)
		case rulemultiply:
			node = node.next
			b := c.Rulee3(node)
			a.Mul(a, b)
		case ruledivide:
			node = node.next
			b := c.Rulee3(node)
			a.Quo(a, b)
		case rulemodulus:
			node = node.next
			b := c.Rulee3(node)
			if a.Denom().Cmp(big.NewInt(1)) == 0 && b.Denom().Cmp(big.NewInt(1)) == 0 {
				a.Num().Mod(a.Num(), b.Num())
			}
		}
		node = node.next
	}
	return a
}

func (c *Calculator) Rulee3(node *node32) *big.Rat {
	node = node.up
	var a *big.Rat
	for node != nil {
		switch node.pegRule {
		case rulee4:
			a = c.Rulee4(node)
		case ruleexponentiation:
			node = node.next
			b := c.Rulee4(node)
			x, _, _ := big.ParseFloat(a.FloatString(1024), 10, 1024, big.ToNearestEven)
			y, _, _ := big.ParseFloat(b.FloatString(1024), 10, 1024, big.ToNearestEven)
			pow := bigfloat.Pow(x, y)
			a.SetString(pow.String())
		}
		node = node.next
	}
	return a
}

func (c *Calculator) Rulee4(node *node32) *big.Rat {
	node = node.up
	minus := false
	for node != nil {
		switch node.pegRule {
		case rulevalue:
			a := c.Rulevalue(node)
			if minus {
				a.Neg(a)
			}
			return a
		case ruleminus:
			minus = true
		}
		node = node.next
	}
	return nil
}

func (c *Calculator) Rulevalue(node *node32) *big.Rat {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulenumber:
			a := big.NewRat(1, 1)
			a.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
			return a
		case rulesub:
			return c.Rulesub(node)
		}
		node = node.next
	}
	return nil
}

func (c *Calculator) Rulesub(node *node32) *big.Rat {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulee1:
			return c.Rulee1(node)
		}
		node = node.next
	}
	return nil
}
