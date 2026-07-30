// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc

import (
	"fmt"
	"math/big"

	"github.com/ALTree/bigfloat"
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
	// OperationIntegrate integrates the equation
	OperationIntegrate
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
	// OperationAbs computes the absolute value / modulus
	OperationAbs
	// OperationConj computes the complex conjugate
	OperationConj
	// OperationRe extracts the real part
	OperationRe
	// OperationIm extracts the imaginary part
	OperationIm
	// OperationNotation is E notation operation
	OperationNotation
)

// Node is a node in an expression binary tree
type Node struct {
	Operation   Operation
	Value       string
	Matrix      *Matrix
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
			if n.Value == "" || n.Value == "1" || n.Value == "+1" {
				return "i"
			}
			if n.Value == "-1" {
				return "-i"
			}
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
		case OperationAbs:
			return "abs(" + process(n.Left) + ")"
		case OperationConj:
			return "conj(" + process(n.Left) + ")"
		case OperationRe:
			return "re(" + process(n.Left) + ")"
		case OperationIm:
			return "im(" + process(n.Left) + ")"
		}
		return ""
	}
	return process(n)
}

// Eval evaluates an expression
func (n *Node) Eval() (*Node, *Matrix) {
	var expression *Node
	var process func(n *Node) *Matrix
	process = func(n *Node) *Matrix {
		if n == nil {
			ra := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a := NewMatrix(prec)
			a.Values = [][]Rational{[]Rational{*ra}}
			return &a
		}
		switch n.Operation {
		case OperationNoop:
			ra := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a := NewMatrix(prec)
			a.Values = [][]Rational{[]Rational{*ra}}
			return &a
		case OperationAdd:
			a := NewMatrix(prec)
			a.Add(process(n.Left), process(n.Right))
			return &a
		case OperationSubtract:
			a := NewMatrix(prec)
			a.Sub(process(n.Left), process(n.Right))
			return &a
		case OperationMultiply:
			a := NewMatrix(prec)
			a.Mul(process(n.Left), process(n.Right))
			return &a
		case OperationDivide:
			a := NewMatrix(prec)
			a.Div(process(n.Left), process(n.Right))
			return &a
		case OperationModulus:
			a := NewMatrix(prec)
			a.Values = [][]Rational{{*NewRational(big.NewRat(0, 0), big.NewRat(0, 1))}}
			left := process(n.Left)
			right := process(n.Right)
			if left.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 && right.Values[0][0].A.Denom().Cmp(big.NewInt(1)) == 0 {
				a.Values[0][0].A.Num().Mod(left.Values[0][0].A.Num(), right.Values[0][0].A.Num())
			}
			return &a
		case OperationExponentiation:
			a := NewMatrix(prec)
			right := process(n.Right).Values[0][0]
			a.Pow(process(n.Left), &right)
			return &a
		case OperationNegate:
			a := NewMatrix(prec)
			a.Neg(process(n.Left))
			return &a
		case OperationMatrix:
			return n.Matrix
		case OperationVariable:
			ra := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a := NewMatrix(prec)
			a.Values = [][]Rational{[]Rational{*ra}}
			return &a
		case OperationImaginary:
			a := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			// Bare "i" has empty Value and means 1i.
			if n.Value == "" {
				a.A.SetInt64(1)
			} else {
				a.A.SetString(n.Value)
			}
			a.A, a.B = a.B, a.A
			b := NewMatrix(prec)
			b.Values = [][]Rational{[]Rational{*a}}
			return &b
		case OperationNumber:
			a := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			a.A.SetString(n.Value)

			b := NewMatrix(prec)
			b.Values = [][]Rational{[]Rational{*a}}
			return &b
		case OperationNotation:
			a := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			left := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			left.A.SetString(n.Left.Value)
			right := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			right.A.SetString(n.Right.Value)
			c := NewRational(big.NewRat(10, 1), big.NewRat(0, 1))
			x := NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
			x.SetRat(c)
			y := NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
			y.SetRat(right)
			x.Pow(x, y).Rat(right)
			a.Mul(left, right)

			if n.Left.Operation == OperationImaginary {
				a.A, a.B = a.B, a.A
			}
			b := NewMatrix(prec)
			b.Values = [][]Rational{[]Rational{*a}}
			return &b
		case OperationNaturalExponentiation:
			a := NewMatrix(prec)
			a.Exp(process(n.Left))
			return &a
		case OperationNatural:
			a := NewRational(big.NewRat(1, 1), big.NewRat(0, 1))
			b := NewMatrix(prec)
			b.Values = [][]Rational{[]Rational{*a}}
			b.Exp(&b)
			return &b
		case OperationPI:
			a := big.NewRat(1, 1)
			bigfloat.PI(prec).Rat(a)
			b := NewRational(a, big.NewRat(0, 1))
			c := NewMatrix(prec)
			c.Values = [][]Rational{[]Rational{*b}}
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
		case OperationIntegrate:
			integral := n.Left.Integrate()
			if integral != nil {
				integral = integral.Simplify()
			}
			expression = integral
			return nil
		case OperationSimplify:
			expression = n.Left.Simplify()
			return nil
		case OperationNaturalLogarithm:
			a := NewMatrix(prec)
			a.Log(process(n.Left))
			return &a
		case OperationSquareRoot:
			a := NewMatrix(prec)
			a.Sqrt(process(n.Left))
			return &a
		case OperationCosine:
			a := NewMatrix(prec)
			a.Cos(process(n.Left))
			return &a
		case OperationSine:
			a := NewMatrix(prec)
			a.Sin(process(n.Left))
			return &a
		case OperationTangent:
			a := NewMatrix(prec)
			a.Tan(process(n.Left))
			return &a
		case OperationAbs:
			a := NewMatrix(prec)
			a.Abs(process(n.Left))
			return &a
		case OperationConj:
			a := NewMatrix(prec)
			a.Conj(process(n.Left))
			return &a
		case OperationRe:
			a := NewMatrix(prec)
			a.Re(process(n.Left))
			return &a
		case OperationIm:
			a := NewMatrix(prec)
			a.Im(process(n.Left))
			return &a
		}
		a := NewMatrix(prec)
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
		case OperationImaginary, OperationNumber, OperationNotation,
			OperationNatural, OperationPI:
			// Complex/real constants differentiate to 0
			return &Node{Operation: OperationNumber, Value: "0"}
		case OperationAbs, OperationConj, OperationRe, OperationIm:
			// Constant complex helpers → 0; otherwise leave as-is (not fully expanded)
			if isConstantExpr(n.Left) {
				return &Node{Operation: OperationNumber, Value: "0"}
			}
			return &Node{Operation: n.Operation, Left: process(n.Left)}
		case OperationNaturalExponentiation:
			a := &Node{
				Operation: OperationMultiply,
				Left:      n,
				Right:     process(n.Left),
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

// isConstantExpr reports whether the expression contains no free variables.
// Imaginary literals, re/im/conj/abs of constants, e, and pi count as constant.
func isConstantExpr(n *Node) bool {
	if n == nil {
		return true
	}
	switch n.Operation {
	case OperationVariable:
		return false
	case OperationDerivative, OperationIntegrate, OperationSimplify, OperationSetPrec:
		return false
	case OperationAbs, OperationConj, OperationRe, OperationIm,
		OperationNegate, OperationNaturalLogarithm, OperationSquareRoot,
		OperationCosine, OperationSine, OperationTangent, OperationNaturalExponentiation:
		return isConstantExpr(n.Left)
	}
	return isConstantExpr(n.Left) && isConstantExpr(n.Right)
}

// findVariable returns the first variable name in the expression, or "x".
func findVariable(n *Node) string {
	var found string
	var walk func(*Node)
	walk = func(m *Node) {
		if m == nil || found != "" {
			return
		}
		if m.Operation == OperationVariable {
			found = m.Value
			return
		}
		walk(m.Left)
		walk(m.Right)
	}
	walk(n)
	if found != "" {
		return found
	}
	return "x"
}

// partsPriority returns a LIATE priority for integration by parts.
// Lower values are preferred as u (Log, Inverse trig, Algebraic, Trig, Exp).
func partsPriority(n *Node) int {
	if n == nil {
		return 99
	}
	switch n.Operation {
	case OperationNaturalLogarithm:
		return 1
	case OperationVariable:
		return 3
	case OperationExponentiation:
		// log(x)^n counts as logarithmic for LIATE
		if n.Left != nil && n.Left.Operation == OperationNaturalLogarithm {
			return 1
		}
		if n.Left != nil && n.Left.Operation == OperationVariable && isConstantExpr(n.Right) {
			return 3 // algebraic x^n
		}
		// sin(x)^n, cos(x)^n — treat as trig
		if n.Left != nil && (n.Left.Operation == OperationSine ||
			n.Left.Operation == OperationCosine || n.Left.Operation == OperationTangent) {
			return 4
		}
		if isConstantExpr(n.Left) {
			return 5 // exponential a^x
		}
		return 3
	case OperationSquareRoot:
		return 3
	case OperationSine, OperationCosine, OperationTangent:
		return 4
	case OperationNaturalExponentiation:
		return 5
	case OperationNegate:
		return partsPriority(n.Left)
	case OperationMultiply:
		l, r := partsPriority(n.Left), partsPriority(n.Right)
		if l < r {
			return l
		}
		return r
	default:
		if isConstantExpr(n) {
			return 0
		}
		return 50
	}
}

// productKey returns a order-independent key for a product (or the expression string).
func productKey(n *Node) string {
	if n == nil {
		return ""
	}
	if n.Operation == OperationMultiply && n.Left != nil && n.Right != nil {
		a, b := n.Left.String(), n.Right.String()
		if a > b {
			a, b = b, a
		}
		return a + "*" + b
	}
	if n.Operation == OperationNegate && n.Left != nil {
		return "-(" + productKey(n.Left) + ")"
	}
	return n.String()
}

// containsNode reports whether root contains target by pointer identity.
func containsNode(root, target *Node) bool {
	if root == nil {
		return false
	}
	if root == target {
		return true
	}
	return containsNode(root.Left, target) || containsNode(root.Right, target)
}

func nodeNumber(v string) *Node {
	return &Node{Operation: OperationNumber, Value: v}
}

func nodeZero() *Node { return nodeNumber("0") }
func nodeOne() *Node  { return nodeNumber("1") }

func nodeNegate(n *Node) *Node {
	if n == nil {
		return nil
	}
	if n.Operation == OperationNegate {
		return n.Left // --x = x
	}
	if isNumeric(n.Operation) && n.Equals(0) {
		return nodeZero()
	}
	return &Node{Operation: OperationNegate, Left: n}
}

func numericInt(n *Node) (int64, bool) {
	if n == nil {
		return 0, false
	}
	if isNumeric(n.Operation) {
		v := big.NewInt(0)
		if _, ok := v.SetString(n.Value, 10); !ok {
			return 0, false
		}
		if !v.IsInt64() {
			return 0, false
		}
		return v.Int64(), true
	}
	switch n.Operation {
	case OperationNegate:
		if v, ok := numericInt(n.Left); ok {
			return -v, true
		}
	case OperationAdd:
		a, ok1 := numericInt(n.Left)
		b, ok2 := numericInt(n.Right)
		if ok1 && ok2 {
			return a + b, true
		}
	case OperationSubtract:
		a, ok1 := numericInt(n.Left)
		b, ok2 := numericInt(n.Right)
		if ok1 && ok2 {
			return a - b, true
		}
	case OperationMultiply:
		a, ok1 := numericInt(n.Left)
		b, ok2 := numericInt(n.Right)
		if ok1 && ok2 {
			return a * b, true
		}
	case OperationDivide:
		a, ok1 := numericInt(n.Left)
		b, ok2 := numericInt(n.Right)
		if ok1 && ok2 && b != 0 && a%b == 0 {
			return a / b, true
		}
	}
	return 0, false
}

func nodeAdd(a, b *Node) *Node {
	if a == nil || b == nil {
		return nil
	}
	if isNumeric(a.Operation) && a.Equals(0) {
		return b
	}
	if isNumeric(b.Operation) && b.Equals(0) {
		return a
	}
	// a + (-b) = a - b
	if b.Operation == OperationNegate {
		return nodeSub(a, b.Left)
	}
	// (-a) + b = b - a
	if a.Operation == OperationNegate {
		return nodeSub(b, a.Left)
	}
	if x, ok := numericInt(a); ok {
		if y, ok := numericInt(b); ok {
			return nodeNumber(fmt.Sprintf("%d", x+y))
		}
	}
	return &Node{Operation: OperationAdd, Left: a, Right: b}
}

func nodeSub(a, b *Node) *Node {
	if a == nil || b == nil {
		return nil
	}
	if isNumeric(b.Operation) && b.Equals(0) {
		return a
	}
	// a - (-b) = a + b
	if b.Operation == OperationNegate {
		return nodeAdd(a, b.Left)
	}
	// 0 - b = -b
	if isNumeric(a.Operation) && a.Equals(0) {
		return nodeNegate(b)
	}
	if x, ok := numericInt(a); ok {
		if y, ok := numericInt(b); ok {
			return nodeNumber(fmt.Sprintf("%d", x-y))
		}
	}
	return &Node{Operation: OperationSubtract, Left: a, Right: b}
}

func nodeMul(a, b *Node) *Node {
	if a == nil || b == nil {
		return nil
	}
	if isNumeric(a.Operation) && a.Equals(0) {
		return nodeZero()
	}
	if isNumeric(b.Operation) && b.Equals(0) {
		return nodeZero()
	}
	if isNumeric(a.Operation) && a.Equals(1) {
		return b
	}
	if isNumeric(b.Operation) && b.Equals(1) {
		return a
	}
	// (-a)*b = -(a*b), a*(-b) = -(a*b)
	if a.Operation == OperationNegate {
		return nodeNegate(nodeMul(a.Left, b))
	}
	if b.Operation == OperationNegate {
		return nodeNegate(nodeMul(a, b.Left))
	}
	// (-1)*x handled above; fold small integer products for coefficients
	if isNumeric(a.Operation) && isNumeric(b.Operation) {
		x, y := big.NewInt(0), big.NewInt(0)
		if _, ok := x.SetString(a.Value, 10); ok {
			if _, ok := y.SetString(b.Value, 10); ok {
				x.Mul(x, y)
				return nodeNumber(x.String())
			}
		}
	}
	return &Node{Operation: OperationMultiply, Left: a, Right: b}
}

func nodeDiv(a, b *Node) *Node {
	if a == nil || b == nil {
		return nil
	}
	if isNumeric(b.Operation) && b.Equals(1) {
		return a
	}
	if isNumeric(b.Operation) && b.Equals(-1) {
		return nodeNegate(a)
	}
	return &Node{Operation: OperationDivide, Left: a, Right: b}
}

// collectSelf separates n into coeff*mark + rest, assuming mark appears linearly.
// Returns (nil, nil) if the dependence on mark is not linear.
func collectSelf(n, mark *Node) (coeff *Node, rest *Node) {
	if n == nil {
		return nodeZero(), nodeZero()
	}
	if n == mark {
		return nodeOne(), nodeZero()
	}
	if !containsNode(n, mark) {
		return nodeZero(), n
	}
	switch n.Operation {
	case OperationAdd:
		c1, r1 := collectSelf(n.Left, mark)
		c2, r2 := collectSelf(n.Right, mark)
		if c1 == nil || c2 == nil {
			return nil, nil
		}
		return nodeAdd(c1, c2), nodeAdd(r1, r2)
	case OperationSubtract:
		c1, r1 := collectSelf(n.Left, mark)
		c2, r2 := collectSelf(n.Right, mark)
		if c1 == nil || c2 == nil {
			return nil, nil
		}
		return nodeSub(c1, c2), nodeSub(r1, r2)
	case OperationNegate:
		c, r := collectSelf(n.Left, mark)
		if c == nil {
			return nil, nil
		}
		return nodeNegate(c), nodeNegate(r)
	case OperationMultiply:
		leftHas := containsNode(n.Left, mark)
		rightHas := containsNode(n.Right, mark)
		if leftHas && rightHas {
			return nil, nil // nonlinear
		}
		if rightHas {
			c, r := collectSelf(n.Right, mark)
			if c == nil {
				return nil, nil
			}
			return nodeMul(n.Left, c), nodeMul(n.Left, r)
		}
		if leftHas {
			c, r := collectSelf(n.Left, mark)
			if c == nil {
				return nil, nil
			}
			return nodeMul(n.Right, c), nodeMul(n.Right, r)
		}
		return nodeZero(), n
	case OperationDivide:
		if containsNode(n.Right, mark) {
			return nil, nil
		}
		if containsNode(n.Left, mark) {
			c, r := collectSelf(n.Left, mark)
			if c == nil {
				return nil, nil
			}
			return nodeDiv(c, n.Right), nodeDiv(r, n.Right)
		}
		return nodeZero(), n
	default:
		return nil, nil
	}
}

// solveSelfIntegral solves I = rest + coeff*I for I, i.e. I = rest / (1 - coeff).
func solveSelfIntegral(result, mark *Node) *Node {
	if result == nil {
		return nil
	}
	if !containsNode(result, mark) {
		return result
	}
	coeff, rest := collectSelf(result, mark)
	if coeff == nil || rest == nil {
		return nil
	}
	// denom = 1 - coeff
	denom := nodeSub(nodeOne(), coeff)
	if denom == nil {
		return nil
	}
	if isNumeric(denom.Operation) && denom.Equals(0) {
		return nil // unsolvable
	}
	return nodeDiv(rest, denom)
}

// powerOfVar returns (variable name, exponent node) if n is x or x^k, else ("", nil).
func powerOfVar(n *Node) (string, *Node) {
	if n == nil {
		return "", nil
	}
	switch n.Operation {
	case OperationVariable:
		return n.Value, nodeOne()
	case OperationExponentiation:
		if n.Left != nil && n.Left.Operation == OperationVariable && isConstantExpr(n.Right) {
			return n.Left.Value, n.Right
		}
	case OperationSquareRoot:
		// sqrt(x) = x^(1/2)
		if n.Left != nil && n.Left.Operation == OperationVariable {
			return n.Left.Value, &Node{
				Operation: OperationDivide,
				Left:      nodeOne(),
				Right:     nodeNumber("2"),
			}
		}
	}
	return "", nil
}

// linearForm reports whether n is of the form a*x + b (a, b constant; a may be 1).
// Returns (a, b, true) on success. b may be nil meaning 0.
func linearForm(n *Node) (a, b *Node, ok bool) {
	if n == nil {
		return nil, nil, false
	}
	switch n.Operation {
	case OperationVariable:
		return nodeOne(), nodeZero(), true
	case OperationNegate:
		if n.Left != nil && n.Left.Operation == OperationVariable {
			return nodeNegate(nodeOne()), nodeZero(), true
		}
		if la, lb, lok := linearForm(n.Left); lok {
			return nodeNegate(la), nodeNegate(lb), true
		}
	case OperationMultiply:
		if isConstantExpr(n.Left) && n.Right != nil && n.Right.Operation == OperationVariable {
			return n.Left, nodeZero(), true
		}
		if isConstantExpr(n.Right) && n.Left != nil && n.Left.Operation == OperationVariable {
			return n.Right, nodeZero(), true
		}
	case OperationAdd:
		// a*x + b or b + a*x
		if isConstantExpr(n.Right) {
			if la, lb, lok := linearForm(n.Left); lok && (lb == nil || (isNumeric(lb.Operation) && lb.Equals(0))) {
				return la, n.Right, true
			}
		}
		if isConstantExpr(n.Left) {
			if la, lb, lok := linearForm(n.Right); lok && (lb == nil || (isNumeric(lb.Operation) && lb.Equals(0))) {
				return la, n.Left, true
			}
		}
	case OperationSubtract:
		// a*x - b or b - a*x
		if isConstantExpr(n.Right) {
			if la, lb, lok := linearForm(n.Left); lok && (lb == nil || (isNumeric(lb.Operation) && lb.Equals(0))) {
				return la, nodeNegate(n.Right), true
			}
		}
		if isConstantExpr(n.Left) {
			if la, lb, lok := linearForm(n.Right); lok && (lb == nil || (isNumeric(lb.Operation) && lb.Equals(0))) {
				return nodeNegate(la), n.Left, true
			}
		}
	}
	return nil, nil, false
}

// chainScale divides an antiderivative F(ax+b) by the linear coefficient a.
// Supports real and complex constant a (e.g. ∫ e^(i x) dx = e^(i x)/i).
func chainScale(F, a *Node) *Node {
	if a == nil {
		return F
	}
	if isNumeric(a.Operation) && a.Equals(1) {
		return F
	}
	if isNumeric(a.Operation) && a.Equals(-1) {
		return nodeNegate(F)
	}
	// 1/i = -i, so F/i = -i·F
	if isImaginaryUnit(a) {
		negI := &Node{Operation: OperationNegate, Left: &Node{Operation: OperationImaginary, Value: "1"}}
		return nodeMul(negI, F)
	}
	// F / (-i) = i·F
	if a.Operation == OperationNegate && isImaginaryUnit(a.Left) {
		return nodeMul(&Node{Operation: OperationImaginary, Value: "1"}, F)
	}
	return nodeDiv(F, a)
}

// matchPowLog reports whether n is x^p or log(x)^k (k defaults to 1 for bare log).
func matchPowLog(n *Node) (isPow, isLog bool, exp int64, ok bool) {
	if n == nil {
		return false, false, 0, false
	}
	if n.Operation == OperationVariable {
		return true, false, 1, true
	}
	if n.Operation == OperationExponentiation && n.Left != nil &&
		n.Left.Operation == OperationVariable && isConstantExpr(n.Right) {
		if k, okk := numericInt(n.Right); okk {
			return true, false, k, true
		}
	}
	if n.Operation == OperationNaturalLogarithm && n.Left != nil &&
		n.Left.Operation == OperationVariable {
		return false, true, 1, true
	}
	if n.Operation == OperationExponentiation && n.Left != nil &&
		n.Left.Operation == OperationNaturalLogarithm &&
		n.Left.Left != nil && n.Left.Left.Operation == OperationVariable &&
		isConstantExpr(n.Right) {
		if k, okk := numericInt(n.Right); okk && k >= 1 {
			return false, true, k, true
		}
	}
	return false, false, 0, false
}

// integratePowLog computes ∫ x^n log(x)^k dx for integer n≠-1, k≥0.
// Formula: x^(n+1)/(n+1) · log^k - k/(n+1) · ∫ x^n log^(k-1) dx
func integratePowLog(a, b *Node, process func(*Node) *Node) *Node {
	var nExp, kLog int64
	var x *Node
	// a = x^n, b = log^k or log
	if isP, _, e, ok := matchPowLog(a); ok && isP {
		if _, isL, k, ok2 := matchPowLog(b); ok2 && isL {
			nExp, kLog = e, k
			if a.Operation == OperationVariable {
				x = a
			} else {
				x = a.Left
			}
		}
	}
	// a = log^k, b = x^n
	if x == nil {
		if _, isL, k, ok := matchPowLog(a); ok && isL {
			if isP, _, e, ok2 := matchPowLog(b); ok2 && isP {
				nExp, kLog = e, k
				if b.Operation == OperationVariable {
					x = b
				} else {
					x = b.Left
				}
			}
		}
	}
	if x == nil || kLog < 0 || kLog > 6 {
		return nil
	}
	if nExp == -1 {
		// ∫ x^(-1) log^k = log^(k+1) / (k+1)
		logNode := &Node{Operation: OperationNaturalLogarithm, Left: x}
		if kLog == 0 {
			return logNode
		}
		next := &Node{
			Operation: OperationExponentiation,
			Left:      logNode,
			Right:     intNode(kLog + 1),
		}
		return nodeDiv(next, intNode(kLog+1))
	}

	// Recursive reduction on k
	var reduce func(k int64) *Node
	reduce = func(k int64) *Node {
		// ∫ x^n = x^(n+1)/(n+1)
		xp1 := &Node{
			Operation: OperationExponentiation,
			Left:      x,
			Right:     intNode(nExp + 1),
		}
		base := nodeDiv(xp1, intNode(nExp+1))
		if k == 0 {
			return base
		}
		logNode := &Node{Operation: OperationNaturalLogarithm, Left: x}
		var logPow *Node
		if k == 1 {
			logPow = logNode
		} else {
			logPow = &Node{
				Operation: OperationExponentiation,
				Left:      logNode,
				Right:     intNode(k),
			}
		}
		// x^(n+1)/(n+1) * log^k - k/(n+1) * ∫ x^n log^(k-1)
		term := nodeMul(base, logPow)
		rest := reduce(k - 1)
		coeff := nodeDiv(intNode(k), intNode(nExp+1))
		return nodeSub(term, nodeMul(coeff, rest))
	}
	_ = process
	return reduce(kLog)
}

// foldConstant returns a canonical integer node if n evaluates to an integer constant.
func foldConstant(n *Node) *Node {
	if n == nil || !isConstantExpr(n) {
		return n
	}
	v, ok := numericInt(n)
	if !ok {
		return n
	}
	already := isNumeric(n.Operation) ||
		(n.Operation == OperationNegate && n.Left != nil && isNumeric(n.Left.Operation))
	if already {
		return n
	}
	if v < 0 {
		return nodeNegate(nodeNumber(fmt.Sprintf("%d", -v)))
	}
	return nodeNumber(fmt.Sprintf("%d", v))
}

// rewriteIntegrand applies algebraic rewrites that help integration (bottom-up):
// constant folding, a*(b/c) → (a*b)/c, x^m*x^n → x^(m+n), signs, etc.
func rewriteIntegrand(n *Node) *Node {
	if n == nil {
		return nil
	}

	// Bottom-up: rewrite children first.
	switch n.Operation {
	case OperationNegate:
		left := rewriteIntegrand(n.Left)
		if left != n.Left {
			n = &Node{Operation: OperationNegate, Left: left}
		}
		if n.Left != nil && n.Left.Operation == OperationNegate {
			return rewriteIntegrand(n.Left.Left)
		}
		return foldConstant(n)
	case OperationAdd, OperationSubtract, OperationMultiply, OperationDivide:
		left := rewriteIntegrand(n.Left)
		right := rewriteIntegrand(n.Right)
		if left != n.Left || right != n.Right {
			n = &Node{Operation: n.Operation, Left: left, Right: right}
		}
		if folded := foldConstant(n); folded != n {
			return folded
		}
	case OperationExponentiation:
		left := rewriteIntegrand(n.Left)
		right := rewriteIntegrand(n.Right)
		if left != n.Left || right != n.Right {
			n = &Node{Operation: OperationExponentiation, Left: left, Right: right}
		}
		// x^0 → 1, x^1 → x; fold arithmetic exponents to bare numbers
		if n.Left != nil && n.Left.Operation == OperationVariable && n.Right != nil {
			if k, ok := numericInt(n.Right); ok {
				if k == 0 {
					return nodeOne()
				}
				if k == 1 {
					return n.Left
				}
				// Fold nested arithmetic exponents to a bare integer (not ±1,0).
				// Do not rewrite x^(-n) ↔ 1/x^n here — that loops with divide rules.
				canonical := isNumeric(n.Right.Operation) ||
					(n.Right.Operation == OperationNegate && n.Right.Left != nil &&
						isNumeric(n.Right.Left.Operation))
				if !canonical {
					return &Node{
						Operation: OperationExponentiation,
						Left:      n.Left,
						Right:     intNode(k),
					}
				}
			}
		}
		return n
	default:
		return foldConstant(n)
	}

	// Multiply / divide structure rules (children already rewritten).
	switch n.Operation {
	case OperationMultiply:
		l, r := n.Left, n.Right
		if l == nil || r == nil {
			return n
		}
		// 1*f → f, 0*f → 0
		if isNumeric(l.Operation) && l.Equals(1) {
			return r
		}
		if isNumeric(r.Operation) && r.Equals(1) {
			return l
		}
		if isNumeric(l.Operation) && l.Equals(0) {
			return nodeZero()
		}
		if isNumeric(r.Operation) && r.Equals(0) {
			return nodeZero()
		}
		// (-a)*b → -(a*b), a*(-b) → -(a*b)
		if l.Operation == OperationNegate {
			return rewriteIntegrand(nodeNegate(nodeMul(l.Left, r)))
		}
		if r.Operation == OperationNegate {
			return rewriteIntegrand(nodeNegate(nodeMul(l, r.Left)))
		}
		// (c * f) * g → c * (f * g) when c is constant
		if l.Operation == OperationMultiply && isConstantExpr(l.Left) {
			return rewriteIntegrand(nodeMul(l.Left, nodeMul(l.Right, r)))
		}
		if l.Operation == OperationMultiply && isConstantExpr(l.Right) {
			return rewriteIntegrand(nodeMul(l.Right, nodeMul(l.Left, r)))
		}
		// Pull constant from right only when left is not already a bare constant
		// (avoids c*(d*f) ↔ d*(c*f) rewrite loops).
		if !isConstantExpr(l) {
			if r.Operation == OperationMultiply && isConstantExpr(r.Left) {
				return rewriteIntegrand(nodeMul(r.Left, nodeMul(l, r.Right)))
			}
			if r.Operation == OperationMultiply && isConstantExpr(r.Right) {
				return rewriteIntegrand(nodeMul(r.Right, nodeMul(l, r.Left)))
			}
		}
		// a * (b/c) → (a*b)/c
		if r.Operation == OperationDivide {
			return rewriteIntegrand(&Node{
				Operation: OperationDivide,
				Left:      nodeMul(l, r.Left),
				Right:     r.Right,
			})
		}
		// (a/b) * c → (a*c)/b
		if l.Operation == OperationDivide {
			return rewriteIntegrand(&Node{
				Operation: OperationDivide,
				Left:      nodeMul(l.Left, r),
				Right:     l.Right,
			})
		}
		// x^m * x^n → x^(m+n)
		v1, e1 := powerOfVar(l)
		v2, e2 := powerOfVar(r)
		if v1 != "" && v1 == v2 {
			return rewriteIntegrand(&Node{
				Operation: OperationExponentiation,
				Left:      &Node{Operation: OperationVariable, Value: v1},
				Right:     nodeAdd(e1, e2),
			})
		}
	case OperationDivide:
		// x^m / x^n → x^(m-n)
		if v1, e1 := powerOfVar(n.Left); v1 != "" {
			if v2, e2 := powerOfVar(n.Right); v2 == v1 {
				return rewriteIntegrand(&Node{
					Operation: OperationExponentiation,
					Left:      &Node{Operation: OperationVariable, Value: v1},
					Right:     nodeSub(e1, e2),
				})
			}
		}
		// 1/x^n → x^(-n) for n≠1 (1/x stays as division for ∫ c/x = c log x).
		// 1/sqrt(x) → x^(-1/2)
		if isConstantExpr(n.Left) {
			if v, e := powerOfVar(n.Right); v != "" {
				// Skip pure 1/x (exponent 1) to avoid fighting divide integration rules.
				if k, ok := numericInt(e); ok && k == 1 {
					// leave as c/x
				} else {
					negExp := nodeNegate(e)
					if kk, okk := numericInt(e); okk {
						negExp = intNode(-kk)
					}
					pow := &Node{
						Operation: OperationExponentiation,
						Left:      &Node{Operation: OperationVariable, Value: v},
						Right:     negExp,
					}
					if isNumeric(n.Left.Operation) && n.Left.Equals(1) {
						return rewriteIntegrand(pow)
					}
					return rewriteIntegrand(nodeMul(n.Left, pow))
				}
			}
			// 1/(c*x) → (1/c)*(1/x)
			if n.Right != nil && n.Right.Operation == OperationMultiply {
				if isConstantExpr(n.Right.Left) && n.Right.Right != nil &&
					n.Right.Right.Operation == OperationVariable {
					return rewriteIntegrand(nodeMul(
						nodeDiv(n.Left, n.Right.Left),
						&Node{
							Operation: OperationDivide,
							Left:      nodeOne(),
							Right:     n.Right.Right,
						},
					))
				}
				if isConstantExpr(n.Right.Right) && n.Right.Left != nil &&
					n.Right.Left.Operation == OperationVariable {
					return rewriteIntegrand(nodeMul(
						nodeDiv(n.Left, n.Right.Right),
						&Node{
							Operation: OperationDivide,
							Left:      nodeOne(),
							Right:     n.Right.Left,
						},
					))
				}
			}
		}
		if n.Right != nil && n.Right.Operation == OperationVariable {
			name := n.Right.Value
			left := n.Left
			if left != nil && left.Operation == OperationVariable && left.Value == name {
				return nodeOne()
			}
			if v, e := powerOfVar(left); v == name {
				return rewriteIntegrand(&Node{
					Operation: OperationExponentiation,
					Left:      &Node{Operation: OperationVariable, Value: name},
					Right:     nodeSub(e, nodeOne()),
				})
			}
			if left != nil && left.Operation == OperationMultiply {
				if isConstantExpr(left.Left) {
					if v, e := powerOfVar(left.Right); v == name {
						return rewriteIntegrand(nodeMul(left.Left, &Node{
							Operation: OperationExponentiation,
							Left:      &Node{Operation: OperationVariable, Value: name},
							Right:     nodeSub(e, nodeOne()),
						}))
					}
				}
				if isConstantExpr(left.Right) {
					if v, e := powerOfVar(left.Left); v == name {
						return rewriteIntegrand(nodeMul(left.Right, &Node{
							Operation: OperationExponentiation,
							Left:      &Node{Operation: OperationVariable, Value: name},
							Right:     nodeSub(e, nodeOne()),
						}))
					}
				}
			}
			if left != nil && left.Operation == OperationDivide && isConstantExpr(left.Right) {
				if v, e := powerOfVar(left.Left); v == name {
					return rewriteIntegrand(nodeDiv(&Node{
						Operation: OperationExponentiation,
						Left:      &Node{Operation: OperationVariable, Value: name},
						Right:     nodeSub(e, nodeOne()),
					}, left.Right))
				}
			}
		}
	}
	return n
}

// Integrate integrates the expression symbolically (indefinite integral).
// The constant of integration is omitted. Assumes a single free variable.
// Products of non-constants are handled with integration by parts (LIATE),
// including cyclic cases such as ∫ e^x sin(x) dx.
func (n *Node) Integrate() *Node {
	variable := findVariable(n)
	varNode := func() *Node {
		return &Node{
			Operation: OperationVariable,
			Value:     variable,
		}
	}
	number := func(v string) *Node {
		return nodeNumber(v)
	}

	// Per-integrand markers for cyclic integration by parts.
	// When ∫f reappears while computing ∫f, we return marks[key] and later
	// solve I = rest + coeff*I for that specific integrand only.
	marks := map[string]*Node{}
	const maxPartsDepth = 8
	partsDepth := 0
	integrating := make([]string, 0, maxPartsDepth)

	markFor := func(key string) *Node {
		if m, ok := marks[key]; ok {
			return m
		}
		m := &Node{Operation: OperationNoop, Value: "__ibp__"}
		marks[key] = m
		return m
	}

	isAnyMark := func(n *Node) bool {
		if n == nil {
			return false
		}
		for _, m := range marks {
			if n == m {
				return true
			}
		}
		return false
	}

	var process func(n *Node) *Node
	process = func(n *Node) *Node {
		if n == nil {
			return nil
		}

		// Algebraic rewrites (e.g. after IBP produces (x^2/2)*(1/x)).
		if rewritten := rewriteIntegrand(n); rewritten != n {
			return process(rewritten)
		}

		key := productKey(n)
		for _, s := range integrating {
			if s == key {
				return markFor(s)
			}
		}

		switch n.Operation {
		case OperationNoop:
			return n
		case OperationAdd:
			left, right := process(n.Left), process(n.Right)
			if left == nil || right == nil {
				return nil
			}
			return &Node{
				Operation: OperationAdd,
				Left:      left,
				Right:     right,
			}
		case OperationSubtract:
			left, right := process(n.Left), process(n.Right)
			if left == nil || right == nil {
				return nil
			}
			return &Node{
				Operation: OperationSubtract,
				Left:      left,
				Right:     right,
			}
		case OperationMultiply:
			// Constant multiple rule: ∫ c·f = c·∫f
			if isConstantExpr(n.Left) {
				right := process(n.Right)
				if right == nil {
					return nil
				}
				// Preserve cyclic marks: c * mark
				return &Node{
					Operation: OperationMultiply,
					Left:      n.Left,
					Right:     right,
				}
			}
			if isConstantExpr(n.Right) {
				left := process(n.Left)
				if left == nil {
					return nil
				}
				return &Node{
					Operation: OperationMultiply,
					Left:      n.Right,
					Right:     left,
				}
			}
			// Expand (a±b)*c for integration when one factor is a sum
			if n.Left.Operation == OperationAdd || n.Left.Operation == OperationSubtract {
				return process(&Node{
					Operation: n.Left.Operation,
					Left:      &Node{Operation: OperationMultiply, Left: n.Left.Left, Right: n.Right},
					Right:     &Node{Operation: OperationMultiply, Left: n.Left.Right, Right: n.Right},
				})
			}
			if n.Right.Operation == OperationAdd || n.Right.Operation == OperationSubtract {
				return process(&Node{
					Operation: n.Right.Operation,
					Left:      &Node{Operation: OperationMultiply, Left: n.Left, Right: n.Right.Left},
					Right:     &Node{Operation: OperationMultiply, Left: n.Left, Right: n.Right.Right},
				})
			}

			// ∫ x^n · log(x)^k via reduction (more reliable than generic IBP).
			if r := integratePowLog(n.Left, n.Right, process); r != nil {
				return r
			}
			if r := integratePowLog(n.Right, n.Left, process); r != nil {
				return r
			}

			// Integration by parts: ∫ u dv = u v - ∫ v du
			if partsDepth >= maxPartsDepth {
				return nil
			}

			tryParts := func(u, dv *Node) *Node {
				du := u.Derivative()
				if du == nil {
					return nil
				}
				if s := du.Simplify(); s != nil {
					du = s
				}
				du = rewriteIntegrand(du)
				if du == nil {
					return nil
				}

				// ∫ dv must succeed without relying on further IBP of the same product.
				partsDepth++
				v := process(dv)
				partsDepth--
				if v == nil || isAnyMark(v) {
					return nil
				}

				uv := nodeMul(u, v)
				if isNumeric(du.Operation) && du.Equals(0) {
					return uv
				}

				vdu := rewriteIntegrand(nodeMul(v, du))
				mark := markFor(key)
				partsDepth++
				integrating = append(integrating, key)
				integralVdu := process(vdu)
				integrating = integrating[:len(integrating)-1]
				partsDepth--
				if integralVdu == nil {
					return nil
				}

				result := nodeSub(uv, integralVdu)
				// Only solve for this integrand's own mark (not an ancestor's).
				if containsNode(result, mark) {
					return solveSelfIntegral(result, mark)
				}
				return result
			}

			// LIATE: lower priority becomes u. Only try the preferred order
			// first; fall back to the swapped order once if needed.
			pl, pr := partsPriority(n.Left), partsPriority(n.Right)
			if pl <= pr {
				if result := tryParts(n.Left, n.Right); result != nil {
					return result
				}
				return tryParts(n.Right, n.Left)
			}
			if result := tryParts(n.Right, n.Left); result != nil {
				return result
			}
			return tryParts(n.Left, n.Right)
		case OperationDivide:
			// ∫ f/c = (1/c)·∫f when c is constant
			if isConstantExpr(n.Right) {
				left := process(n.Left)
				if left == nil {
					return nil
				}
				return &Node{
					Operation: OperationDivide,
					Left:      left,
					Right:     n.Right,
				}
			}
			// ∫ c/x = c·log(x)
			if isConstantExpr(n.Left) && n.Right.Operation == OperationVariable {
				return &Node{
					Operation: OperationMultiply,
					Left:      n.Left,
					Right: &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Right,
					},
				}
			}
			// ∫ c/(a*x+b) = (c/a)·log(a*x+b)
			if isConstantExpr(n.Left) {
				if a, _, ok := linearForm(n.Right); ok {
					logTerm := &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Right,
					}
					scaled := chainScale(logTerm, a)
					if isNumeric(n.Left.Operation) && n.Left.Equals(1) {
						return scaled
					}
					return nodeMul(n.Left, scaled)
				}
			}
			// ∫ log(x)/x = (log(x))^2 / 2  (substitution u=log(x))
			if n.Left != nil && n.Left.Operation == OperationNaturalLogarithm &&
				n.Left.Left != nil && n.Left.Left.Operation == OperationVariable &&
				n.Right != nil && n.Right.Operation == OperationVariable &&
				n.Left.Left.Value == n.Right.Value {
				return &Node{
					Operation: OperationDivide,
					Left: &Node{
						Operation: OperationExponentiation,
						Left:      n.Left,
						Right:     number("2"),
					},
					Right: number("2"),
				}
			}
			// ∫ cos(x)/sin(x) = log(sin(x)), ∫ sin(x)/cos(x) = -log(cos(x))
			if n.Left != nil && n.Right != nil {
				if n.Left.Operation == OperationCosine && n.Right.Operation == OperationSine &&
					nodesEqual(n.Left.Left, n.Right.Left) {
					return &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Right,
					}
				}
				if n.Left.Operation == OperationSine && n.Right.Operation == OperationCosine &&
					nodesEqual(n.Left.Left, n.Right.Left) {
					return &Node{
						Operation: OperationNegate,
						Left: &Node{
							Operation: OperationNaturalLogarithm,
							Left:      n.Right,
						},
					}
				}
			}
			// ∫ f/x when f is not constant: treat as f · x^(-1)
			if partsDepth < maxPartsDepth && n.Right.Operation == OperationVariable {
				inv := &Node{
					Operation: OperationExponentiation,
					Left:      n.Right,
					Right:     nodeNegate(nodeOne()),
				}
				product := rewriteIntegrand(&Node{
					Operation: OperationMultiply,
					Left:      n.Left,
					Right:     inv,
				})
				partsDepth++
				result := process(product)
				partsDepth--
				return result
			}
			// ∫ 1/x^n rewritten by rewriteIntegrand; fall through product form
			if partsDepth < maxPartsDepth && isConstantExpr(n.Left) {
				if _, e := powerOfVar(n.Right); e != nil {
					rewritten := rewriteIntegrand(n)
					if rewritten != n {
						return process(rewritten)
					}
				}
			}
			return nil
		case OperationModulus:
			return nil
		case OperationExponentiation:
			// ∫ x^n dx
			if n.Left != nil && n.Left.Operation == OperationVariable && isConstantExpr(n.Right) {
				// Special case: ∫ x^(-1) = log(x)
				if k, ok := numericInt(n.Right); ok && k == -1 {
					return &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Left,
					}
				}
				if (isNumeric(n.Right.Operation) && n.Right.Equals(-1)) ||
					(n.Right.Operation == OperationNegate && n.Right.Left != nil &&
						isNumeric(n.Right.Left.Operation) && n.Right.Left.Equals(1)) {
					return &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Left,
					}
				}
				// ∫ x^n = x^(n+1)/(n+1)
				exp := &Node{
					Operation: OperationAdd,
					Left:      n.Right,
					Right:     number("1"),
				}
				power := &Node{
					Operation: OperationExponentiation,
					Left:      n.Left,
					Right:     exp,
				}
				return &Node{
					Operation: OperationDivide,
					Left:      power,
					Right:     exp,
				}
			}
			// ∫ (a*x+b)^n dx = (a*x+b)^(n+1) / ((n+1)*a)  for n ≠ -1
			if isConstantExpr(n.Right) {
				if a, _, ok := linearForm(n.Left); ok {
					if k, okk := numericInt(n.Right); okk && k == -1 {
						return chainScale(&Node{
							Operation: OperationNaturalLogarithm,
							Left:      n.Left,
						}, a)
					}
					exp := nodeAdd(n.Right, nodeOne())
					power := &Node{
						Operation: OperationExponentiation,
						Left:      n.Left,
						Right:     exp,
					}
					return chainScale(nodeDiv(power, exp), a)
				}
			}
			// ∫ sin(x)^2 = x/2 - sin(x)cos(x)/2
			if n.Left != nil && n.Left.Operation == OperationSine &&
				n.Left.Left != nil && n.Left.Left.Operation == OperationVariable &&
				isNumeric(n.Right.Operation) && n.Right.Equals(2) {
				x := n.Left.Left
				halfX := nodeDiv(x, number("2"))
				sincos := nodeDiv(nodeMul(n.Left, &Node{
					Operation: OperationCosine,
					Left:      x,
				}), number("2"))
				return nodeSub(halfX, sincos)
			}
			// ∫ cos(x)^2 = x/2 + sin(x)cos(x)/2
			if n.Left != nil && n.Left.Operation == OperationCosine &&
				n.Left.Left != nil && n.Left.Left.Operation == OperationVariable &&
				isNumeric(n.Right.Operation) && n.Right.Equals(2) {
				x := n.Left.Left
				halfX := nodeDiv(x, number("2"))
				sincos := nodeDiv(nodeMul(&Node{
					Operation: OperationSine,
					Left:      x,
				}, n.Left), number("2"))
				return nodeAdd(halfX, sincos)
			}
			// ∫ tan(x)^2 = tan(x) - x  (since tan^2 = sec^2 - 1)
			if n.Left != nil && n.Left.Operation == OperationTangent &&
				n.Left.Left != nil && n.Left.Left.Operation == OperationVariable &&
				isNumeric(n.Right.Operation) && n.Right.Equals(2) {
				return nodeSub(n.Left, n.Left.Left)
			}
			// ∫ log(x)^2 dx via IBP: x·log(x)^2 - 2∫ log(x) dx
			if n.Left != nil && n.Left.Operation == OperationNaturalLogarithm &&
				n.Left.Left != nil && n.Left.Left.Operation == OperationVariable &&
				isConstantExpr(n.Right) {
				if k, ok := numericInt(n.Right); ok && k >= 1 && k <= 4 {
					// General: ∫ log^n = x log^n - n ∫ log^(n-1)
					x := n.Left.Left
					logn := n
					var ibpLog func(power int) *Node
					ibpLog = func(power int) *Node {
						if power == 0 {
							return x
						}
						if power == 1 {
							return nodeSub(nodeMul(x, n.Left), x)
						}
						logp := &Node{
							Operation: OperationExponentiation,
							Left:      n.Left,
							Right:     number(fmt.Sprintf("%d", power)),
						}
						if int64(power) == k {
							logp = logn
						}
						rest := ibpLog(power - 1)
						return nodeSub(nodeMul(x, logp), nodeMul(number(fmt.Sprintf("%d", power)), rest))
					}
					return ibpLog(int(k)) // k fits in int (bounded ≤4)
				}
			}
			// ∫ a^x dx = a^x / log(a) for constant a
			if isConstantExpr(n.Left) && n.Right != nil && n.Right.Operation == OperationVariable {
				return &Node{
					Operation: OperationDivide,
					Left:      n,
					Right: &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Left,
					},
				}
			}
			// ∫ a^(c*x) = a^(c*x) / (c·log(a))
			if isConstantExpr(n.Left) && n.Right != nil {
				if a, b, ok := linearForm(n.Right); ok && (b == nil || (isNumeric(b.Operation) && b.Equals(0))) {
					return chainScale(&Node{
						Operation: OperationDivide,
						Left:      n,
						Right: &Node{
							Operation: OperationNaturalLogarithm,
							Left:      n.Left,
						},
					}, a)
				}
			}
			// ∫ c^k for constant base and exponent → c^k · x
			if isConstantExpr(n.Left) && isConstantExpr(n.Right) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		case OperationNegate:
			left := process(n.Left)
			if left == nil {
				return nil
			}
			return &Node{
				Operation: OperationNegate,
				Left:      left,
			}
		case OperationVariable:
			// ∫ x dx = x^2 / 2
			return &Node{
				Operation: OperationDivide,
				Left: &Node{
					Operation: OperationExponentiation,
					Left:      n,
					Right:     number("2"),
				},
				Right: number("2"),
			}
		case OperationImaginary, OperationNumber, OperationNotation, OperationNatural, OperationPI:
			// ∫ c dx = c · x
			return &Node{
				Operation: OperationMultiply,
				Left:      n,
				Right:     varNode(),
			}
		case OperationNaturalExponentiation:
			// ∫ e^x dx = e^x
			if n.Left != nil && n.Left.Operation == OperationVariable {
				return n
			}
			// ∫ e^c dx = e^c · x
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			// ∫ e^(a*x+b) dx = e^(a*x+b) / a
			if a, _, ok := linearForm(n.Left); ok {
				return chainScale(n, a)
			}
			// ∫ e^(c·x) dx = e^(c·x) / c  (explicit multiply form)
			if n.Left != nil && n.Left.Operation == OperationMultiply {
				var coeff *Node
				if isConstantExpr(n.Left.Left) && n.Left.Right != nil && n.Left.Right.Operation == OperationVariable {
					coeff = n.Left.Left
				} else if isConstantExpr(n.Left.Right) && n.Left.Left != nil && n.Left.Left.Operation == OperationVariable {
					coeff = n.Left.Right
				}
				if coeff != nil {
					return &Node{
						Operation: OperationDivide,
						Left:      n,
						Right:     coeff,
					}
				}
			}
			// ∫ e^(-x) when Negate(variable)
			if n.Left != nil && n.Left.Operation == OperationNegate &&
				n.Left.Left != nil && n.Left.Left.Operation == OperationVariable {
				return nodeNegate(n)
			}
			return nil
		case OperationNaturalLogarithm:
			// ∫ log(x) dx = x·log(x) - x
			if n.Left != nil && n.Left.Operation == OperationVariable {
				return &Node{
					Operation: OperationSubtract,
					Left: &Node{
						Operation: OperationMultiply,
						Left:      n.Left,
						Right:     n,
					},
					Right: n.Left,
				}
			}
			// ∫ log(a*x+b) dx via IBP / formula:
			// ((a*x+b)·log(a*x+b) - (a*x+b)) / a
			if a, _, ok := linearForm(n.Left); ok {
				arg := n.Left
				return chainScale(nodeSub(nodeMul(arg, n), arg), a)
			}
			// ∫ log(c) dx = log(c) · x
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		case OperationSquareRoot:
			// ∫ sqrt(x) dx = (2/3) · x · sqrt(x)
			if n.Left != nil && n.Left.Operation == OperationVariable {
				twoThirds := &Node{
					Operation: OperationDivide,
					Left:      number("2"),
					Right:     number("3"),
				}
				xSqrt := &Node{
					Operation: OperationMultiply,
					Left:      n.Left,
					Right:     n,
				}
				return &Node{
					Operation: OperationMultiply,
					Left:      twoThirds,
					Right:     xSqrt,
				}
			}
			// ∫ sqrt(a*x+b) = (2/3)* (a*x+b)^(3/2) / a
			if a, _, ok := linearForm(n.Left); ok {
				// (2/3) * (ax+b) * sqrt(ax+b) / a
				body := nodeMul(n.Left, n)
				return chainScale(nodeMul(nodeDiv(number("2"), number("3")), body), a)
			}
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		case OperationCosine:
			// ∫ cos(x) dx = sin(x); ∫ cos(ax+b) = sin(ax+b)/a
			if n.Left != nil && n.Left.Operation == OperationVariable {
				return &Node{
					Operation: OperationSine,
					Left:      n.Left,
				}
			}
			if a, _, ok := linearForm(n.Left); ok {
				return chainScale(&Node{
					Operation: OperationSine,
					Left:      n.Left,
				}, a)
			}
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		case OperationSine:
			// ∫ sin(x) dx = -cos(x); ∫ sin(ax+b) = -cos(ax+b)/a
			if n.Left != nil && n.Left.Operation == OperationVariable {
				return &Node{
					Operation: OperationNegate,
					Left: &Node{
						Operation: OperationCosine,
						Left:      n.Left,
					},
				}
			}
			if a, _, ok := linearForm(n.Left); ok {
				return chainScale(nodeNegate(&Node{
					Operation: OperationCosine,
					Left:      n.Left,
				}), a)
			}
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		case OperationTangent:
			// ∫ tan(x) dx = -log(cos(x)); ∫ tan(ax+b) = -log(cos(ax+b))/a
			if n.Left != nil && n.Left.Operation == OperationVariable {
				return &Node{
					Operation: OperationNegate,
					Left: &Node{
						Operation: OperationNaturalLogarithm,
						Left: &Node{
							Operation: OperationCosine,
							Left:      n.Left,
						},
					},
				}
			}
			if a, _, ok := linearForm(n.Left); ok {
				return chainScale(nodeNegate(&Node{
					Operation: OperationNaturalLogarithm,
					Left: &Node{
						Operation: OperationCosine,
						Left:      n.Left,
					},
				}), a)
			}
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		}
		return nil
	}
	return process(n)
}

