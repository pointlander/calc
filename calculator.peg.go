package main

// Code generated by peg calculator.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulee
	rulee1
	rulee2
	rulee3
	rulee4
	rulevalue
	rulenumber
	rulesub
	ruleadd
	ruleminus
	rulemultiply
	ruledivide
	rulemodulus
	ruleexponentiation
	ruleopen
	ruleclose
	rulesp
	rulePegText
)

var rul3s = [...]string{
	"Unknown",
	"e",
	"e1",
	"e2",
	"e3",
	"e4",
	"value",
	"number",
	"sub",
	"add",
	"minus",
	"multiply",
	"divide",
	"modulus",
	"exponentiation",
	"open",
	"close",
	"sp",
	"PegText",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Calculator struct {
	Buffer string
	buffer []rune
	rules  [19]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Calculator) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Calculator) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Calculator
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *Calculator) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Calculator) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func Pretty(pretty bool) func(*Calculator) error {
	return func(p *Calculator) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*Calculator) error {
	return func(p *Calculator) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *Calculator) Init(options ...func(*Calculator) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 e <- <(sp e1 !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulesp]() {
					goto l0
				}
				if !_rules[rulee1]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
				add(rulee, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 e1 <- <(e2 ((add e2) / (minus e2))*)> */
		func() bool {
			position3, tokenIndex3 := position, tokenIndex
			{
				position4 := position
				if !_rules[rulee2]() {
					goto l3
				}
			l5:
				{
					position6, tokenIndex6 := position, tokenIndex
					{
						position7, tokenIndex7 := position, tokenIndex
						if !_rules[ruleadd]() {
							goto l8
						}
						if !_rules[rulee2]() {
							goto l8
						}
						goto l7
					l8:
						position, tokenIndex = position7, tokenIndex7
						if !_rules[ruleminus]() {
							goto l6
						}
						if !_rules[rulee2]() {
							goto l6
						}
					}
				l7:
					goto l5
				l6:
					position, tokenIndex = position6, tokenIndex6
				}
				add(rulee1, position4)
			}
			return true
		l3:
			position, tokenIndex = position3, tokenIndex3
			return false
		},
		/* 2 e2 <- <(e3 ((multiply e3) / (divide e3) / (modulus e3))*)> */
		func() bool {
			position9, tokenIndex9 := position, tokenIndex
			{
				position10 := position
				if !_rules[rulee3]() {
					goto l9
				}
			l11:
				{
					position12, tokenIndex12 := position, tokenIndex
					{
						position13, tokenIndex13 := position, tokenIndex
						if !_rules[rulemultiply]() {
							goto l14
						}
						if !_rules[rulee3]() {
							goto l14
						}
						goto l13
					l14:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[ruledivide]() {
							goto l15
						}
						if !_rules[rulee3]() {
							goto l15
						}
						goto l13
					l15:
						position, tokenIndex = position13, tokenIndex13
						if !_rules[rulemodulus]() {
							goto l12
						}
						if !_rules[rulee3]() {
							goto l12
						}
					}
				l13:
					goto l11
				l12:
					position, tokenIndex = position12, tokenIndex12
				}
				add(rulee2, position10)
			}
			return true
		l9:
			position, tokenIndex = position9, tokenIndex9
			return false
		},
		/* 3 e3 <- <(e4 (exponentiation e4)*)> */
		func() bool {
			position16, tokenIndex16 := position, tokenIndex
			{
				position17 := position
				if !_rules[rulee4]() {
					goto l16
				}
			l18:
				{
					position19, tokenIndex19 := position, tokenIndex
					if !_rules[ruleexponentiation]() {
						goto l19
					}
					if !_rules[rulee4]() {
						goto l19
					}
					goto l18
				l19:
					position, tokenIndex = position19, tokenIndex19
				}
				add(rulee3, position17)
			}
			return true
		l16:
			position, tokenIndex = position16, tokenIndex16
			return false
		},
		/* 4 e4 <- <((minus value) / value)> */
		func() bool {
			position20, tokenIndex20 := position, tokenIndex
			{
				position21 := position
				{
					position22, tokenIndex22 := position, tokenIndex
					if !_rules[ruleminus]() {
						goto l23
					}
					if !_rules[rulevalue]() {
						goto l23
					}
					goto l22
				l23:
					position, tokenIndex = position22, tokenIndex22
					if !_rules[rulevalue]() {
						goto l20
					}
				}
			l22:
				add(rulee4, position21)
			}
			return true
		l20:
			position, tokenIndex = position20, tokenIndex20
			return false
		},
		/* 5 value <- <(number / sub)> */
		func() bool {
			position24, tokenIndex24 := position, tokenIndex
			{
				position25 := position
				{
					position26, tokenIndex26 := position, tokenIndex
					if !_rules[rulenumber]() {
						goto l27
					}
					goto l26
				l27:
					position, tokenIndex = position26, tokenIndex26
					if !_rules[rulesub]() {
						goto l24
					}
				}
			l26:
				add(rulevalue, position25)
			}
			return true
		l24:
			position, tokenIndex = position24, tokenIndex24
			return false
		},
		/* 6 number <- <(<('-'? [0-9]+ ('.' [0-9]*)?)> sp)> */
		func() bool {
			position28, tokenIndex28 := position, tokenIndex
			{
				position29 := position
				{
					position30 := position
					{
						position31, tokenIndex31 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l31
						}
						position++
						goto l32
					l31:
						position, tokenIndex = position31, tokenIndex31
					}
				l32:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l28
					}
					position++
				l33:
					{
						position34, tokenIndex34 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l34
						}
						position++
						goto l33
					l34:
						position, tokenIndex = position34, tokenIndex34
					}
					{
						position35, tokenIndex35 := position, tokenIndex
						if buffer[position] != rune('.') {
							goto l35
						}
						position++
					l37:
						{
							position38, tokenIndex38 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l38
							}
							position++
							goto l37
						l38:
							position, tokenIndex = position38, tokenIndex38
						}
						goto l36
					l35:
						position, tokenIndex = position35, tokenIndex35
					}
				l36:
					add(rulePegText, position30)
				}
				if !_rules[rulesp]() {
					goto l28
				}
				add(rulenumber, position29)
			}
			return true
		l28:
			position, tokenIndex = position28, tokenIndex28
			return false
		},
		/* 7 sub <- <(open e1 close)> */
		func() bool {
			position39, tokenIndex39 := position, tokenIndex
			{
				position40 := position
				if !_rules[ruleopen]() {
					goto l39
				}
				if !_rules[rulee1]() {
					goto l39
				}
				if !_rules[ruleclose]() {
					goto l39
				}
				add(rulesub, position40)
			}
			return true
		l39:
			position, tokenIndex = position39, tokenIndex39
			return false
		},
		/* 8 add <- <('+' sp)> */
		func() bool {
			position41, tokenIndex41 := position, tokenIndex
			{
				position42 := position
				if buffer[position] != rune('+') {
					goto l41
				}
				position++
				if !_rules[rulesp]() {
					goto l41
				}
				add(ruleadd, position42)
			}
			return true
		l41:
			position, tokenIndex = position41, tokenIndex41
			return false
		},
		/* 9 minus <- <('-' sp)> */
		func() bool {
			position43, tokenIndex43 := position, tokenIndex
			{
				position44 := position
				if buffer[position] != rune('-') {
					goto l43
				}
				position++
				if !_rules[rulesp]() {
					goto l43
				}
				add(ruleminus, position44)
			}
			return true
		l43:
			position, tokenIndex = position43, tokenIndex43
			return false
		},
		/* 10 multiply <- <('*' sp)> */
		func() bool {
			position45, tokenIndex45 := position, tokenIndex
			{
				position46 := position
				if buffer[position] != rune('*') {
					goto l45
				}
				position++
				if !_rules[rulesp]() {
					goto l45
				}
				add(rulemultiply, position46)
			}
			return true
		l45:
			position, tokenIndex = position45, tokenIndex45
			return false
		},
		/* 11 divide <- <('/' sp)> */
		func() bool {
			position47, tokenIndex47 := position, tokenIndex
			{
				position48 := position
				if buffer[position] != rune('/') {
					goto l47
				}
				position++
				if !_rules[rulesp]() {
					goto l47
				}
				add(ruledivide, position48)
			}
			return true
		l47:
			position, tokenIndex = position47, tokenIndex47
			return false
		},
		/* 12 modulus <- <('%' sp)> */
		func() bool {
			position49, tokenIndex49 := position, tokenIndex
			{
				position50 := position
				if buffer[position] != rune('%') {
					goto l49
				}
				position++
				if !_rules[rulesp]() {
					goto l49
				}
				add(rulemodulus, position50)
			}
			return true
		l49:
			position, tokenIndex = position49, tokenIndex49
			return false
		},
		/* 13 exponentiation <- <('^' sp)> */
		func() bool {
			position51, tokenIndex51 := position, tokenIndex
			{
				position52 := position
				if buffer[position] != rune('^') {
					goto l51
				}
				position++
				if !_rules[rulesp]() {
					goto l51
				}
				add(ruleexponentiation, position52)
			}
			return true
		l51:
			position, tokenIndex = position51, tokenIndex51
			return false
		},
		/* 14 open <- <('(' sp)> */
		func() bool {
			position53, tokenIndex53 := position, tokenIndex
			{
				position54 := position
				if buffer[position] != rune('(') {
					goto l53
				}
				position++
				if !_rules[rulesp]() {
					goto l53
				}
				add(ruleopen, position54)
			}
			return true
		l53:
			position, tokenIndex = position53, tokenIndex53
			return false
		},
		/* 15 close <- <(')' sp)> */
		func() bool {
			position55, tokenIndex55 := position, tokenIndex
			{
				position56 := position
				if buffer[position] != rune(')') {
					goto l55
				}
				position++
				if !_rules[rulesp]() {
					goto l55
				}
				add(ruleclose, position56)
			}
			return true
		l55:
			position, tokenIndex = position55, tokenIndex55
			return false
		},
		/* 16 sp <- <(' ' / '\t')*> */
		func() bool {
			{
				position58 := position
			l59:
				{
					position60, tokenIndex60 := position, tokenIndex
					{
						position61, tokenIndex61 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l62
						}
						position++
						goto l61
					l62:
						position, tokenIndex = position61, tokenIndex61
						if buffer[position] != rune('\t') {
							goto l60
						}
						position++
					}
				l61:
					goto l59
				l60:
					position, tokenIndex = position60, tokenIndex60
				}
				add(rulesp, position58)
			}
			return true
		},
		nil,
	}
	p.rules = _rules
	return nil
}
