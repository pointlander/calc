// Copyright 2020 The Calc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc

import (
	"math/big"

	"github.com/ALTree/bigfloat"
)

// Matrix is a matrix
type Matrix struct {
	Prec   uint
	Values [][]Rational
}

// NewMatrix make a new matrix
func NewMatrix(prec uint) Matrix {
	return Matrix{
		Prec: prec,
	}
}

// Add adds two matricies
func (m *Matrix) Add(a, b *Matrix) *Matrix {
	asingular := len(a.Values) == 1 && len(a.Values[0]) == 1
	bsingular := len(b.Values) == 1 && len(b.Values[0]) == 1
	if asingular || bsingular {
		if bsingular {
			a, b = b, a
		}
		value, values := a.Values[0][0], [][]Rational{}
		for _, b := range b.Values {
			var row []Rational
			for _, bb := range b {
				ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
				row = append(row, *(ab.Add(&bb, &value)))
			}
			values = append(values, row)
		}
		m.Values = values
		return m
	}

	values := [][]Rational{}
	for i, b := range b.Values {
		var row []Rational
		for j, bb := range b {
			ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
			row = append(row, *(ab.Add(&bb, &a.Values[i][j])))
		}
		values = append(values, row)
	}
	m.Values = values
	return m
}

// Sub subtracts two matricies
func (m *Matrix) Sub(a, b *Matrix) *Matrix {
	if len(a.Values) == 1 && len(a.Values[0]) == 1 &&
		len(b.Values) == 1 && len(b.Values[0]) == 1 {
		aa, bb, values := a.Values[0][0], b.Values[0][0], [][]Rational{}
		var row []Rational
		ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
		row = append(row, *(ab.Sub(&aa, &bb)))
		values = append(values, row)
		m.Values = values
		return m
	} else if len(a.Values) == 1 && len(a.Values[0]) == 1 {
		value, values := a.Values[0][0], [][]Rational{}
		for _, b := range b.Values {
			var row []Rational
			for _, bb := range b {
				ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
				row = append(row, *(ab.Sub(&value, &bb)))
			}
			values = append(values, row)
		}
		m.Values = values
		return m
	} else if len(b.Values) == 1 && len(b.Values[0]) == 1 {
		value, values := b.Values[0][0], [][]Rational{}
		for _, a := range a.Values {
			var row []Rational
			for _, aa := range a {
				ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
				row = append(row, *(ab.Sub(&aa, &value)))
			}
			values = append(values, row)
		}
		m.Values = values
		return m
	}

	values := [][]Rational{}
	for i, b := range b.Values {
		var row []Rational
		for j, bb := range b {
			ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
			row = append(row, *(ab.Sub(&a.Values[i][j], &bb)))
		}
		values = append(values, row)
	}
	m.Values = values
	return m
}

// Mul multiplies two matricies
func (m *Matrix) Mul(a, b *Matrix) *Matrix {
	asingular := len(a.Values) == 1 && len(a.Values[0]) == 1
	bsingular := len(b.Values) == 1 && len(b.Values[0]) == 1
	if asingular || bsingular {
		if bsingular {
			a, b = b, a
		}
		value, values := a.Values[0][0], [][]Rational{}
		for _, b := range b.Values {
			var row []Rational
			for _, bb := range b {
				ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
				row = append(row, *(ab.Mul(&bb, &value)))
			}
			values = append(values, row)
		}
		m.Values = values
		return m
	}

	values := [][]Rational{}
	for x := 0; x < len(a.Values); x++ {
		var row []Rational
		for y := 0; y < len(b.Values[0]); y++ {
			sum := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
			for z := 0; z < len(b.Values); z++ {
				ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
				ab.Mul(&a.Values[x][z], &b.Values[z][y])
				sum.Add(sum, ab)
			}
			row = append(row, *sum)
		}
		values = append(values, row)
	}
	m.Values = values
	return m
}

