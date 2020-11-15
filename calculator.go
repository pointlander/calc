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
		case rulepi:
			// https://en.wikipedia.org/wiki/Gauss%E2%80%93Legendre_algorithm
			a, b, t, p :=
				big.NewFloat(1).SetPrec(1024), big.NewFloat(1).SetPrec(1024),
				big.NewFloat(1).SetPrec(1024), big.NewFloat(1).SetPrec(1024)
			b = b.Quo(b, bigfloat.Sqrt(big.NewFloat(2).SetPrec(1024)))
			t = t.Quo(t, big.NewFloat(4).SetPrec(1024))
			for {
				an := big.NewFloat(0).SetPrec(1024)
				an.Add(a, b)
				an.Quo(an, big.NewFloat(2).SetPrec(1024))

				bn := big.NewFloat(0).SetPrec(1024)
				bn.Mul(a, b)
				bn = bigfloat.Sqrt(bn)

				diff := big.NewFloat(0).SetPrec(1024)
				diff.Sub(a, b)
				diffn := big.NewFloat(0).SetPrec(1024)
				diffn.Sub(an, bn)
				if diff.Cmp(diffn) == 0 {
					break
				}

				tn := big.NewFloat(0).SetPrec(1024)
				tn.Sub(a, an)
				tn.Mul(tn, tn)
				tn.Mul(tn, p)
				tn.Sub(t, tn)

				pn := big.NewFloat(0).SetPrec(1024)
				pn.Mul(p, big.NewFloat(2).SetPrec(1024))

				a, b, t, p = an, bn, tn, pn
			}
			pi := big.NewFloat(0).SetPrec(1024)
			pi.Add(a, b)
			pi.Mul(pi, pi)
			pi.Quo(pi, t.Mul(t, big.NewFloat(4).SetPrec(1024)))
			x := big.NewRat(1, 1)
			x.SetString(pi.Text('f', 1024))
			return x
		case ruleexp:
			node := node.up
			for node != nil {
				if node.pegRule == rulevalue {
					a := c.Rulevalue(node)
					x, _, _ := big.ParseFloat(a.FloatString(1024), 10, 1024, big.ToNearestEven)
					exp := bigfloat.Exp(x)
					a.SetString(exp.String())
					return a
				}
				node = node.next
			}
		case rulelog:
			node := node.up
			for node != nil {
				if node.pegRule == rulevalue {
					a := c.Rulevalue(node)
					x, _, _ := big.ParseFloat(a.FloatString(1024), 10, 1024, big.ToNearestEven)
					log := bigfloat.Log(x)
					a.SetString(log.String())
					return a
				}
				node = node.next
			}
		case rulesqrt:
			node := node.up
			for node != nil {
				if node.pegRule == rulevalue {
					a := c.Rulevalue(node)
					x, _, _ := big.ParseFloat(a.FloatString(1024), 10, 1024, big.ToNearestEven)
					sqrt := bigfloat.Sqrt(x)
					a.SetString(sqrt.String())
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
