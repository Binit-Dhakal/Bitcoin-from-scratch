package encryption

import (
	"errors"
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
	a *big.Int
	b *big.Int
	x *big.Int
	y *big.Int
}

func OpOnBig(x *big.Int, y *big.Int, optype OP_TYPE) *big.Int {
	var op big.Int
	switch optype {
	case ADD:
		return op.Add(x, y)
	case SUB:
		return op.Sub(x, y)
	case MUL:
		return op.Mul(x, y)
	case DIV:
		return op.Div(x, y)
	case EXP:
		return op.Exp(x, y, nil)
	}

	panic("wrong operation type")
}

func NewPoint(a *big.Int, b *big.Int, x *big.Int, y *big.Int) (*Point, error) {
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
		return nil, InvalidPoint
	}

	return &p, nil
}

func isValid(a *big.Int, b *big.Int, x *big.Int, y *big.Int) bool {
	// y^2 = x^3 + a*x + b
	lhs := OpOnBig(y, big.NewInt(int64(2)), EXP)
	rhs := OpOnBig(OpOnBig(OpOnBig(x, big.NewInt(int64(3)), EXP), OpOnBig(a, x, MUL), ADD), b, ADD)
	return lhs.Cmp(rhs) == 0
}

func (p *Point) Equal(other *Point) bool {
	return p.a.Cmp(other.a) == 0 && p.x.Cmp(other.x) == 0 && p.b.Cmp(other.b) == 0 && p.y.Cmp(other.y) == 0
}

func (p *Point) Add(other *Point) (*Point, error) {
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 {
		return nil, InvalidPoint
	}

	if p.x == nil {
		return other, nil
	}

	if other.x == nil {
		return p, nil
	}

	// CASE: vertical line
	if p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) != 0 {
		return &Point{a: p.a, b: p.b, x: nil, y: nil}, nil
	}

	// CASE: vertical line at tip
	if p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) == 0 && p.y.Cmp(big.NewInt(int64(0))) == 0 {
		return &Point{a: p.a, b: p.b, x: nil, y: nil}, nil
	}

	var s *big.Int // slope
	// CASE: x1 != x2
	// CASE: x1 == x2 and also y1 == y2
	if p.x.Cmp(other.x) != 0 {
		// slope = (y_2 - y_1)/(x_2 - x_1)
		num := OpOnBig(other.y, p.y, SUB)
		den := OpOnBig(other.x, p.x, SUB)
		s = OpOnBig(num, den, DIV)
	} else {
		// slope = (3 * x * x + a)/(2 * y)
		num := OpOnBig(OpOnBig(big.NewInt(int64(3)), OpOnBig(p.x, big.NewInt(int64(2)), EXP), MUL), p.a, ADD)
		den := OpOnBig(big.NewInt(int64(2)), p.y, MUL)
		s = OpOnBig(num, den, DIV)
	}

	// x_3 = s^2 - x_1 - x_2
	// y_3 = s*(x_1 - x_3) - y_1
	x_3 := OpOnBig(OpOnBig(OpOnBig(s, big.NewInt(int64(2)), EXP), p.x, SUB), other.x, SUB)
	y_3 := OpOnBig(OpOnBig(OpOnBig(p.x, x_3, SUB), s, MUL), p.y, SUB)

	return &Point{x: x_3, y: y_3, a: p.a, b: p.b}, nil
}
