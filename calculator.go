// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc

import (
	"math/big"
	"strings"

	"github.com/ALTree/bigfloat"
	complex "github.com/pointlander/c0mpl3x"
)

var prec uint = 1024

// ValueType is a value type
type ValueType int

const (
	// ValueTypeMatrix is a matrix value type
	ValueTypeMatrix ValueType = iota
	// ValueTypeExpression is an expression value type
	ValueTypeExpression
)

// Value is a value
type Value struct {
	ValueType  ValueType
	Matrix     *complex.Matrix
	Expression *Node
}

// Eval evaluates the expression
func (c *Calculator) Eval() Value {
	return c.Rulee(c.AST())
}

// Rulee is a root expresion
func (c *Calculator) Rulee(node *node32) Value {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulee1:
			return c.Rulee1(node)
		}
		node = node.next
	}
	return Value{}
}

// Rulee1 deals with addation or subtraction
func (c *Calculator) Rulee1(node *node32) Value {
	node = node.up
	var a Value
	for node != nil {
		switch node.pegRule {
		case rulee2:
			a = c.Rulee2(node)
		case ruleadd:
			node = node.next
			b := c.Rulee2(node)
			a.Matrix.Add(a.Matrix, b.Matrix)
		case ruleminus:
			node = node.next
			b := c.Rulee2(node)
			a.Matrix.Sub(a.Matrix, b.Matrix)
		}
		node = node.next
	}
	return a
}

// Rulee2 deals with multiplication, division, or modulus
func (c *Calculator) Rulee2(node *node32) Value {
	node = node.up
	var a Value
	for node != nil {
		switch node.pegRule {
		case rulee3:
			a = c.Rulee3(node)
		case rulemultiply:
			node = node.next
			b := c.Rulee3(node)
			a.Matrix.Mul(a.Matrix, b.Matrix)
		case ruledivide:
			node = node.next
			b := c.Rulee3(node)
			a.Matrix.Div(a.Matrix, b.Matrix)
		case rulemodulus:
			node = node.next
			b := c.Rulee3(node)
			if a.Matrix.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 && b.Matrix.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 {
				a.Matrix.Values[0][0].A.Num().Mod(a.Matrix.Values[0][0].A.Num(), b.Matrix.Values[0][0].A.Num())
			}
		}
		node = node.next
	}
	return a
}

// Rulee3 deals with exponentiation
func (c *Calculator) Rulee3(node *node32) Value {
	node = node.up
	var a Value
	for node != nil {
		switch node.pegRule {
		case rulee4:
			a = c.Rulee4(node)
		case ruleexponentiation:
			node = node.next
			b := c.Rulee4(node)
			a.Matrix.Pow(a.Matrix, &b.Matrix.Values[0][0])
		}
		node = node.next
	}
	return a
}

// Rulee4 negates a number
func (c *Calculator) Rulee4(node *node32) Value {
	node = node.up
	minus := false
	for node != nil {
		switch node.pegRule {
		case rulevalue:
			a := c.Rulevalue(node)
			if minus {
				a.Matrix.Neg(a.Matrix)
			}
			return a
		case ruleminus:
			minus = true
		}
		node = node.next
	}
	return Value{}
}

// Rulevalue evaluates the value
func (c *Calculator) Rulevalue(node *node32) Value {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulematrix:
			return c.Rulematrix(node)
		case ruleimaginary:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			node := node.up
			for node != nil {
				switch node.pegRule {
				case ruledecimal:
					a.A.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
				case rulenotation:
					b := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
					b.A.SetString(strings.TrimSpace(string(c.buffer[node.up.begin:node.up.end])))
					c := complex.NewRational(big.NewRat(10, 1), big.NewRat(0, 1))
					x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					x.SetRat(c)
					y := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					y.SetRat(b)
					x.Pow(x, y).Rat(b)
					a.Mul(a, b)
				}
				node = node.next
			}

			a.A, a.B = a.B, a.A
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return Value{
				ValueType: ValueTypeMatrix,
				Matrix:    &b,
			}
		case rulenumber:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			node := node.up
			for node != nil {
				switch node.pegRule {
				case ruledecimal:
					a.A.SetString(strings.TrimSpace(string(c.buffer[node.begin:node.end])))
				case rulenotation:
					b := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
					b.A.SetString(strings.TrimSpace(string(c.buffer[node.up.begin:node.up.end])))
					c := complex.NewRational(big.NewRat(10, 1), big.NewRat(0, 1))
					x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					x.SetRat(c)
					y := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
					y.SetRat(b)
					x.Pow(x, y).Rat(b)
					a.Mul(a, b)
				}
				node = node.next
			}

			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return Value{
				ValueType: ValueTypeMatrix,
				Matrix:    &b,
			}
		case ruleexp1:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Matrix.Exp(a.Matrix)
					return a
				}
				node = node.next
			}
		case ruleexp2:
			node := node.up
			for node != nil {
				if node.pegRule == rulevalue {
					a := c.Rulevalue(node)
					a.Matrix.Exp(a.Matrix)
					return a
				}
				node = node.next
			}
		case rulenatural:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return Value{
				ValueType: ValueTypeMatrix,
				Matrix:    b.Exp(&b),
			}
		case rulepi:
			a := big.NewRat(1, 1)
			bigfloat.PI(prec).Rat(a)
			b := complex.NewRational(a, big.NewRat(0, 1))
			c := complex.NewMatrix(prec)
			c.Values = [][]complex.Rational{[]complex.Rational{*b}}
			return Value{
				ValueType: ValueTypeMatrix,
				Matrix:    &c,
			}
		case ruleprec:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					prec = uint(a.Matrix.Values[0][0].A.Num().Uint64())
					return a
				}
				node = node.next
			}
		case rulesimplify:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					return c.Rulesimplify(node)
				}
				node = node.next
			}
		case rulederivative:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					return c.Rulederivative(node)
				}
				node = node.next
			}
		case rulelog:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Matrix.Log(a.Matrix)
					return a
				}
				node = node.next
			}
		case rulesqrt:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Matrix.Sqrt(a.Matrix)
					return a
				}
				node = node.next
			}
		case rulecos:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Matrix.Cos(a.Matrix)
					return a
				}
				node = node.next
			}
		case rulesin:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Matrix.Sin(a.Matrix)
					return a
				}
				node = node.next
			}
		case ruletan:
			node := node.up
			for node != nil {
				if node.pegRule == rulee1 {
					a := c.Rulee1(node)
					a.Matrix.Tan(a.Matrix)
					return a
				}
				node = node.next
			}
		case rulesub:
			return c.Rulesub(node)
		}
		node = node.next
	}
	return Value{}
}