// Div divides two matricies
func (m *Matrix) Div(a, b *Matrix) *Matrix {
	if len(a.Values) == 1 && len(a.Values[0]) == 1 &&
		len(b.Values) == 1 && len(b.Values[0]) == 1 {
		aa, bb, values := a.Values[0][0], b.Values[0][0], [][]Rational{}
		var row []Rational
		ab := NewRational(big.NewRat(0, 1), big.NewRat(0, 1))
		row = append(row, *(ab.Div(&aa, &bb)))
		values = append(values, row)
		m.Values = values
		return m
	}

	panic("can't divide non 1x1 matrices")
}

func (m *Matrix) apply(a *Matrix, function func(a *Rational) *Rational) *Matrix {
	values := [][]Rational{}
	for _, a := range a.Values {
		var row []Rational
		for _, aa := range a {
			row = append(row, *(function(&aa)))
		}
		values = append(values, row)
	}
	m.Values = values
	return m
}

func (m *Matrix) apply2(a *Matrix, y *Rational, function func(a *Rational, b *Rational) *Rational) *Matrix {
	values := [][]Rational{}
	for _, a := range a.Values {
		var row []Rational
		for _, aa := range a {
			row = append(row, *(function(&aa, y)))
		}
		values = append(values, row)
	}
	m.Values = values
	return m
}

// Abs computes the absolute value of the entries of the matrix
func (m *Matrix) Abs(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Abs(x).Rat(a)
		return a
	})
}

// Conj computes the complex conjugate of a
func (m *Matrix) Conj(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Conj(x).Rat(a)
		return a
	})
}

// Re extracts the real part of each entry
func (m *Matrix) Re(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		a.B.SetInt64(0)
		return a
	})
}

// Im extracts the imaginary part of each entry (as a real number)
func (m *Matrix) Im(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		a.A.Set(a.B)
		a.B.SetInt64(0)
		return a
	})
}

// Sqrt computes the square root of the matrix
func (m *Matrix) Sqrt(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Sqrt(x).Rat(a)
		return a
	})
}

// Atan2 computes atan2 of x
// https://en.wikipedia.org/wiki/Atan2
func (m *Matrix) Atan2(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Atan2(x).Rat(a)
		return a
	})
}

// Arg computes arg(x + yi) = tan-1(y/x)
// https://mathworld.wolfram.com/ComplexArgument.html
func (m *Matrix) Arg(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Arg(x).Rat(a)
		return a
	})
}

// Exp computes e^x for a complex number
// https://www.wolframalpha.com/input/?i=e%5E%28x+%2B+yi%29
func (m *Matrix) Exp(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Exp(x).Rat(a)
		return a
	})
}

// Cos computes cosine of a number
// https://www.wolframalpha.com/input/?i=cos%28x+%2B+yi%29
func (m *Matrix) Cos(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Cos(x).Rat(a)
		return a
	})
}

// Sin computes sine of a number
// https://www.wolframalpha.com/input/?i=sin%28x+%2B+yi%29
func (m *Matrix) Sin(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Sin(x).Rat(a)
		return a
	})
}

// Tan computes tangent of a number
// https://en.wikipedia.org/wiki/Trigonometric_functions
func (m *Matrix) Tan(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Tan(x).Rat(a)
		return a
	})
}

// Log computes the natural log of x
// https://en.wikipedia.org/wiki/Complex_logarithm
func (m *Matrix) Log(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		x.Log(x).Rat(a)
		return a
	})
}

// Pow computes x**y
// https://mathworld.wolfram.com/ComplexExponentiation.html
func (m *Matrix) Pow(x *Matrix, y *Rational) *Matrix {
	return m.apply2(x, y, func(a, b *Rational) *Rational {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(a)
		y := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		y.SetRat(b)
		x.Pow(x, y).Rat(a)
		return a
	})
}