// numeric marks real (non-complex) literal kinds used for 0/1 folding.
// OperationImaginary is intentionally excluded — "1i" must not be treated as 1.
var numeric = map[Operation]bool{
	OperationNumber:   true,
	OperationNotation: true,
}

func isNumeric(operation Operation) bool {
	return numeric[operation]
}

// isImaginaryUnit reports whether n is the pure imaginary unit i (1i).
func isImaginaryUnit(n *Node) bool {
	if n == nil || n.Operation != OperationImaginary {
		return false
	}
	return n.Value == "" || n.Value == "1" || n.Value == "+1"
}

// imagCoeff returns the integer coefficient of a pure imaginary literal (k·i).
func imagCoeff(n *Node) (int64, bool) {
	if n == nil || n.Operation != OperationImaginary {
		return 0, false
	}
	if n.Value == "" || n.Value == "1" || n.Value == "+1" {
		return 1, true
	}
	if n.Value == "-1" {
		return -1, true
	}
	v := new(big.Int)
	if _, ok := v.SetString(n.Value, 10); !ok || !v.IsInt64() {
		return 0, false
	}
	return v.Int64(), true
}

func imagNode(k int64) *Node {
	if k == 0 {
		return nodeZero()
	}
	if k < 0 {
		return &Node{Operation: OperationNegate, Left: imagNode(-k)}
	}
	return &Node{Operation: OperationImaginary, Value: fmt.Sprintf("%d", k)}
}

