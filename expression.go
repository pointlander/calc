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
			a.A.SetString(n.Value)

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

// isConstantExpr reports whether the expression contains no variables.
func isConstantExpr(n *Node) bool {
	if n == nil {
		return true
	}
	if n.Operation == OperationVariable {
		return false
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
		if n.Left.Operation == OperationVariable && isConstantExpr(n.Right) {
			return 3 // algebraic x^n
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
	}
	return "", nil
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
				canonical := isNumeric(n.Right.Operation) ||
					(n.Right.Operation == OperationNegate && n.Right.Left != nil &&
						isNumeric(n.Right.Left.Operation))
				if !canonical {
					var exp *Node
					if k < 0 {
						exp = nodeNegate(nodeNumber(fmt.Sprintf("%d", -k)))
					} else {
						exp = nodeNumber(fmt.Sprintf("%d", k))
					}
					return &Node{
						Operation: OperationExponentiation,
						Left:      n.Left,
						Right:     exp,
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
		if r.Operation == OperationMultiply && isConstantExpr(r.Left) {
			return rewriteIntegrand(nodeMul(r.Left, nodeMul(l, r.Right)))
		}
		if r.Operation == OperationMultiply && isConstantExpr(r.Right) {
			return rewriteIntegrand(nodeMul(r.Right, nodeMul(l, r.Left)))
		}
		if l.Operation == OperationMultiply && isConstantExpr(l.Right) {
			return rewriteIntegrand(nodeMul(l.Right, nodeMul(l.Left, r)))
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
			// ∫ x^n / x^m via rewriteIntegrand handles pure powers;
			// ∫ f/x when f is not constant: treat as f * x^(-1) with depth guard
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
			return nil
		case OperationModulus:
			return nil
		case OperationExponentiation:
			// ∫ x^n dx
			if n.Left.Operation == OperationVariable && isConstantExpr(n.Right) {
				// Special case: ∫ x^(-1) = log(x)
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
			// ∫ a^x dx = a^x / log(a) for constant a
			if isConstantExpr(n.Left) && n.Right.Operation == OperationVariable {
				return &Node{
					Operation: OperationDivide,
					Left:      n,
					Right: &Node{
						Operation: OperationNaturalLogarithm,
						Left:      n.Left,
					},
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
			if n.Left.Operation == OperationVariable {
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
			// ∫ e^(c·x) dx = e^(c·x) / c
			if n.Left.Operation == OperationMultiply {
				var coeff *Node
				if isConstantExpr(n.Left.Left) && n.Left.Right.Operation == OperationVariable {
					coeff = n.Left.Left
				} else if isConstantExpr(n.Left.Right) && n.Left.Left.Operation == OperationVariable {
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
			return nil
		case OperationNaturalLogarithm:
			// ∫ log(x) dx = x·log(x) - x
			if n.Left.Operation == OperationVariable {
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
			// ∫ sqrt(x) dx = (2/3) · x^(3/2) = (2/3) · x · sqrt(x)
			if n.Left.Operation == OperationVariable {
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
			if isConstantExpr(n.Left) {
				return &Node{
					Operation: OperationMultiply,
					Left:      n,
					Right:     varNode(),
				}
			}
			return nil
		case OperationCosine:
			// ∫ cos(x) dx = sin(x)
			if n.Left.Operation == OperationVariable {
				return &Node{
					Operation: OperationSine,
					Left:      n.Left,
				}
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
			// ∫ sin(x) dx = -cos(x)
			if n.Left.Operation == OperationVariable {
				return &Node{
					Operation: OperationNegate,
					Left: &Node{
						Operation: OperationCosine,
						Left:      n.Left,
					},
				}
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
			// ∫ tan(x) dx = -log(cos(x))
			if n.Left.Operation == OperationVariable {
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
