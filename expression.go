// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc

import (
	"math/big"

	"github.com/ALTree/bigfloat"
	complex "github.com/pointlander/c0mpl3x"
)

// Operation is a mathematical operation
type Operation uint

const (
	// OperationNoop is a noop
	OperationNoop Operation = iota
	// OperationAdd adds two numbers
	OperationAdd
	// OperationSubtract subtracts two numbers
	OperationSubtract
	// OperationMultiply multiplies two numbers
	OperationMultiply
	// OperationDivide divides two numbers
	OperationDivide
	// OperationModulus computes the modulus of two numbers
	OperationModulus
	// OperationExponentiation raises a number to a number
	OperationExponentiation
	// OperationNegate changes the sign of a number
	OperationNegate
	// OperationVariable is a variable
	OperationVariable
	// OperationImaginary is an imaginary number
	OperationImaginary
	// OperationNumber is a real number
	OperationNumber
	// OperationMatrix is a matrix
	OperationMatrix
	// OperationNaturalExponentiation raises the natural number to a power
	OperationNaturalExponentiation
	// OperationNatural is the constant e
	OperationNatural
	// OperationPI is the constant pi
	OperationPI
	// OperationSetPrec sets the precision
	OperationSetPrec
	// OperationSimplify simplify the expression
	OperationSimplify
	// OperationDerivative compute the symbolic derivative
	OperationDerivative
	// OperationNaturalLogarithm os the natural logarithm
	OperationNaturalLogarithm
	// OperationSquareRoot computes the square root of a number
	OperationSquareRoot
	// OperationCosine computes the cosine of a number
	OperationCosine
	// OperationSine computes the sine of a number
	OperationSine
	// OperationTangent computes the tangent of a number
	OperationTangent
	// OperationNotation is E notation operation
	OperationNotation
)

// Node is a node in an expression binary tree
type Node struct {
	Operation   Operation
	Value       string
	Matrix      *complex.Matrix
	Left, Right *Node
}

// Equals test if value is equal to x
func (n *Node) Equals(x int64) bool {
	if n.Operation == OperationNotation {
		a := big.NewInt(0)
		a.SetString(n.Left.Value, 10)
		b := big.NewInt(10)
		c := big.NewInt(0)
		c.SetString(n.Right.Value, 10)
		b.Exp(b, c, nil)
		a.Mul(a, b)
		return a.Cmp(big.NewInt(x)) == 0
	}
	value := big.NewInt(0)
	value.SetString(n.Value, 10)
	return value.Cmp(big.NewInt(x)) == 0
}

// String returns the string form of the equation
func (n *Node) String() string {
	var process func(n *Node) string
	process = func(n *Node) string {
		if n == nil {
			return ""
		}
		switch n.Operation {
		case OperationNoop:
			return "(" + process(n.Left) + "???" + process(n.Right) + ")"
		case OperationAdd:
			return "(" + process(n.Left) + " + " + process(n.Right) + ")"
		case OperationSubtract:
			return "(" + process(n.Left) + " - " + process(n.Right) + ")"
		case OperationMultiply:
			return "(" + process(n.Left) + " * " + process(n.Right) + ")"
		case OperationDivide:
			return "(" + process(n.Left) + " / " + process(n.Right) + ")"
		case OperationModulus:
			return "(" + process(n.Left) + " % " + process(n.Right) + ")"
		case OperationExponentiation:
			return "(" + process(n.Left) + "^" + process(n.Right) + ")"
		case OperationNegate:
			return "-(" + process(n.Left) + ")"
		case OperationVariable:
			return n.Value
		case OperationImaginary:
			return n.Value + "i"
		case OperationNumber:
			return n.Value
		case OperationNotation:
			if n.Left.Operation == OperationImaginary {
				return n.Left.Value + "e" + process(n.Right) + "i"
			}
			return process(n.Left) + "e" + process(n.Right)
		case OperationNaturalExponentiation:
			return "(e^" + process(n.Left) + ")"
		case OperationNatural:
			return "e"
		case OperationPI:
			return "pi"
		case OperationNaturalLogarithm:
			return "log(" + process(n.Left) + ")"
		case OperationSquareRoot:
			return "sqrt(" + process(n.Left) + ")"
		case OperationCosine:
			return "cos(" + process(n.Left) + ")"
		case OperationSine:
			return "sin(" + process(n.Left) + ")"
		case OperationTangent:
			return "tan(" + process(n.Left) + ")"
		}
		return ""
	}
	return process(n)
}