func intNode(v int64) *Node {
	if v < 0 {
		return &Node{
			Operation: OperationNegate,
			Left:      nodeNumber(fmt.Sprintf("%d", -v)),
		}
	}
	return nodeNumber(fmt.Sprintf("%d", v))
}

// asRational returns n as a reduced integer fraction num/den when possible.
func asRational(n *Node) (num, den int64, ok bool) {
	if n == nil {
		return 0, 0, false
	}
	if v, ok := numericInt(n); ok {
		return v, 1, true
	}
	if n.Operation == OperationDivide {
		a, ok1 := numericInt(n.Left)
		b, ok2 := numericInt(n.Right)
		if ok1 && ok2 && b != 0 {
			return a, b, true
		}
	}
	if n.Operation == OperationNegate {
		if a, b, ok := asRational(n.Left); ok {
			return -a, b, true
		}
	}
	return 0, 0, false
}

func ratNode(num, den int64) *Node {
	if den == 0 {
		return &Node{Operation: OperationNumber, Value: "+Inf"}
	}
	if den < 0 {
		num, den = -num, -den
	}
	g := gcdInt(num, den)
	num /= g
	den /= g
	if den == 1 {
		return intNode(num)
	}
	return &Node{Operation: OperationDivide, Left: intNode(num), Right: intNode(den)}
}