// Neg negates the rational
func (m *Matrix) Neg(a *Matrix) *Matrix {
	return m.apply(a, func(a *Rational) *Rational {
		a.Neg(a)
		return a
	})
}

func (m *Matrix) String() string {
	if len(m.Values) == 1 && len(m.Values[0]) == 1 {
		x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
		x.SetRat(&m.Values[0][0])
		return x.String()
	}

	s, last := "[", len(m.Values)-1
	for i, row := range m.Values {
		lastColumn := len(row) - 1
		for j, value := range row {
			x := NewFloat(big.NewFloat(0).SetPrec(m.Prec), big.NewFloat(0).SetPrec(m.Prec))
			x.SetRat(&value)
			s += x.String()
			if j < lastColumn {
				s += " "
			}
		}
		if i < last {
			s += ";"
		}
	}
	return s + "]"
}

// Float is an imaginary number
type Float struct {
	A *big.Float // Real part
	B *big.Float // Imaginary part
}

// NewFloat creates a new imaginary number
func NewFloat(a, b *big.Float) *Float {
	return &Float{
		A: a,
		B: b,
	}
}

// Abs computes the absolute value of a
func (f *Float) Abs(a *Float) *Float {
	f.A.Mul(a.A, a.A)
	f.B.Mul(a.B, a.B)
	f.A = bigfloat.Sqrt(f.A.Add(f.A, f.B))
	f.B = big.NewFloat(0)
	return f
}

// Add add two imaginary numbers
func (f *Float) Add(a, b *Float) *Float {
	f.A.Add(a.A, b.A)
	f.B.Add(a.B, b.B)
	return f
}

// Sub subtracts two imaginary numbers
func (f *Float) Sub(a, b *Float) *Float {
	f.A.Sub(a.A, b.A)
	f.B.Sub(a.B, b.B)
	return f
}

// Mul multiples two imaginary numbers
func (f *Float) Mul(a, b *Float) *Float {
	x1, x2, x3, x4 :=
		big.NewFloat(0).SetPrec(f.A.Prec()), big.NewFloat(0).SetPrec(f.B.Prec()),
		big.NewFloat(0).SetPrec(f.B.Prec()), big.NewFloat(0).SetPrec(f.A.Prec())
	x1.Mul(a.A, b.A) // a*a
	x2.Mul(a.A, b.B) // a*ib
	x3.Mul(a.B, b.A) // ib*a
	x4.Mul(a.B, b.B) // i^2 * b = -b
	f.A.Add(x1, x4.Neg(x4))
	f.B.Add(x2, x3)
	return f
}

// Conj computes the complex conjugate of a
func (f *Float) Conj(a *Float) *Float {
	f.A.Set(a.A)
	f.B.Neg(a.B)
	return f
}

// Div divides two imaginary numbers: (a+bi)/(c+di) = ((ac+bd)+(bc-ad)i)/(c²+d²).
func (f *Float) Div(a, b *Float) *Float {
	prec := f.A.Prec()
	if f.B.Prec() > prec {
		prec = f.B.Prec()
	}
	aa := new(big.Float).Copy(a.A)
	bb := new(big.Float).Copy(a.B)
	cc := new(big.Float).Copy(b.A)
	dd := new(big.Float).Copy(b.B)
	// |b|² = c² + d²
	denom := new(big.Float).SetPrec(prec)
	denom.Add(new(big.Float).Mul(cc, cc), new(big.Float).Mul(dd, dd))
	// real = (ac + bd) / |b|²
	// imag = (bc - ad) / |b|²
	ac := new(big.Float).Mul(aa, cc)
	bd := new(big.Float).Mul(bb, dd)
	bc := new(big.Float).Mul(bb, cc)
	ad := new(big.Float).Mul(aa, dd)
	f.A.Quo(new(big.Float).Add(ac, bd), denom)
	f.B.Quo(new(big.Float).Sub(bc, ad), denom)
	return f
}

