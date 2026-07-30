// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc

import (
	"math/big"
	"strings"

	complex "github.com/pointlander/c0mpl3x"
)

//go:generate go tool peg calculator.peg

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
func (c *Calculator[T]) Eval() Value {
	value := c.Convert(c.AST().up)
	e, m := value.Expression.Eval()
	return Value{
		Matrix:     m,
		Expression: e,
	}
}

// Convert converts to an expression
func (c *Calculator[T]) Convert(n *node[T]) Value {
	var (
		convert      func(n *node[T]) (a *Node)
		convertValue func(n *node[T]) (a *Node)
	)
	convertValue = func(node *node[T]) (a *Node) {
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
			case rulematrix:
				node = node.up
				x := complex.NewMatrix(prec)
				x.Values = make([][]complex.Rational, 1)
				for node != nil {
					switch node.pegRule {
					case rulee1:
						a, end := convert(node), len(x.Values)-1
						if a.Value != "" {
							ra := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
							ra.A.SetString(a.Value)
							x.Values[end] = append(x.Values[end], *ra)
							break
						}
						panic("matrix within matrix not allowed")
					case rulerow:
						x.Values = append(x.Values, make([]complex.Rational, 0, 8))
					}
					node = node.next
				}
				a = &Node{
					Operation: OperationMatrix,
					Matrix:    &x,
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
							Left:      convertValue(node),
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
			case ruleprec:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationSetPrec,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulesimplify:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationSimplify,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
			case rulederivative:
				node := node.up
				for node != nil {
					if node.pegRule == rulee1 {
						a = &Node{
							Operation: OperationDerivative,
							Left:      convert(node),
						}
						return a
					}
					node = node.next
				}
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
			if node == nil {
				break
			}
			node = node.next
		}
		return a
	}
	convert = func(node *node[T]) (a *Node) {
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
		Expression: convert(n),
	}
}