func gcdInt(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a == 0 {
		return 1
	}
	return a
}

// nodesEqual reports structural equality of two expression trees.
func nodesEqual(a, b *Node) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Operation != b.Operation || a.Value != b.Value {
		return false
	}
	return nodesEqual(a.Left, b.Left) && nodesEqual(a.Right, b.Right)
}

// mulConst multiplies expression f by integer coefficient c, canceling
// against a division denominator when possible.
func mulConst(c int64, f *Node) *Node {
	if f == nil {
		return nil
	}
	if c == 0 {
		return nodeZero()
	}
	if c == 1 {
		return f
	}
	if c == -1 {
		return &Node{Operation: OperationNegate, Left: f}
	}
	// c * (a/b) with integer b: cancel gcd(c, b)
	if f.Operation == OperationDivide {
		if d, ok := numericInt(f.Right); ok && d != 0 {
			g := gcdInt(c, d)
			c2, d2 := c/g, d/g
			num := f.Left
			if c2 != 1 && c2 != -1 {
				num = &Node{Operation: OperationMultiply, Left: intNode(c2), Right: num}
			} else if c2 == -1 {
				num = &Node{Operation: OperationNegate, Left: num}
			}
			if d2 == 1 {
				return num
			}
			if d2 == -1 {
				return &Node{Operation: OperationNegate, Left: num}
			}
			return &Node{Operation: OperationDivide, Left: num, Right: intNode(d2)}
		}
	}
	// c * (-f) = (-c) * f
	if f.Operation == OperationNegate {
		return mulConst(-c, f.Left)
	}
	return &Node{Operation: OperationMultiply, Left: intNode(c), Right: f}
}