// Sqrt computes the square root of the complex number
// https://www.johndcook.com/blog/2020/06/09/complex-square-root/
func (f *Float) Sqrt(a *Float) *Float {
	x := big.NewFloat(0).SetPrec(f.A.Prec())
	y := big.NewFloat(0).SetPrec(f.B.Prec())
	l := big.NewFloat(0).SetPrec(f.A.Prec())
	x.Mul(a.A, a.A)
	y.Mul(a.B, a.B)
	l.Add(x, y)
	l = bigfloat.Sqrt(l)

	aa := big.NewFloat(0).SetPrec(f.A.Prec())
	aa.Add(l, a.A)
	aa.Quo(aa, big.NewFloat(2).SetPrec(f.A.Prec()))
	aa = bigfloat.Sqrt(aa)

	f.B.Sub(l, a.A)
	f.B.Quo(f.B, big.NewFloat(2).SetPrec(f.B.Prec()))
	f.B = bigfloat.Sqrt(f.B)
	f.B.Mul(big.NewFloat(float64(a.B.Sign())).SetPrec(f.B.Prec()), f.B)
	f.A = aa

	return f
}

// Atan2 computes atan2(y, x) for the complex value x+yi (same as Arg).
// https://en.wikipedia.org/wiki/Atan2
func (f *Float) Atan2(x *Float) *Float {
	return f.Arg(x)
}

// safeArctan computes arctan(t). bigfloat.Arctan does not converge for |t| ≥ 1,
// so reduce via arctan(t) = sign(t)·π/2 − arctan(1/t) when |t| > 1.
func safeArctan(t *big.Float) *big.Float {
	prec := t.Prec()
	if prec == 0 {
		prec = 64
	}
	abs := new(big.Float).Abs(t)
	one := big.NewFloat(1).SetPrec(prec)
	cmp := abs.Cmp(one)
	if cmp == 0 {
		// ±π/4
		q := new(big.Float).Quo(bigfloat.PI(prec), big.NewFloat(4).SetPrec(prec))
		if t.Sign() < 0 {
			q.Neg(q)
		}
		return q
	}
	if cmp > 0 {
		// sign(t)*π/2 − arctan(1/t)
		inv := new(big.Float).Quo(one, t)
		at := bigfloat.Arctan(inv)
		halfPi := new(big.Float).Quo(bigfloat.PI(prec), big.NewFloat(2).SetPrec(prec))
		if t.Sign() < 0 {
			halfPi.Neg(halfPi)
		}
		return new(big.Float).Sub(halfPi, at)
	}
	return bigfloat.Arctan(t)
}

// Arg computes arg(x + yi) using atan2(y, x) in (−π, π].
// https://mathworld.wolfram.com/ComplexArgument.html
func (f *Float) Arg(x *Float) *Float {
	a := x.A
	b := x.B
	prec := a.Prec()
	if b.Prec() > prec {
		prec = b.Prec()
	}
	if prec == 0 {
		prec = 64
	}
	f.B = big.NewFloat(0).SetPrec(prec)
	zero := big.NewFloat(0).SetPrec(prec)
	pi := bigfloat.PI(prec)
	halfPi := new(big.Float).Quo(new(big.Float).Copy(pi), big.NewFloat(2).SetPrec(prec))

	cmpA := a.Cmp(zero)
	cmpB := b.Cmp(zero)

	switch {
	case cmpA > 0:
		// atan(y/x)
		ratio := new(big.Float).SetPrec(prec).Quo(b, a)
		f.A.Set(safeArctan(ratio))
	case cmpA < 0 && cmpB >= 0:
		// atan(y/x) + π
		ratio := new(big.Float).SetPrec(prec).Quo(b, a)
		f.A.Set(safeArctan(ratio))
		f.A.Add(f.A, pi)
	case cmpA < 0 && cmpB < 0:
		// atan(y/x) - π
		ratio := new(big.Float).SetPrec(prec).Quo(b, a)
		f.A.Set(safeArctan(ratio))
		f.A.Sub(f.A, pi)
	case cmpA == 0 && cmpB > 0:
		f.A.Set(halfPi)
	case cmpA == 0 && cmpB < 0:
		f.A.Neg(halfPi)
	default:
		// 0+0i — undefined; return +Inf by convention
		f.A.SetInf(false)
	}
	return f
}

