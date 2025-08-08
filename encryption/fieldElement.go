package encryption

import (
	"errors"
	"fmt"
	"math/big"
)

var ErrOpInDifferentField = errors.New("Cannot do operation in different fields")

type FieldElement struct {
	num   *big.Int
	prime *big.Int
}

func NewFieldElement(num *big.Int, prime *big.Int) *FieldElement {
	// not checking for identity point
	if prime.Cmp(big.NewInt(int64(0))) != 0 && prime.Cmp(num) == -1 {
		panic(fmt.Errorf("Num not in range of 0 to %d", prime))
	}

	return &FieldElement{
		prime: prime,
		num:   num,
	}
}

func (f *FieldElement) Copy() *FieldElement {
	return &FieldElement{
		num:   new(big.Int).Set(f.num),
		prime: new(big.Int).Set(f.prime),
	}
}

func (f *FieldElement) Equal(other *FieldElement) bool {
	if other == nil {
		return false
	}

	return f.prime.Cmp(other.prime) == 0 && f.num.Cmp(other.num) == 0
}

func (f *FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", f.prime, f.num)
}

func (f *FieldElement) CheckOrder(other *FieldElement) bool {
	if f.prime.Cmp(other.prime) != 0 {
		return false
	}
	return true
}

func (f *FieldElement) Add(other *FieldElement) *FieldElement {
	var op big.Int
	res := op.Mod(op.Add(f.num, other.num), f.prime)

	return NewFieldElement(res, f.prime)
}

func (f *FieldElement) Sub(other *FieldElement) *FieldElement {
	var op big.Int
	res := op.Mod(op.Sub(f.num, other.num), f.prime)

	return NewFieldElement(res, f.prime)
}

func (f *FieldElement) Mul(other *FieldElement) *FieldElement {
	var op big.Int
	res := op.Mod(op.Mul(f.num, other.num), f.prime)

	return NewFieldElement(res, f.prime)
}

func (f *FieldElement) Exponent(power *big.Int) *FieldElement {
	var op big.Int
	op.Mod(power, op.Sub(f.prime, big.NewInt(int64(1))))

	res := op.Exp(f.num, &op, f.prime)

	return NewFieldElement(res, f.prime)
}

func (f *FieldElement) ScalarMul(val *big.Int) *FieldElement {
	var op big.Int
	res := op.Mul(f.num, val)
	res = op.Mod(res, f.prime)
	return NewFieldElement(res, f.prime)
}

func (f *FieldElement) Division(other *FieldElement) *FieldElement {
	// using Fermat's Little Theorem
	// (a/b) = a.b^(p-2)
	var power big.Int
	var op big.Int
	o := other.Exponent(power.Sub(other.prime, big.NewInt(int64(2))))
	t := op.Mul(f.num, o.num)
	op.Mod(t, f.prime)

	return NewFieldElement(&op, f.prime)
}
