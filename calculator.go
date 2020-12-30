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

// Eval evaluates the expression
func (c *Calculator) Eval() *complex.Matrix {
	return c.Rulee(c.AST())
}

// Rulee is a root expresion
func (c *Calculator) Rulee(node *node32) *complex.Matrix {
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

// Rulee1 deals with addation or subtraction
func (c *Calculator) Rulee1(node *node32) *complex.Matrix {
	node = node.up
	var a *complex.Matrix
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

// Rulee2 deals with multiplication, division, or modulus
func (c *Calculator) Rulee2(node *node32) *complex.Matrix {
	node = node.up
	var a *complex.Matrix
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
			if a.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 && b.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 {
				a.Values[0][0].A.Num().Mod(a.Values[0][0].A.Num(), b.Values[0][0].A.Num())
			}
		}
		node = node.next
	}
	return a
}

// Rulee3 deals with exponentiation
func (c *Calculator) Rulee3(node *node32) *complex.Matrix {
	node = node.up
	var a *complex.Matrix
	for node != nil {
		switch node.pegRule {
		case rulee4:
			a = c.Rulee4(node)
		case ruleexponentiation:
			node = node.next
			b := c.Rulee4(node)
			a.Pow(a, &b.Values[0][0])
		}
		node = node.next
	}
	return a
}

// Rulee4 negates a number
func (c *Calculator) Rulee4(node *node32) *complex.Matrix {
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

// Rulevalue evaluates the value
func (c *Calculator) Rulevalue(node *node32) *complex.Matrix {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulematrix:
			return c.Rulematrix(node)
		case ruleimaginary:
			a := complex.NewRational(big.NewRat(0, 1), big.NewRat(1, 1))
			a.B.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return &b
		case rulenumber:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a.A.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return &b
		case ruleexp1:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Exp(a)
					return a
				}
				node = node.next
			}
		case ruleexp2:
			node := node.up
			for node != nil {
				if node.pegRule == rulevalue {
					a := c.Rulevalue(node)
					a.Exp(a)
					return a
				}
				node = node.next
			}
		case rulepi:
			a := big.NewRat(1, 1)
			bigfloat.PI(prec).Rat(a)
			b := complex.NewRational(a, big.NewRat(0, 1))
			c := complex.NewMatrix(prec)
			c.Values = [][]complex.Rational{[]complex.Rational{*b}}
			return &c
		case ruleprec:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					prec = uint(a.Values[0][0].A.Num().Uint64())
					return a
				}
				node = node.next
			}
		case rulederivative:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					c.Rulederivative(node)
					return &complex.Matrix{}
				}
				node = node.next
			}
		case rulelog:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Log(a)
					return a
				}
				node = node.next
			}
		case rulesqrt:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Sqrt(a)
					return a
				}
				node = node.next
			}
		case rulecos:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Cos(a)
					return a
				}
				node = node.next
			}
		case rulesin:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Sin(a)
					return a
				}
				node = node.next
			}
		case ruletan:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Tan(a)
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

// Rulematrix computes the matrix
func (c *Calculator) Rulematrix(node *node32) *complex.Matrix {
	node = node.up
	x := complex.NewMatrix(prec)
	x.Values = make([][]complex.Rational, 1)
	for node != nil {
		switch node.pegRule {
		case rulee1:
			a, end := c.Rulee1(node), len(x.Values)-1
			if len(a.Values) == 1 && len(a.Values[0]) == 1 {
				x.Values[end] = append(x.Values[end], a.Values[0][0])
				break
			}
			panic("matrix within matrix not allowed")
		case rulerow:
			x.Values = append(x.Values, make([]complex.Rational, 0, 8))
		}
		node = node.next
	}
	return &x
}

// Rulederivative computes the symbolic derivative of a number
func (c *Calculator) Rulederivative(node *node32) *complex.Matrix {
	var (
		convert      func(node *node32) (a *Node)
		convertValue func(node *node32) (a *Node)
	)
	convertValue = func(node *node32) (a *Node) {
		node = node.up
		for node != nil {
			switch node.pegRule {
			case rulevalue:
				a = convertValue(node)
			case ruleminus:
				a = &Node{
					Operation: OperationNegate,
					Left:      convertValue(node),
				}
			case ruleimaginary:
				a = &Node{
					Operation: OperationImaginary,
					Value:     strings.TrimSpace(string(c.buffer[node.begin:node.end])),
				}
			case rulenumber:
				a = &Node{
					Operation: OperationNumber,
					Value:     strings.TrimSpace(string(c.buffer[node.begin:node.end])),
				}
			case ruleexp1:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationNaturalExponentiation,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case ruleexp2:
				node := node.up
				for node != nil {
					if node.pegRule == rulevalue {
						a = &Node{
							Operation: OperationNaturalExponentiation,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulepi:
				a = &Node{
					Operation: OperationPI,
				}
				return a
			case rulelog:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationNaturalLogarithm,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulesqrt:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationSquareRoot,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulecos:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationCosine,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulesin:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationSine,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case ruletan:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationTangent,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulesub:
				return convert(node)
			}
			node = node.next
		}
		return a
	}
	convert = func(node *node32) (a *Node) {
		node = node.up
		for node != nil {
			switch node.pegRule {
			case rulee2, rulee3:
				a = convert(node)
			case ruleadd:
				node = node.next
				a = &Node{
					Operation: OperationAdd,
					Left:      a,
					Right:     convert(node),
				}
			case ruleminus:
				node = node.next
				a = &Node{
					Operation: OperationSubtract,
					Left:      a,
					Right:     convert(node),
				}
			case rulemultiply:
				node = node.next
				a = &Node{
					Operation: OperationMultiply,
					Left:      a,
					Right:     convert(node),
				}
			case ruledivide:
				node = node.next
				a = &Node{
					Operation: OperationDivide,
					Left:      a,
					Right:     convert(node),
				}
			case rulemodulus:
				node = node.next
				a = &Node{
					Operation: OperationModulus,
					Left:      a,
					Right:     convert(node),
				}
			case rulee4:
				a = convertValue(node)
			case ruleexponentiation:
				node = node.next
				a = &Node{
					Operation: OperationExponentiation,
					Left:      a,
					Right:     convertValue(node),
				}
			}
			node = node.next
		}
		return a
	}
	convert(node)
	return nil
}

// Rulesub computes the subexpression
func (c *Calculator) Rulesub(node *node32) *complex.Matrix {
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