// Exp computes e^x for a complex number
// https://www.wolframalpha.com/input/?i=e%5E%28x+%2B+yi%29
func (f *Float) Exp(x *Float) *Float {
	exp := bigfloat.Exp(x.A)
	cos := bigfloat.Cos(x.B)
	sin := bigfloat.Sin(x.B)
	f.A.Mul(exp, cos)
	f.B.Mul(exp, sin)
	return f
}

// Cos computes cosine of a number
// https://www.wolframalpha.com/input/?i=cos%28x+%2B+yi%29
func (f *Float) Cos(x *Float) *Float {
	y1 := big.NewFloat(0).SetPrec(x.B.Prec())
	y1.Set(x.B)
	y1.Neg(y1)
	x1 := big.NewFloat(0).SetPrec(x.A.Prec())
	x1.Set(x.A)
	a := NewFloat(y1, x1)

	y2 := big.NewFloat(0).SetPrec(x.B.Prec())
	y2.Set(x.B)
	x2 := big.NewFloat(0).SetPrec(x.A.Prec())
	x2.Set(x.A)
	x2.Neg(x2)
	b := NewFloat(y2, x2)

	a.Exp(a)
	b.Exp(b)
	a.Add(a, b)
	x3 := NewFloat(big.NewFloat(.5).SetPrec(x.A.Prec()), big.NewFloat(0).SetPrec(x.B.Prec()))
	f.Mul(a, x3)
	return f
}

// Sin computes sine of a number
// https://www.wolframalpha.com/input/?i=sin%28x+%2B+yi%29
func (f *Float) Sin(x *Float) *Float {
	y1 := big.NewFloat(0).SetPrec(x.B.Prec())
	y1.Set(x.B)
	x1 := big.NewFloat(0).SetPrec(x.A.Prec())
	x1.Set(x.A)
	x1.Neg(x1)
	a := NewFloat(y1, x1)

	y2 := big.NewFloat(0).SetPrec(x.B.Prec())
	y2.Set(x.B)
	y2.Neg(y2)
	x2 := big.NewFloat(0).SetPrec(x.A.Prec())
	x2.Set(x.A)
	b := NewFloat(y2, x2)

	a.Exp(a)
	b.Exp(b)
	a.Sub(a, b)
	x3 := NewFloat(big.NewFloat(0).SetPrec(x.A.Prec()), big.NewFloat(.5).SetPrec(x.B.Prec()))
	f.Mul(a, x3)
	return f
}

// Tan computes tangent of a number
// https://en.wikipedia.org/wiki/Trigonometric_functions
func (f *Float) Tan(x *Float) *Float {
	a := big.NewFloat(0).SetPrec(x.A.Prec())
	a.Set(x.A)
	b := big.NewFloat(0).SetPrec(x.B.Prec())
	b.Set(x.B)
	y := NewFloat(a, b)
	x.Sin(x)
	y.Cos(y)
	f.Div(x, y)
	return f
}

// Log computes the natural log of x: log|z| + i·arg(z)
// https://en.wikipedia.org/wiki/Complex_logarithm
func (f *Float) Log(x *Float) *Float {
	a := x.A
	b := x.B
	prec := a.Prec()
	if b.Prec() > prec {
		prec = b.Prec()
	}
	aa := new(big.Float).Mul(a, a)
	bb := new(big.Float).Mul(b, b)
	mod := bigfloat.Sqrt(new(big.Float).Add(aa, bb))
	real := bigfloat.Log(mod)
	arg := NewFloat(big.NewFloat(0).SetPrec(prec), big.NewFloat(0).SetPrec(prec))
	arg.Arg(x)
	f.A.Set(real)
	f.B.Set(arg.A)
	return f
}

