// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc

import (
	"fmt"
	"math/big"
	"math/rand"
	"sort"
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
	// OperationNaturalExponentiation raises the natural number to a power
	OperationNaturalExponentiation
	// OperationNatural is the constant e
	OperationNatural
	// OperationPI is the constant pi
	OperationPI
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

var (
	// BinaryOperations are binary operations
	BinaryOperations = []Operation{
		OperationAdd,
		OperationSubtract,
		OperationMultiply,
		OperationDivide,
		//OperationModulus,
		OperationExponentiation,
		OperationNotation,
	}
	// IsBinaryOperation is a map of binary operations
	IsBinaryOperation = make(map[Operation]bool)
	// UnaryOperations are unary operations
	UnaryOperations = []Operation{
		OperationNegate,
		OperationNaturalExponentiation,
		OperationNaturalLogarithm,
		OperationSquareRoot,
		OperationCosine,
		OperationSine,
		OperationTangent,
	}
	// IsUnaryOperation is a map of unary operations
	IsUnaryOperation = make(map[Operation]bool)
	// Numbers are numbers
	Numbers = []Operation{
		OperationImaginary,
		OperationNumber,
	}
	// IsNumber is a map of number operations
	IsNumber = make(map[Operation]bool)
	// Constants are constants
	Constants = []Operation{
		OperationNatural,
		OperationPI,
	}
	// IsConstant is a map of constant operations
	IsConstant = make(map[Operation]bool)
)

func init() {
	for _, operation := range BinaryOperations {
		IsBinaryOperation[operation] = true
	}
	for _, operation := range UnaryOperations {
		IsUnaryOperation[operation] = true
	}
	for _, operation := range Numbers {
		IsNumber[operation] = true
	}
	for _, operation := range Constants {
		IsConstant[operation] = true
	}
}

// Node is a node in an expression binary tree
type Node struct {
	Operation   Operation
	Value       string
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
			value := big.NewInt(0)
			value.SetString(n.Right.Value, 10)
			value.Sub(value, big.NewInt(1))
			constant := &Node{
				Operation: OperationNumber,
				Value:     value.String(),
			}
			exp := &Node{
				Operation: OperationExponentiation,
				Left:      n.Left,
				Right:     constant,
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
			value := &Node{
				Operation: OperationNumber,
				Value:     "0.5",
			}
			multiply := &Node{
				Operation: OperationMultiply,
				Left:      value,
				Right:     process(n.Left),
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      multiply,
				Right:     n,
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

// Integrate takes the integral of the equation
func (n *Node) Integrate() *Node {
	rnd := rand.New(rand.NewSource(1))

	var flatten func(n **Node, nodes []**Node) []**Node
	flatten = func(n **Node, nodes []**Node) []**Node {
		if n == nil || *n == nil {
			return nodes
		}
		nodes = append(nodes, n)
		nodes = flatten(&(*n).Left, nodes)
		nodes = flatten(&(*n).Right, nodes)
		return nodes
	}

	var cp func(n *Node) *Node
	cp = func(n *Node) *Node {
		if n == nil {
			return nil
		}
		new := Node{}
		new = *n
		new.Left = cp(n.Left)
		new.Right = cp(n.Right)
		return &new
	}

	mutate := func(n *Node) *Node {
		n = cp(n)
		nodes := flatten(&n, nil)
		selected := nodes[rnd.Intn(len(nodes))]
		if rnd.Intn(2) == 0 {
			switch rnd.Intn(3) {
			case 0:
				*selected = &Node{
					Operation: BinaryOperations[rnd.Intn(len(BinaryOperations))],
					Left:      *selected,
				}
				if (*selected).Operation == OperationExponentiation {
					(*selected).Right = &Node{
						Operation: OperationNumber,
						Value:     fmt.Sprintf("%d", rnd.Intn(10)),
					}
				} else {
					switch rnd.Intn(3) {
					case 0:
						if rnd.Intn(2) == 0 {
							(*selected).Right = &Node{
								Operation: OperationNumber,
								Value:     fmt.Sprintf("%d", rnd.Intn(10)),
							}
						} else {
							(*selected).Right = &Node{
								Operation: OperationImaginary,
								Value:     fmt.Sprintf("%d", rnd.Intn(10)),
							}
						}
					case 1:
						(*selected).Right = &Node{
							Operation: Constants[rnd.Intn(len(Constants))],
						}
					case 2:
						(*selected).Right = &Node{
							Operation: OperationVariable,
							Value:     "x",
						}
					}
				}
			case 1:
				*selected = &Node{
					Operation: UnaryOperations[rnd.Intn(len(UnaryOperations))],
					Left:      *selected,
				}
			case 2:
				switch rnd.Intn(3) {
				case 0:
					if rnd.Intn(2) == 0 {
						*selected = &Node{
							Operation: OperationNumber,
							Value:     fmt.Sprintf("%d", rnd.Intn(10)),
						}
					} else {
						*selected = &Node{
							Operation: OperationImaginary,
							Value:     fmt.Sprintf("%d", rnd.Intn(10)),
						}
					}
				case 1:
					*selected = &Node{
						Operation: Constants[rnd.Intn(len(Constants))],
					}
				case 2:
					*selected = &Node{
						Operation: OperationVariable,
						Value:     "x",
					}
				}
			}
		} else {
			if IsBinaryOperation[(*selected).Operation] {
				(*selected).Operation = BinaryOperations[rnd.Intn(len(BinaryOperations))]
			} else if IsUnaryOperation[(*selected).Operation] {
				(*selected).Operation = UnaryOperations[rnd.Intn(len(UnaryOperations))]
			} else if IsNumber[(*selected).Operation] {
				value := big.NewInt(0)
				value.SetString((*selected).Value, 10)
				switch rnd.Intn(3) {
				case 0:
					value.SetInt64(0)
				case 1:
					value.Add(value, big.NewInt(1))
				case 2:
					value.Sub(value, big.NewInt(1))
				}
				(*selected).Value = value.String()
			} else if IsConstant[(*selected).Operation] {
				if rnd.Intn(2) == 0 {
					(*selected).Operation = OperationNumber
					(*selected).Value = "0"
				} else {
					(*selected).Operation = Constants[rnd.Intn(len(Constants))]
				}
			} else {
				if rnd.Intn(2) == 0 {
					(*selected).Operation = OperationNumber
					(*selected).Value = "0"
				}
			}
		}
		return n
	}

	var difference func(a, b *Node, diff int) int
	difference = func(a, b *Node, diff int) int {
		if a != nil && b != nil {
			if IsNumber[a.Operation] && IsNumber[b.Operation] {
				anumber, bnumber := big.NewInt(0), big.NewInt(0)
				anumber.SetString(a.Value, 10)
				bnumber.SetString(b.Value, 10)
				anumber.Sub(anumber, bnumber)
				anumber.Abs(anumber)
				diff += int(anumber.Int64())
			} else if IsConstant[a.Operation] && IsConstant[b.Operation] &&
				a.Value != b.Value {
				diff++
			} else if a.Operation == OperationVariable && b.Operation == OperationVariable &&
				a.Value != b.Value {
				diff++
			} else if a.Operation != b.Operation {
				diff += 8
			}
			diff = difference(a.Left, b.Left, diff)
			diff = difference(a.Right, b.Right, diff)
		} else if a == nil && b != nil {
			diff++
			diff = difference(nil, b.Left, diff)
			diff = difference(nil, b.Right, diff)
		} else if a != nil && b == nil {
			diff++
			diff = difference(a.Left, nil, diff)
			diff = difference(a.Right, nil, diff)
		}
		return diff
	}

	equations := make([]*Node, 0, 8)
	for i := 0; i < 1000; i++ {
		equations = append(equations, cp(n))
	}
	for {
		length := len(equations)
		for i := 0; i < length; i++ {
			if rnd.Intn(3) == 0 {
				mutated := mutate(equations[i])
				mutations := rnd.Intn(3)
				for j := 0; j < mutations; j++ {
					mutated = mutate(mutated)
				}
				equations = append(equations, mutated)
			}
		}
		sort.Slice(equations, func(i, j int) bool {
			a := equations[i].Derivative()
			b := equations[j].Derivative()
			ae := difference(a, n, len(a.String()))
			be := difference(b, n, len(b.String()))
			return ae < be
		})
		equations = equations[:1000]
		for i := 0; i < 10; i++ {
			fmt.Println(equations[i].String(), equations[i].Derivative().String())
		}
		fmt.Println("")
		if equations[0].Derivative().Simplify().String() == n.String() {
			break
		}
	}
	return equations[0]
}