// extractFactor returns cofactor such that n = cofactor * factor, if possible.
func extractFactor(n, factor *Node) (*Node, bool) {
	if n == nil || factor == nil {
		return nil, false
	}
	if nodesEqual(n, factor) {
		return nodeOne(), true
	}
	if n.Operation == OperationNegate {
		if c, ok := extractFactor(n.Left, factor); ok {
			return &Node{Operation: OperationNegate, Left: c}, true
		}
		return nil, false
	}
	if n.Operation == OperationMultiply {
		if nodesEqual(n.Left, factor) {
			return n.Right, true
		}
		if nodesEqual(n.Right, factor) {
			return n.Left, true
		}
		// c * (a * factor) or c * (factor * a)
		if n.Right != nil && n.Right.Operation == OperationMultiply {
			if nodesEqual(n.Right.Right, factor) {
				return &Node{Operation: OperationMultiply, Left: n.Left, Right: n.Right.Left}, true
			}
			if nodesEqual(n.Right.Left, factor) {
				return &Node{Operation: OperationMultiply, Left: n.Left, Right: n.Right.Right}, true
			}
		}
		if n.Left != nil && n.Left.Operation == OperationMultiply {
			if nodesEqual(n.Left.Right, factor) {
				return &Node{Operation: OperationMultiply, Left: n.Left.Left, Right: n.Right}, true
			}
			if nodesEqual(n.Left.Left, factor) {
				return &Node{Operation: OperationMultiply, Left: n.Left.Right, Right: n.Right}, true
			}
		}
	}
	return nil, false
}