// Pow computes x**y
// https://mathworld.wolfram.com/ComplexExponentiation.html
func (f *Float) Pow(x *Float, y *Float) *Float {
	if x.A.Cmp(big.NewFloat(0).SetPrec(x.A.Prec())) == 0 &&
		x.B.Cmp(big.NewFloat(0).SetPrec(x.B.Prec())) == 0 &&
		y.A.Cmp(big.NewFloat(0).SetPrec(y.A.Prec())) == 0 &&
		y.B.Cmp(big.NewFloat(0).SetPrec(y.B.Prec())) == 0 {
		f.A.SetInf(false)
		return f
	}

	a := big.NewFloat(0).SetPrec(x.A.Prec())
	a.Set(x.A)
	aa := big.NewFloat(0).SetPrec(x.A.Prec())
	aa.Mul(a, a)
	b := big.NewFloat(0).SetPrec(x.B.Prec())
	b.Set(x.B)
	bb := big.NewFloat(0).SetPrec(x.B.Prec())
	bb.Mul(b, b)
	sum := big.NewFloat(0).SetPrec(x.A.Prec())
	sum.Add(aa, bb)
	c := big.NewFloat(0).SetPrec(y.A.Prec())
	c.Set(y.A)
	cc := big.NewFloat(0).SetPrec(c.Prec())
	cc.Quo(c, big.NewFloat(2).SetPrec(c.Prec()))
	e := bigfloat.Pow(sum, cc)
	d := big.NewFloat(0).SetPrec(y.B.Prec())
	d.Set(y.B)
	arg := NewFloat(big.NewFloat(0).SetPrec(a.Prec()),
		big.NewFloat(0).SetPrec(b.Prec()))
	arg.Arg(x)
	exp := big.NewFloat(0).SetPrec(arg.A.Prec())
	exp.Mul(d, arg.A)
	exp.Neg(exp)
	exp = bigfloat.Exp(exp)
	e.Mul(e, exp)

	i := big.NewFloat(0).SetPrec(c.Prec())
	i.Mul(c, arg.A)
	j := big.NewFloat(0).SetPrec(d.Prec())
	j.Mul(d, bigfloat.Log(sum))
	j.Quo(j, big.NewFloat(2).SetPrec(d.Prec()))
	i.Add(i, j)
	cos := bigfloat.Cos(i)
	cos.Mul(e, cos)
	f.A.Set(cos)
	sin := bigfloat.Sin(i)
	sin.Mul(e, sin)
	f.B.Set(sin)
	return f
}

// SetRat sets the value to a rational
func (f *Float) SetRat(r *Rational) {
	f.A.SetRat(r.A)
	f.B.SetRat(r.B)
}

// Rat stores the float in a rational
func (f *Float) Rat(r *Rational) {
	f.A.Rat(r.A)
	f.B.Rat(r.B)
}

// isTiny reports whether x is negligible for display (underflow noise).
func isTiny(x *big.Float) bool {
	if x.Sign() == 0 {
		return true
	}
	abs := new(big.Float).Abs(x)
	// Treat magnitudes smaller than 1e-20 as zero for pretty-printing.
	return abs.Cmp(big.NewFloat(1e-20)) < 0
}

// String returns a string of the complex number.
func (f *Float) String() string {
	reTiny, imTiny := isTiny(f.A), isTiny(f.B)
	if imTiny {
		if reTiny {
			return "0"
		}
		return f.A.String()
	}
	if reTiny {
		if f.B.Cmp(big.NewFloat(1)) == 0 {
			return "i"
		}
		if f.B.Cmp(big.NewFloat(-1)) == 0 {
			return "-i"
		}
		return f.B.String() + "i"
	}
	// Format as a ± bi
	bAbs := new(big.Float).Abs(f.B)
	sign := " + "
	if f.B.Sign() < 0 {
		sign = " - "
	}
	if bAbs.Cmp(big.NewFloat(1)) == 0 {
		return f.A.String() + sign + "i"
	}
	return f.A.String() + sign + bAbs.String() + "i"
}

