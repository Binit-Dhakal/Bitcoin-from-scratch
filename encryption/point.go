package encryption

import (
	"errors"
	"fmt"
	"math/big"
)

var InvalidPoint = errors.New("Point is not in the curve")

// curve of the form {y^2 = x^3+ax+b}
type Point struct {
	a *FieldElement
	b *FieldElement
	x *FieldElement
	y *FieldElement
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
		return nil, InvalidPoint
	}

	return &p, nil
}

func isValid(a *FieldElement, b *FieldElement, x *FieldElement, y *FieldElement) bool {
	if !a.CheckOrder(b) || !b.CheckOrder(x) || !x.CheckOrder(y) {
		return false
	}

	// y^2 = x^3 + a*x + b
	lhs := y.Exponent(big.NewInt(2))

	x3 := x.Exponent(big.NewInt(3))
	ax := x.Mul(a)

	rhs := x3.Add(ax)
	rhs = rhs.Add(b)

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
		num := other.y.Sub(p.y)
		den := other.x.Sub(p.x)
		s = num.Division(den)
	} else {
		// slope = (3 * x * x + a)/(2 * y)
		x2 := p.x.Exponent(big.NewInt(2))
		num := x2.ScalarMul(big.NewInt(3))
		num = num.Add(p.a)

		den := p.y.ScalarMul(big.NewInt(2))
		s = num.Division(den)
	}

	// x_3 = s^2 - x_1 - x_2
	// y_3 = s*(x_1 - x_3) - y_1
	s2 := s.Exponent(big.NewInt(2))
	x_3 := s2.Sub(p.x).Sub(other.x)
	y_3 := p.x.Sub(x_3).Mul(s).Sub(p.y)

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