// Rulematrix computes the matrix
func (c *Calculator) Rulematrix(node *node32) Value {
	node = node.up
	x := complex.NewMatrix(prec)
	x.Values = make([][]complex.Rational, 1)
	for node != nil {
		switch node.pegRule {
		case rulee1:
			a, end := c.Rulee1(node), len(x.Values)-1
			if len(a.Matrix.Values) == 1 && len(a.Matrix.Values[0]) == 1 {
				x.Values[end] = append(x.Values[end], a.Matrix.Values[0][0])
				break
			}
			panic("matrix within matrix not allowed")
		case rulerow:
			x.Values = append(x.Values, make([]complex.Rational, 0, 8))
		}
		node = node.next
	}
	return Value{
		ValueType: ValueTypeMatrix,
		Matrix:    &x,
	}
}

// Convert converts to an expression
func (c *Calculator) Convert(node *node32) Value {
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
			case rulevariable:
				a = &Node{
					Operation: OperationVariable,
					Value:     strings.TrimSpace(string(c.buffer[node.begin:node.end])),
				}
			case ruleimaginary:
				node := node.up
				a = &Node{
					Operation: OperationImaginary,
				}
				for node != nil {
					switch node.pegRule {
					case ruledecimal:
						a.Value = strings.TrimSpace(string(c.buffer[node.begin:node.end]))
					case rulenotation:
						a.Left = &Node{
							Operation: OperationImaginary,
							Value:     a.Value,
						}
						a.Operation = OperationNotation
						a.Value = ""
						a.Right = &Node{
							Operation: OperationNumber,
							Value:     strings.TrimSpace(string(c.buffer[node.up.begin:node.up.end])),
						}
					}
					node = node.next
				}
				return a
			case rulenumber:
				node := node.up
				a = &Node{
					Operation: OperationNumber,
				}
				for node != nil {
					switch node.pegRule {
					case ruledecimal:
						a.Value = strings.TrimSpace(string(c.buffer[node.begin:node.end]))
					case rulenotation:
						a.Left = &Node{
							Operation: OperationNumber,
							Value:     a.Value,
						}
						a.Operation = OperationNotation
						a.Value = ""
						a.Right = &Node{
							Operation: OperationNumber,
							Value:     strings.TrimSpace(string(c.buffer[node.up.begin:node.up.end])),
						}
					}
					node = node.next
				}
				return a
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
			case rulenatural:
				a = &Node{
					Operation: OperationNatural,
				}
				return a
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
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						return convert(node)
					}
					node = node.next
				}
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
	return Value{
		ValueType:  ValueTypeExpression,
		Expression: convert(node),
	}
}

// Rulesimplify simplifies the expression
func (c *Calculator) Rulesimplify(node *node32) Value {
	expression := c.Convert(node).Expression
	if expression != nil {
		expression = expression.Simplify()
	}
	return Value{
		ValueType:  ValueTypeExpression,
		Expression: expression,
	}
}

// Rulederivative computes the symbolic derivative of a number
func (c *Calculator) Rulederivative(node *node32) Value {
	expression := c.Convert(node).Expression
	derivative := expression.Derivative()
	if derivative != nil {
		derivative = derivative.Simplify()
	}
	return Value{
		ValueType:  ValueTypeExpression,
		Expression: derivative,
	}
}

// Rulesub computes the subexpression
func (c *Calculator) Rulesub(node *node32) Value {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulee1:
			return c.Rulee1(node)
		}
		node = node.next
	}
	return Value{}
}