// Rational is an imaginary number
type Rational struct {
	A *big.Rat // Real part
	B *big.Rat // Imaginary part
}

// NewRational creates a new imaginary number
func NewRational(a, b *big.Rat) *Rational {
	return &Rational{
		A: a,
		B: b,
	}
}

// Add add two imaginary numbers
func (r *Rational) Add(a, b *Rational) *Rational {
	r.A.Add(a.A, b.A)
	r.B.Add(a.B, b.B)
	return r
}

// Sub subtracts two imaginary numbers
func (r *Rational) Sub(a, b *Rational) *Rational {
	r.A.Sub(a.A, b.A)
	r.B.Sub(a.B, b.B)
	return r
}

// Mul multiples two imaginary numbers
func (r *Rational) Mul(a, b *Rational) *Rational {
	x1, x2, x3, x4 :=
		big.NewRat(0, 1), big.NewRat(0, 1),
		big.NewRat(0, 1), big.NewRat(0, 1)
	x1.Mul(a.A, b.A) // a*a
	x2.Mul(a.A, b.B) // a*ib
	x3.Mul(a.B, b.A) // ib*a
	x4.Mul(a.B, b.B) // i^2 * b = -b
	r.A.Add(x1, x4.Neg(x4))
	r.B.Add(x2, x3)
	return r
}

// Conj computes the complex conjugate of a
func (r *Rational) Conj(a *Rational) *Rational {
	r.A.Set(a.A)
	r.B.Neg(a.B)
	return r
}

// Div divides two imaginary numbers: (a+bi)/(c+di) = ((ac+bd)+(bc-ad)i)/(c²+d²).
// Copies inputs so callers may pass overlapping or shared rationals safely.
func (r *Rational) Div(a, b *Rational) *Rational {
	aa, bb := new(big.Rat).Set(a.A), new(big.Rat).Set(a.B)
	cc, dd := new(big.Rat).Set(b.A), new(big.Rat).Set(b.B)
	// |b|² = c² + d²
	denom := new(big.Rat).Add(new(big.Rat).Mul(cc, cc), new(big.Rat).Mul(dd, dd))
	// real = (ac + bd) / |b|²
	// imag = (bc - ad) / |b|²
	ac := new(big.Rat).Mul(aa, cc)
	bd := new(big.Rat).Mul(bb, dd)
	bc := new(big.Rat).Mul(bb, cc)
	ad := new(big.Rat).Mul(aa, dd)
	r.A.Quo(new(big.Rat).Add(ac, bd), denom)
	r.B.Quo(new(big.Rat).Sub(bc, ad), denom)
	return r
}

// Neg negates the rational
func (r *Rational) Neg(a *Rational) *Rational {
	r.A.Neg(a.A)
	r.B.Neg(a.B)
	return r
}

// String returns a string of the complex rational.
func (r *Rational) String() string {
	zero := big.NewRat(0, 1)
	if r.B.Cmp(zero) == 0 {
		return r.A.String()
	}
	if r.A.Cmp(zero) == 0 {
		if r.B.Cmp(big.NewRat(1, 1)) == 0 {
			return "i"
		}
		if r.B.Cmp(big.NewRat(-1, 1)) == 0 {
			return "-i"
		}
		return r.B.String() + "i"
	}
	bAbs := new(big.Rat).Abs(r.B)
	sign := " + "
	if r.B.Sign() < 0 {
		sign = " - "
	}
	if bAbs.Cmp(big.NewRat(1, 1)) == 0 {
		return r.A.String() + sign + "i"
	}
	return r.A.String() + sign + bAbs.String() + "i"
}
