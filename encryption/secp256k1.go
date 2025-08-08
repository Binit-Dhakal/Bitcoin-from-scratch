package encryption

import (
	"fmt"
	"math/big"
)

type Secp256k1 struct {
	n    *big.Int
	p    *big.Int
	a, b *FieldElement
}

func NewSecp256k1() *Secp256k1 {
	twoExp256 := new(big.Int).Lsh(big.NewInt(1), 256)
	twoExp32 := new(big.Int).Lsh(big.NewInt(1), 32)
	p := new(big.Int).Sub(twoExp256, twoExp32)
	p.Sub(p, big.NewInt(977))

	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	a := NewFieldElement(big.NewInt(0), p)
	b := NewFieldElement(big.NewInt(7), p)

	return &Secp256k1{
		p: p,
		n: n,
		a: a,
		b: b,
	}
}

func (s *Secp256k1) Prime() *big.Int {
	return s.p
}

func (s *Secp256k1) Order() *big.Int {
	return s.n
}

func (s *Secp256k1) A() *FieldElement {
	return s.a
}

func (s *Secp256k1) B() *FieldElement {
	return s.b
}

func (s *Secp256k1) NewFieldElement(num *big.Int) *FieldElement {
	return NewFieldElement(num, s.p)
}

func (s *Secp256k1) NewPoint(x *big.Int, y *big.Int) (*Point, error) {
	feX := s.NewFieldElement(x)
	feY := s.NewFieldElement(y)
	if x != nil {
		p, err := NewPoint(s.a, s.b, feX, feY)
		if err != nil {
			return nil, err
		}

		return p, nil
	}

	// in case of identity point
	return NewPoint(s.a, s.b, nil, nil)
}

func (s *Secp256k1) ScalarMul(point *Point, scalar *big.Int) (*Point, error) {
	if !point.a.Equal(s.a) || !point.b.Equal(s.b) {
		return nil, fmt.Errorf("point is not in the curve")
	}
	scalar.Mod(scalar, s.n)
	return point.ScalarMul(scalar)
}

func (s *Secp256k1) Add(p1, p2 *Point) (*Point, error) {
	if !p1.a.Equal(s.a) || !p1.b.Equal(s.b) || !p2.a.Equal(s.a) || !p2.b.Equal(s.b) {
		return nil, fmt.Errorf("one or both points are not in the curve")
	}

	return p1.Add(p2)
}