func factorCandidates(n *Node) []*Node {
	if n == nil {
		return nil
	}
	out := []*Node{n}
	if n.Operation == OperationMultiply {
		out = append(out, n.Left, n.Right)
		if n.Right != nil && n.Right.Operation == OperationMultiply {
			out = append(out, n.Right.Left, n.Right.Right)
		}
		if n.Left != nil && n.Left.Operation == OperationMultiply {
			out = append(out, n.Left.Left, n.Left.Right)
		}
	}
	if n.Operation == OperationNegate && n.Left != nil {
		out = append(out, factorCandidates(n.Left)...)
	}
	return out
}

// factorCommon factors a common multiplier out of a sum or difference:
// a*c ± b*c → (a±b)*c, including nested products like a*c ± k*(b*c).
func factorCommon(op Operation, left, right *Node) *Node {
	if left == nil || right == nil {
		return nil
	}
	seen := map[string]bool{}
	try := func(factor *Node) *Node {
		if factor == nil {
			return nil
		}
		// Don't factor out pure numbers (would just reshuffle constants).
		if _, ok := numericInt(factor); ok {
			return nil
		}
		key := factor.String()
		if seen[key] {
			return nil
		}
		seen[key] = true
		a, ok1 := extractFactor(left, factor)
		b, ok2 := extractFactor(right, factor)
		if !ok1 || !ok2 {
			return nil
		}
		// Avoid trivial factoring of the whole expression as 1*expr ± 0.
		if nodesEqual(a, nodeOne()) && nodesEqual(b, nodeOne()) && op == OperationSubtract {
			return nodeZero()
		}
		inner := &Node{Operation: op, Left: a, Right: b}
		return &Node{Operation: OperationMultiply, Left: inner, Right: factor}
	}

	for _, f := range factorCandidates(left) {
		if r := try(f); r != nil {
			return r
		}
	}
	for _, f := range factorCandidates(right) {
		if r := try(f); r != nil {
			return r
		}
	}
	return nil
}