// Eval evaluates an expression
func (n *Node) Eval() (*Node, *complex.Matrix) {
	var expression *Node
	var process func(n *Node) *complex.Matrix
	process = func(n *Node) *complex.Matrix {
		if n == nil {
			ra := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a := complex.NewMatrix(prec)
			a.Values = [][]complex.Rational{[]complex.Rational{*ra}}
			return &a
		}
		switch n.Operation {
		case OperationNoop:
			ra := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a := complex.NewMatrix(prec)
			a.Values = [][]complex.Rational{[]complex.Rational{*ra}}
			return &a
		case OperationAdd:
			a := complex.NewMatrix(prec)
			a.Add(process(n.Left), process(n.Right))
			return &a
		case OperationSubtract:
			a := complex.NewMatrix(prec)
			a.Sub(process(n.Left), process(n.Right))
			return &a
		case OperationMultiply:
			a := complex.NewMatrix(prec)
			a.Mul(process(n.Left), process(n.Right))
			return &a
		case OperationDivide:
			a := complex.NewMatrix(prec)
			a.Div(process(n.Left), process(n.Right))
			return &a
		case OperationModulus:
			a := complex.NewMatrix(prec)
			a.Values = [][]complex.Rational{{*complex.NewRational(big.NewRat(0, 0), big.NewRat(0, 1))}}
			left := process(n.Left)
			right := process(n.Right)
			if left.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 && right.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 {
				a.Values[0][0].A.Num().Mod(left.Values[0][0].A.Num(), right.Values[0][0].A.Num())
			}
			return &a
		case OperationExponentiation:
			a := complex.NewMatrix(prec)
			right := process(n.Right).Values[0][0]
			a.Pow(process(n.Left), &right)
			return &a
		case OperationNegate:
			a := complex.NewMatrix(prec)
			a.Neg(process(n.Left))
			return &a
		case OperationMatrix:
			return n.Matrix
		case OperationVariable:
			ra := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a := complex.NewMatrix(prec)
			a.Values = [][]complex.Rational{[]complex.Rational{*ra}}
			return &a
		case OperationImaginary:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a.A.SetString(n.Value)

			a.A, a.B = a.B, a.A
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return &b
		case OperationNumber:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a.A.SetString(n.Value)

			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return &b
		case OperationNotation:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			left := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			left.A.SetString(n.Left.Value)
			right := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			right.A.SetString(n.Right.Value)
			c := complex.NewRational(big.NewRat(10, 1), big.NewRat(0, 1))
			x := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
			x.SetRat(c)
			y := complex.NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
			y.SetRat(right)
			x.Pow(x, y).Rat(right)
			a.Mul(left, right)

			if n.Left.Operation == OperationImaginary {
				a.A, a.B = a.B, a.A
			}
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			return &b
		case OperationNaturalExponentiation:
			a := complex.NewMatrix(prec)
			a.Exp(process(n.Left))
			return &a
		case OperationNatural:
			a := complex.NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			b := complex.NewMatrix(prec)
			b.Values = [][]complex.Rational{[]complex.Rational{*a}}
			b.Exp(&b)
			return &b
		case OperationPI:
			a := big.NewRat(1, 1)
			bigfloat.PI(prec).Rat(a)
			b := complex.NewRational(a, big.NewRat(0, 1))
			c := complex.NewMatrix(prec)
			c.Values = [][]complex.Rational{[]complex.Rational{*b}}
			return &c
		case OperationSetPrec:
			a := process(n.Left)
			prec = uint(a.Values[0][0].A.Num().Uint64())
			return a
		case OperationDerivative:
			derivative := n.Left.Derivative()
			if derivative != nil {
				derivative = derivative.Simplify()
			}
			expression = derivative
			return nil
		case OperationSimplify:
			derivative := n.Left.Simplify()
			expression = derivative
			return nil
		case OperationNaturalLogarithm:
			a := complex.NewMatrix(prec)
			a.Log(process(n.Left))
			return &a
		case OperationSquareRoot:
			a := complex.NewMatrix(prec)
			a.Sqrt(process(n.Left))
			return &a
		case OperationCosine:
			a := complex.NewMatrix(prec)
			a.Cos(process(n.Left))
			return &a
		case OperationSine:
			a := complex.NewMatrix(prec)
			a.Sin(process(n.Left))
			return &a
		case OperationTangent:
			a := complex.NewMatrix(prec)
			a.Tan(process(n.Left))
			return &a
		}
		a := complex.NewMatrix(prec)
		return &a
	}
	return expression, process(n)
}

