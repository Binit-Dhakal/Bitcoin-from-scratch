package encryption

import "math/big"

type Curve interface {
	Prime() *big.Int
	Order() *big.Int
	A() *FieldElement
	B() *FieldElement

	NewFieldElement(num *big.Int) *FieldElement
	NewPoint(x, y *big.Int) (*Point, error)

	ScalarMul(point *Point, scalar *big.Int) (*Point, error)
	Add(p1, p2 *Point) (*Point, error)
}
