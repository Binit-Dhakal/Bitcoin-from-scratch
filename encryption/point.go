package encryption

import (
	"errors"
	"fmt"
	"math/big"
)

type OP_TYPE int

const (
	ADD OP_TYPE = iota
	SUB
	MUL
	DIV
	EXP
)

var InvalidPoint = errors.New("Point is not in the curve")

// curve of the form {y^2 = x^3+ax+b}
type Point struct {
	a *FieldElement
	b *FieldElement
	x *FieldElement
	y *FieldElement
}

func OpOnBig(x *FieldElement, y *FieldElement, scalar *big.Int, optype OP_TYPE) *FieldElement {
	switch optype {
	case ADD:
		res, _ := x.Add(y)
		return res
	case SUB:
		res, _ := x.Sub(y)
		return res
	case MUL:
		if y != nil {
			res, _ := x.Mul(y)
			return res
		}

		if scalar != nil {
			res, _ := x.ScalarMul(scalar)
			return res
		}
		panic("wrong Multiplication operation parameters")
	case DIV:
		res, _ := x.Division(y)
		return res
	case EXP:
		if scalar == nil {
			panic("Scalar cannot be nil for EXP op")
		}
		res, _ := x.Exponent(scalar)
		return res
	}

	panic("wrong operation type")
}

func NewPoint(a *FieldElement, b *FieldElement, x *FieldElement, y *FieldElement) (*Point, error) {
	p := Point{
		a: a,
		b: b,
		x: x,
		y: y,
	}

	if x == nil && y == nil {
		// identity point and thus no need to check for validity
		return &p, nil
	}

	if !isValid(a, b, x, y) {
		fmt.Println("error")
		return nil, InvalidPoint
	}

	return &p, nil
}

func isValid(a *FieldElement, b *FieldElement, x *FieldElement, y *FieldElement) bool {
	// y^2 = x^3 + a*x + b
	lhs := OpOnBig(y, nil, big.NewInt(int64(2)), EXP)
	x3 := OpOnBig(x, nil, big.NewInt(int64(3)), EXP)
	ax := OpOnBig(a, x, nil, MUL)
	rhs := OpOnBig(OpOnBig(x3, ax, nil, ADD), b, nil, ADD)
	return lhs.Equal(rhs)
}

func (p *Point) Equal(other *Point) bool {
	return p.a.Equal(other.a) && p.x.Equal(other.x) && p.b.Equal(other.b) && p.y.Equal(other.y)
}

func (p *Point) GetCoordinates() (*big.Int, *big.Int) {
	return p.x.num, p.y.num
}

func (p *Point) Copy() *Point {
	return &Point{
		a: p.a.Copy(),
		b: p.b.Copy(),
		x: p.x.Copy(),
		y: p.y.Copy(),
	}
}

func (p *Point) Add(other *Point) (*Point, error) {
	if !p.a.Equal(other.a) || !p.b.Equal(other.b) {
		return nil, InvalidPoint
	}

	if p.x == nil {
		return other, nil
	}

	if other.x == nil {
		return p, nil
	}

	// CASE: vertical line
	if p.x.Equal(other.x) && !p.y.Equal(other.y) {
		return &Point{a: p.a, b: p.b, x: nil, y: nil}, nil
	}

	// CASE: vertical line at tip
	zero := NewFieldElement(big.NewInt(int64(0)), p.x.prime)
	if p.x.Equal(other.x) && p.y.Equal(other.y) && p.y.Equal(zero) {
		return &Point{a: p.a, b: p.b, x: nil, y: nil}, nil
	}

	var s *FieldElement // slope
	// CASE: x1 != x2
	// CASE: x1 == x2 and also y1 == y2
	if !p.x.Equal(other.x) {
		// slope = (y_2 - y_1)/(x_2 - x_1)
		num := OpOnBig(other.y, p.y, nil, SUB)
		den := OpOnBig(other.x, p.x, nil, SUB)
		s = OpOnBig(num, den, nil, DIV)
	} else {
		// slope = (3 * x * x + a)/(2 * y)
		x2 := OpOnBig(p.x, nil, big.NewInt(int64(2)), EXP)
		num := OpOnBig(OpOnBig(x2, nil, big.NewInt(int64(3)), MUL), p.a, nil, ADD)
		den := OpOnBig(p.y, nil, big.NewInt(int64(2)), MUL)
		s = OpOnBig(num, den, nil, DIV)
	}

	// x_3 = s^2 - x_1 - x_2
	// y_3 = s*(x_1 - x_3) - y_1
	s2 := OpOnBig(s, nil, big.NewInt(int64(2)), EXP)
	x_3 := OpOnBig(OpOnBig(s2, p.x, nil, SUB), other.x, nil, SUB)
	y_3 := OpOnBig(OpOnBig(OpOnBig(p.x, x_3, nil, SUB), s, nil, MUL), p.y, nil, SUB)

	return &Point{x: x_3, y: y_3, a: p.a, b: p.b}, nil
}

func (p *Point) ScalarMul(scalar *big.Int) (*Point, error) {
	if scalar == nil {
		return nil, fmt.Errorf("insufficient arguments")
	}

	result := &Point{a: p.a, b: p.b, x: nil, y: nil}
	current := p.Copy()

	coef := new(big.Int).Set(scalar)

	for coef.Sign() > 0 {
		if coef.Bit(0) == 1 {
			if result.x == nil {
				result = current.Copy()
			} else {
				result, _ = result.Add(current)
			}
		}

		current, _ = current.Add(current)
		coef.Rsh(coef, 1)
	}

	return result, nil
}
