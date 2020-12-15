// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/big"
	"strings"

	"github.com/ALTree/bigfloat"
	complex "github.com/pointlander/c0mpl3x"
)

var prec uint = 1024

func (c *Calculator) Eval() *complex.Rational {
	return c.Rulee(c.AST())
}

func (c *Calculator) Rulee(node *node32) *complex.Rational {
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

func (c *Calculator) Rulee1(node *node32) *complex.Rational {
	node = node.up
	var a *complex.Rational
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

func (c *Calculator) Rulee2(node *node32) *complex.Rational {
	node = node.up
	var a *complex.Rational
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
			a.Div(a, b)
		case rulemodulus:
			node = node.next
			b := c.Rulee3(node)
			if a.A.Denom().Cmp(big.NewInt(1)) == 0 && b.A.Denom().Cmp(big.NewInt(1)) == 0 {
				a.A.Num().Mod(a.A.Num(), b.A.Num())
			}
		}
		node = node.next
	}
	return a
}

func (c *Calculator) Rulee3(node *node32) *complex.Rational {
	node = node.up
	var a *complex.Rational
	for node != nil {
		switch node.pegRule {
		case rulee4:
			a = c.Rulee4(node)
		case ruleexponentiation:
			node = node.next
			b := c.Rulee4(node)
			x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
			x.SetRat(a)
			y := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
			y.SetRat(b)
			x.Pow(x, y).Rat(a)
		}
		node = node.next
	}
	return a
}

func (c *Calculator) Rulee4(node *node32) *complex.Rational {
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

func (c *Calculator) Rulevalue(node *node32) *complex.Rational {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case ruleimaginary:
			a := complex.NewRational(big.NewRat(0, 1), big.NewRat(1, 1))
			a.B.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
			return a
		case rulenumber:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a.A.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
			return a
		case ruleexp1:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					x.SetRat(a)
					x.Exp(x).Rat(a)
					return a
				}
				node = node.next
			}
		case ruleexp2:
			node := node.up
			for node != nil {
				if node.pegRule == rulevalue {
					a := c.Rulevalue(node)
					x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					x.SetRat(a)
					x.Exp(x).Rat(a)
					return a
				}
				node = node.next
			}
		case rulepi:
			a := big.NewRat(1, 1)
			bigfloat.PI(prec).Rat(a)
			b := complex.NewRational(a, big.NewRat(0, 1))
			return b
		case ruleprec:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					prec = uint(a.A.Num().Uint64())
					return a
				}
				node = node.next
			}
		case rulelog:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					x.SetRat(a)
					x.Log(x).Rat(a)
					return a
				}
				node = node.next
			}
		case rulesqrt:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					x.SetRat(a)
					x.Sqrt(x).Rat(a)
					return a
				}
				node = node.next
			}
		case rulesub:
			return c.Rulesub(node)
		}
		node = node.next
	}
	return nil
}

func (c *Calculator) Rulesub(node *node32) *complex.Rational {
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