// Derivative takes the derivative of the equation
// https://www.cs.utexas.edu/users/novak/asg-symdif.html#:~:text=Introduction,numeric%20calculations%20based%20on%20formulas.
func (n *Node) Derivative() *Node {
	var process func(n *Node) *Node
	process = func(n *Node) *Node {
		if n == nil {
			return nil
		}
		switch n.Operation {
		case OperationNoop:
			return n
		case OperationAdd:
			a := &Node{
				Operation: OperationAdd,
				Left:      process(n.Left),
				Right:     process(n.Right),
			}
			return a
		case OperationSubtract:
			a := &Node{
				Operation: OperationSubtract,
				Left:      process(n.Left),
				Right:     process(n.Right),
			}
			return a
		case OperationMultiply:
			left := &Node{
				Operation: OperationMultiply,
				Left:      n.Left,
				Right:     process(n.Right),
			}
			right := &Node{
				Operation: OperationMultiply,
				Left:      n.Right,
				Right:     process(n.Left),
			}
			a := &Node{
				Operation: OperationAdd,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationDivide:
			left := &Node{
				Operation: OperationMultiply,
				Left:      n.Right,
				Right:     process(n.Left),
			}
			right := &Node{
				Operation: OperationMultiply,
				Left:      n.Left,
				Right:     process(n.Right),
			}
			difference := &Node{
				Operation: OperationSubtract,
				Left:      left,
				Right:     right,
			}
			square := &Node{
				Operation: OperationExponentiation,
				Left:      n.Right,
				Right: &Node{
					Operation: OperationNumber,
					Value:     "2",
				},
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      difference,
				Right:     square,
			}
			return a
		case OperationModulus:
			return n
		case OperationExponentiation:
			one := &Node{
				Operation: OperationNumber,
				Value:     "1",
			}
			subtract := &Node{
				Operation: OperationSubtract,
				Left:      n.Right,
				Right:     one,
			}
			exp := &Node{
				Operation: OperationExponentiation,
				Left:      n.Left,
				Right:     subtract,
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      n.Right,
				Right:     exp,
			}
			a = &Node{
				Operation: OperationMultiply,
				Left:      a,
				Right:     process(n.Left),
			}
			return a
		case OperationNegate:
			a := &Node{
				Operation: OperationNegate,
				Left:      process(n.Left),
			}
			return a
		case OperationVariable:
			a := &Node{
				Operation: OperationNumber,
				Value:     "1",
			}
			return a
		case OperationImaginary:
			a := &Node{
				Operation: OperationNumber,
				Value:     "0",
			}
			return a
		case OperationNumber:
			a := &Node{
				Operation: OperationNumber,
				Value:     "0",
			}
			return a
		case OperationNotation:
			a := &Node{
				Operation: OperationNumber,
				Value:     "0",
			}
			return a
		case OperationNaturalExponentiation:
			a := &Node{
				Operation: OperationMultiply,
				Left:      n,
				Right:     process(n.Left),
			}
			return a
		case OperationNatural:
			a := &Node{
				Operation: OperationNumber,
				Value:     "0",
			}
			return a
		case OperationPI:
			a := &Node{
				Operation: OperationNumber,
				Value:     "0",
			}
			return a
		case OperationNaturalLogarithm:
			a := &Node{
				Operation: OperationDivide,
				Left:      process(n.Left),
				Right:     n.Left,
			}
			return a
		case OperationSquareRoot:
			value2 := &Node{
				Operation: OperationNumber,
				Value:     "2",
			}
			multiply := &Node{
				Operation: OperationMultiply,
				Left:      value2,
				Right:     n,
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      process(n.Left),
				Right:     multiply,
			}
			return a
		case OperationCosine:
			sin := &Node{
				Operation: OperationSine,
				Left:      n.Left,
			}
			multiply := &Node{
				Operation: OperationMultiply,
				Left:      sin,
				Right:     process(n.Left),
			}
			a := &Node{
				Operation: OperationNegate,
				Left:      multiply,
			}
			return a
		case OperationSine:
			cos := &Node{
				Operation: OperationCosine,
				Left:      n.Left,
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      cos,
				Right:     process(n.Left),
			}
			return a
		case OperationTangent:
			value1 := &Node{
				Operation: OperationNumber,
				Value:     "1",
			}
			value2 := &Node{
				Operation: OperationNumber,
				Value:     "2",
			}
			exp := &Node{
				Operation: OperationExponentiation,
				Left:      n,
				Right:     value2,
			}
			add := &Node{
				Operation: OperationAdd,
				Left:      value1,
				Right:     exp,
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      add,
				Right:     process(n.Left),
			}
			return a
		}
		return nil
	}
	return process(n)
}

var numeric = map[Operation]bool{
	OperationNumber:    true,
	OperationImaginary: true,
	OperationNotation:  true,
}

func isNumeric(operation Operation) bool {
	return numeric[operation]
}

// Simplify simplifies an expression
func (n *Node) Simplify() *Node {
	var process func(n *Node) *Node
	process = func(n *Node) *Node {
		if n == nil {
			return nil
		}
		switch n.Operation {
		case OperationNoop:
			return n
		case OperationAdd:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				return right
			} else if isNumeric(right.Operation) && right.Equals(0) {
				return left
			}
			a := &Node{
				Operation: OperationAdd,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationSubtract:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNegate,
					Left:      right,
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				return left
			}
			a := &Node{
				Operation: OperationSubtract,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationMultiply:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "0",
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "0",
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				return right
			} else if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationDivide:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "0",
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "+Inf",
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationModulus:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationModulus,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationExponentiation:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "0",
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "1",
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "1",
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationExponentiation,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationNegate:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "0",
				}
				return a
			}
			a := &Node{
				Operation: OperationNegate,
				Left:      left,
			}
			return a
		case OperationVariable:
			return n
		case OperationImaginary:
			return n
		case OperationNumber:
			return n
		case OperationNotation:
			return n
		case OperationNaturalExponentiation:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "1",
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				a := &Node{
					Operation: OperationVariable,
					Value:     "e",
				}
				return a
			}
			a := &Node{
				Operation: OperationNaturalExponentiation,
				Left:      left,
			}
			return a
		case OperationNatural:
			return n
		case OperationPI:
			return n
		case OperationNaturalLogarithm:
			left := process(n.Left)
			if left.Operation == OperationNatural {
				return left
			}
			a := &Node{
				Operation: OperationNaturalLogarithm,
				Left:      left,
			}
			return a
		case OperationSquareRoot:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "0",
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				a := &Node{
					Operation: OperationNumber,
					Value:     "1",
				}
				return a
			}
			a := &Node{
				Operation: OperationSquareRoot,
				Left:      left,
			}
			return a
		case OperationCosine:
			a := &Node{
				Operation: OperationCosine,
				Left:      process(n.Left),
			}
			return a
		case OperationSine:
			a := &Node{
				Operation: OperationSine,
				Left:      process(n.Left),
			}
			return a
		case OperationTangent:
			a := &Node{
				Operation: OperationTangent,
				Left:      process(n.Left),
			}
			return a
		}
		return nil
	}
	return process(n)
}