// Simplify simplifies an expression. Repeated passes fold constants,
// cancel factors, remove double negations, and factor common terms.
func (n *Node) Simplify() *Node {
	if n == nil {
		return nil
	}
	cur := n
	for pass := 0; pass < 10; pass++ {
		next := cur.simplifyOnce()
		if next == nil {
			return cur
		}
		if nodesEqual(cur, next) {
			return next
		}
		cur = next
	}
	return cur
}

func (n *Node) simplifyOnce() *Node {
	var process func(n *Node) *Node
	process = func(n *Node) *Node {
		if n == nil {
			return nil
		}
		switch n.Operation {
		case OperationNoop:
			return n
		case OperationDerivative:
			if n.Left == nil {
				return n
			}
			inner := process(n.Left)
			if inner == nil {
				return nil
			}
			d := inner.Derivative()
			if d == nil {
				return nil
			}
			return process(d)
		case OperationIntegrate:
			if n.Left == nil {
				return n
			}
			inner := process(n.Left)
			if inner == nil {
				return nil
			}
			i := inner.Integrate()
			if i == nil {
				return nil
			}
			return process(i)
		case OperationSimplify:
			if n.Left == nil {
				return n
			}
			return process(n.Left)
		case OperationAdd:
			left, right := process(n.Left), process(n.Right)
			if left == nil || right == nil {
				return n
			}
			// constant folding (integers and simple rationals) — reals only
			if a, ad, ok1 := asRational(left); ok1 {
				if b, bd, ok2 := asRational(right); ok2 {
					return ratNode(a*bd+b*ad, ad*bd)
				}
			}
			// (a i) + (b i) → (a+b) i
			if li, ok1 := imagCoeff(left); ok1 {
				if ri, ok2 := imagCoeff(right); ok2 {
					return imagNode(li + ri)
				}
			}
			if isNumeric(left.Operation) && left.Equals(0) {
				return right
			}
			if isNumeric(right.Operation) && right.Equals(0) {
				return left
			}
			// a + (-b) → a - b
			if right.Operation == OperationNegate {
				return process(&Node{Operation: OperationSubtract, Left: left, Right: right.Left})
			}
			// (-a) + b → b - a
			if left.Operation == OperationNegate {
				return process(&Node{Operation: OperationSubtract, Left: right, Right: left.Left})
			}
			if factored := factorCommon(OperationAdd, left, right); factored != nil {
				return process(factored)
			}
			return &Node{Operation: OperationAdd, Left: left, Right: right}

		case OperationSubtract:
			left, right := process(n.Left), process(n.Right)
			if left == nil || right == nil {
				return n
			}
			if a, ad, ok1 := asRational(left); ok1 {
				if b, bd, ok2 := asRational(right); ok2 {
					return ratNode(a*bd-b*ad, ad*bd)
				}
			}
			// (a i) - (b i) → (a-b) i
			if li, ok1 := imagCoeff(left); ok1 {
				if ri, ok2 := imagCoeff(right); ok2 {
					return imagNode(li - ri)
				}
			}
			if isNumeric(right.Operation) && right.Equals(0) {
				return left
			}
			if isNumeric(left.Operation) && left.Equals(0) {
				return process(&Node{Operation: OperationNegate, Left: right})
			}
			// a - (-b) → a + b
			if right.Operation == OperationNegate {
				return process(&Node{Operation: OperationAdd, Left: left, Right: right.Left})
			}
			// a - a → 0
			if nodesEqual(left, right) {
				return nodeZero()
			}
			// a - (b - c) → a - b + c
			if right.Operation == OperationSubtract {
				return process(&Node{
					Operation: OperationAdd,
					Left:      &Node{Operation: OperationSubtract, Left: left, Right: right.Left},
					Right:     right.Right,
				})
			}
			// a - c*(b - d) → a - c*b + c*d  (polynomial expansion)
			if right.Operation == OperationMultiply {
				if c, ok := numericInt(right.Left); ok && right.Right != nil &&
					right.Right.Operation == OperationSubtract {
					cb := &Node{Operation: OperationMultiply, Left: intNode(c), Right: right.Right.Left}
					cd := &Node{Operation: OperationMultiply, Left: intNode(c), Right: right.Right.Right}
					return process(&Node{
						Operation: OperationAdd,
						Left:      &Node{Operation: OperationSubtract, Left: left, Right: cb},
						Right:     cd,
					})
				}
			}
			// Expand c*(a±b) - d only when d shares a factor with a term inside,
			// enabling cleanup like 2*(x*sin+cos) - x^2*cos.
			if left.Operation == OperationMultiply {
				if c, ok := numericInt(left.Left); ok && left.Right != nil {
					innerOp := left.Right.Operation
					if innerOp == OperationAdd || innerOp == OperationSubtract {
						a := left.Right.Left
						b := left.Right.Right
						ca := &Node{Operation: OperationMultiply, Left: intNode(c), Right: a}
						cb := &Node{Operation: OperationMultiply, Left: intNode(c), Right: b}
						if factorCommon(OperationSubtract, ca, right) != nil ||
							factorCommon(OperationSubtract, cb, right) != nil {
							expanded := &Node{Operation: innerOp, Left: ca, Right: cb}
							return process(&Node{Operation: OperationSubtract, Left: expanded, Right: right})
						}
					}
				}
			}
			// (a+b)-c → a+(b-c) or (a-c)+b when that factors
			if left.Operation == OperationAdd {
				if factored := factorCommon(OperationSubtract, left.Right, right); factored != nil {
					return process(&Node{Operation: OperationAdd, Left: left.Left, Right: factored})
				}
				if factored := factorCommon(OperationSubtract, left.Left, right); factored != nil {
					return process(&Node{Operation: OperationAdd, Left: factored, Right: left.Right})
				}
			}
			if factored := factorCommon(OperationSubtract, left, right); factored != nil {
				return process(factored)
			}
			return &Node{Operation: OperationSubtract, Left: left, Right: right}

		case OperationMultiply:
			left, right := process(n.Left), process(n.Right)
			if left == nil || right == nil {
				return n
			}
			if isNumeric(left.Operation) && left.Equals(0) {
				return nodeZero()
			}
			if isNumeric(right.Operation) && right.Equals(0) {
				return nodeZero()
			}
			if isNumeric(left.Operation) && left.Equals(1) {
				return right
			}
			if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			// i * i = -1; (a i)*(b i) = -a*b
			if li, ok1 := imagCoeff(left); ok1 {
				if ri, ok2 := imagCoeff(right); ok2 {
					return intNode(-(li * ri))
				}
			}
			// c * (k i) → (c*k) i
			if x, ok1 := numericInt(left); ok1 {
				if k, ok2 := imagCoeff(right); ok2 {
					return imagNode(x * k)
				}
			}
			if y, ok2 := numericInt(right); ok2 {
				if k, ok1 := imagCoeff(left); ok1 {
					return imagNode(y * k)
				}
			}
			if x, ok1 := numericInt(left); ok1 {
				if y, ok2 := numericInt(right); ok2 {
					return intNode(x * y)
				}
				// c * (-f) or c * (a/k) with integer k: cancel via mulConst
				if right.Operation == OperationNegate {
					return process(mulConst(x, right))
				}
				if right.Operation == OperationDivide {
					if _, ok := numericInt(right.Right); ok {
						return process(mulConst(x, right))
					}
				}
				if x == -1 {
					return process(&Node{Operation: OperationNegate, Left: right})
				}
				return &Node{Operation: OperationMultiply, Left: intNode(x), Right: right}
			}
			if y, ok2 := numericInt(right); ok2 {
				if left.Operation == OperationNegate {
					return process(mulConst(y, left))
				}
				if left.Operation == OperationDivide {
					if _, ok := numericInt(left.Right); ok {
						return process(mulConst(y, left))
					}
				}
				if y == -1 {
					return process(&Node{Operation: OperationNegate, Left: left})
				}
				// prefer constant on the left
				return &Node{Operation: OperationMultiply, Left: intNode(y), Right: left}
			}
			// (-a)*(-b) → a*b
			if left.Operation == OperationNegate && right.Operation == OperationNegate {
				return process(&Node{Operation: OperationMultiply, Left: left.Left, Right: right.Left})
			}
			// (-a)*b → -(a*b)
			if left.Operation == OperationNegate {
				return process(&Node{
					Operation: OperationNegate,
					Left:      &Node{Operation: OperationMultiply, Left: left.Left, Right: right},
				})
			}
			if right.Operation == OperationNegate {
				return process(&Node{
					Operation: OperationNegate,
					Left:      &Node{Operation: OperationMultiply, Left: left, Right: right.Left},
				})
			}
			// a * (b/c) → (a*b)/c
			if right.Operation == OperationDivide {
				return process(&Node{
					Operation: OperationDivide,
					Left:      &Node{Operation: OperationMultiply, Left: left, Right: right.Left},
					Right:     right.Right,
				})
			}
			// (a/b) * c → (a*c)/b
			if left.Operation == OperationDivide {
				return process(&Node{
					Operation: OperationDivide,
					Left:      &Node{Operation: OperationMultiply, Left: left.Left, Right: right},
					Right:     left.Right,
				})
			}
			// (c*a)*b → c*(a*b) when c is constant
			if left.Operation == OperationMultiply {
				if c, ok := numericInt(left.Left); ok {
					inner := process(&Node{
						Operation: OperationMultiply,
						Left:      left.Right,
						Right:     right,
					})
					if c == 1 {
						return inner
					}
					if c == -1 {
						return process(&Node{Operation: OperationNegate, Left: inner})
					}
					if inner != nil && (inner.Operation == OperationDivide || inner.Operation == OperationNegate) {
						return process(mulConst(c, inner))
					}
					return &Node{Operation: OperationMultiply, Left: intNode(c), Right: inner}
				}
			}
			return &Node{Operation: OperationMultiply, Left: left, Right: right}

		case OperationDivide:
			left, right := process(n.Left), process(n.Right)
			if left == nil || right == nil {
				return n
			}
			// integer division when exact
			if x, ok1 := numericInt(left); ok1 {
				if y, ok2 := numericInt(right); ok2 && y != 0 {
					if x%y == 0 {
						return intNode(x / y)
					}
					// reduce fraction x/y
					g := gcdInt(x, y)
					x2, y2 := x/g, y/g
					if y2 < 0 {
						x2, y2 = -x2, -y2
					}
					if y2 == 1 {
						return intNode(x2)
					}
					return &Node{Operation: OperationDivide, Left: intNode(x2), Right: intNode(y2)}
				}
			}
			if isNumeric(left.Operation) && left.Equals(0) {
				return nodeZero()
			}
			if isNumeric(right.Operation) && right.Equals(0) {
				return &Node{Operation: OperationNumber, Value: "+Inf"}
			}
			if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			if y, ok := numericInt(right); ok && y == -1 {
				return process(&Node{Operation: OperationNegate, Left: left})
			}
			// a/i = -i·a ; a/(k i) = -i·a/k
			if k, ok := imagCoeff(right); ok && k != 0 {
				// 1/(k i) = -i/k
				negI := imagNode(-1)
				if k == 1 || k == -1 {
					if k == -1 {
						negI = imagNode(1)
					}
					return process(nodeMul(negI, left))
				}
				return process(nodeDiv(nodeMul(negI, left), intNode(k)))
			}
			// (a i)/(b i) = a/b
			if li, ok1 := imagCoeff(left); ok1 {
				if ri, ok2 := imagCoeff(right); ok2 && ri != 0 {
					return process(nodeDiv(intNode(li), intNode(ri)))
				}
			}
			// a/a → 1
			if nodesEqual(left, right) {
				return nodeOne()
			}
			// (a/b)/c → a/(b*c)
			if left.Operation == OperationDivide {
				return process(&Node{
					Operation: OperationDivide,
					Left:      left.Left,
					Right: &Node{
						Operation: OperationMultiply,
						Left:      left.Right,
						Right:     right,
					},
				})
			}
			// a/(b/c) → (a*c)/b
			if right.Operation == OperationDivide {
				return process(&Node{
					Operation: OperationDivide,
					Left: &Node{
						Operation: OperationMultiply,
						Left:      left,
						Right:     right.Right,
					},
					Right: right.Left,
				})
			}
			// (c*a)/d with integer c,d → cancel when gcd(|c|,|d|) > 1
			if left.Operation == OperationMultiply {
				if c, ok := numericInt(left.Left); ok {
					if d, ok2 := numericInt(right); ok2 && d != 0 {
						if g := gcdInt(c, d); g > 1 {
							return process(mulConst(c, &Node{
								Operation: OperationDivide,
								Left:      left.Right,
								Right:     intNode(d),
							}))
						}
					}
				}
				if c, ok := numericInt(left.Right); ok {
					if d, ok2 := numericInt(right); ok2 && d != 0 {
						if g := gcdInt(c, d); g > 1 {
							return process(mulConst(c, &Node{
								Operation: OperationDivide,
								Left:      left.Left,
								Right:     intNode(d),
							}))
						}
					}
				}
			}
			// (-a)/b → -(a/b), a/(-b) → -(a/b)
			if left.Operation == OperationNegate {
				return process(&Node{
					Operation: OperationNegate,
					Left:      &Node{Operation: OperationDivide, Left: left.Left, Right: right},
				})
			}
			if right.Operation == OperationNegate {
				return process(&Node{
					Operation: OperationNegate,
					Left:      &Node{Operation: OperationDivide, Left: left, Right: right.Left},
				})
			}
			return &Node{Operation: OperationDivide, Left: left, Right: right}

		case OperationModulus:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			return &Node{Operation: OperationModulus, Left: left, Right: right}

		case OperationExponentiation:
			left, right := process(n.Left), process(n.Right)
			// fold constant exponent arithmetic already done by process on right
			if x, ok1 := numericInt(left); ok1 {
				if y, ok2 := numericInt(right); ok2 && y >= 0 && y <= 20 {
					// small non-negative integer powers of integers
					res := int64(1)
					for i := int64(0); i < y; i++ {
						res *= x
					}
					return intNode(res)
				}
			}
			if isNumeric(left.Operation) && left.Equals(0) {
				// 0^0 → 1 by convention here; 0^n → 0 for n>0 handled if right not 0
				if isNumeric(right.Operation) && right.Equals(0) {
					return nodeOne()
				}
				return nodeZero()
			}
			if isNumeric(right.Operation) && right.Equals(0) {
				return nodeOne()
			}
			if isNumeric(left.Operation) && left.Equals(1) {
				return nodeOne()
			}
			if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			// x^(-1) → 1/x, x^(-n) → 1/x^n for display (simplify only)
			if left != nil && left.Operation == OperationVariable {
				if k, ok := numericInt(right); ok && k < 0 {
					if k == -1 {
						return &Node{Operation: OperationDivide, Left: nodeOne(), Right: left}
					}
					return &Node{
						Operation: OperationDivide,
						Left:      nodeOne(),
						Right: &Node{
							Operation: OperationExponentiation,
							Left:      left,
							Right:     intNode(-k),
						},
					}
				}
			}
			return &Node{Operation: OperationExponentiation, Left: left, Right: right}

		case OperationNegate:
			left := process(n.Left)
			if left == nil {
				return n
			}
			// --x → x
			if left.Operation == OperationNegate {
				return left.Left
			}
			if isNumeric(left.Operation) && left.Equals(0) {
				return nodeZero()
			}
			// fold -(-n) already; -k for number stays as Negate for display
			// -(a - b) → b - a
			if left.Operation == OperationSubtract {
				return process(&Node{Operation: OperationSubtract, Left: left.Right, Right: left.Left})
			}
			// -(a + b) → (-a) + (-b) → -a - b, keep as Negate of sum for structure
			// -(c*f) with c int → (-c)*f
			if left.Operation == OperationMultiply {
				if c, ok := numericInt(left.Left); ok {
					return process(mulConst(-c, left.Right))
				}
			}
			// Note: do not rewrite -(a/b) → (-a)/b; the divide rule does the reverse
			// and together they loop.
			return &Node{Operation: OperationNegate, Left: left}

		case OperationVariable:
			return n
		case OperationImaginary:
			return n
		case OperationNumber:
			return n
		case OperationNotation:
			return n
		case OperationRe:
			left := process(n.Left)
			// re(real) = real; re(k i) = 0
			if isNumeric(left.Operation) || left.Operation == OperationNatural || left.Operation == OperationPI {
				return left
			}
			if _, ok := imagCoeff(left); ok {
				return nodeZero()
			}
			return &Node{Operation: OperationRe, Left: left}
		case OperationIm:
			left := process(n.Left)
			// im(real) = 0; im(k i) = k
			if isNumeric(left.Operation) || left.Operation == OperationNatural || left.Operation == OperationPI {
				return nodeZero()
			}
			if k, ok := imagCoeff(left); ok {
				return intNode(k)
			}
			return &Node{Operation: OperationIm, Left: left}
		case OperationConj:
			left := process(n.Left)
			// conj(real) = real; conj(k i) = -k i
			if isNumeric(left.Operation) || left.Operation == OperationNatural || left.Operation == OperationPI {
				return left
			}
			if k, ok := imagCoeff(left); ok {
				return imagNode(-k)
			}
			// conj(a+b) not fully expanded
			return &Node{Operation: OperationConj, Left: left}
		case OperationAbs:
			left := process(n.Left)
			// abs of non-negative real literal stays as abs(...) unless zero
			if isNumeric(left.Operation) && left.Equals(0) {
				return nodeZero()
			}
			return &Node{Operation: OperationAbs, Left: left}
		case OperationNaturalExponentiation:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				return nodeOne()
			} else if isNumeric(left.Operation) && left.Equals(1) {
				return &Node{Operation: OperationNatural}
			}
			return &Node{Operation: OperationNaturalExponentiation, Left: left}
		case OperationNatural:
			return n
		case OperationPI:
			return n
		case OperationNaturalLogarithm:
			left := process(n.Left)
			if left.Operation == OperationNatural {
				return nodeOne() // log(e) = 1
			}
			return &Node{Operation: OperationNaturalLogarithm, Left: left}
		case OperationSquareRoot:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				return nodeZero()
			} else if isNumeric(left.Operation) && left.Equals(1) {
				return nodeOne()
			}
			return &Node{Operation: OperationSquareRoot, Left: left}
		case OperationCosine:
			return &Node{Operation: OperationCosine, Left: process(n.Left)}
		case OperationSine:
			return &Node{Operation: OperationSine, Left: process(n.Left)}
		case OperationTangent:
			return &Node{Operation: OperationTangent, Left: process(n.Left)}
		}
		return nil
	}
	return process(n)
}
