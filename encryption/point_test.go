package encryption

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointADD(t *testing.T) {
	testCases := []struct {
		name       string
		order      int
		parameterA [4]int //a,b,x,y // 0 in x, y mean nil value
		parameterB [4]int
		output     [4]int
	}{
		{"Identity Case: A+I", 223, [4]int{5, 7, -1, -1}, [4]int{5, 7, 0, 0}, [4]int{5, 7, -1, -1}},
		{"Identity Case: I+A", 223, [4]int{5, 7, 0, 0}, [4]int{5, 7, -1, -1}, [4]int{5, 7, -1, -1}},
		{"Vertical Line: slope inf", 223, [4]int{5, 7, -1, -1}, [4]int{5, 7, -1, 1}, [4]int{5, 7, 0, 0}},
		{"different point: x1 != x2", 223, [4]int{5, 7, 2, 5}, [4]int{5, 7, -1, -1}, [4]int{5, 7, 3, 216}},
		{"same point: x1==x2", 223, [4]int{5, 7, -1, -1}, [4]int{5, 7, -1, -1}, [4]int{5, 7, 18, 77}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			p1 := BuildPointStructFromArray(test.order, test.parameterA)
			p2 := BuildPointStructFromArray(test.order, test.parameterB)

			res, err := p1.Add(p2)

			require.NoError(t, err)

			arr := BuildArrayFromPointStruct(res)

			assert.Equal(t, test.output, arr)
		})
	}
}

func TestScalarMul(t *testing.T) {
	testCases := []struct {
		name       string
		order      int
		parameters [4]int
		scalar     int
		output     [4]int
	}{
		{"basic example-5", 223, [4]int{0, 7, 47, 71}, 5, [4]int{0, 7, 126, 96}},
		{"basic example-10", 223, [4]int{0, 7, 47, 71}, 10, [4]int{0, 7, 154, 150}},
		{"basic example-15", 223, [4]int{0, 7, 47, 71}, 15, [4]int{0, 7, 139, 86}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			p1 := BuildPointStructFromArray(test.order, test.parameters)

			res, err := p1.ScalarMul(big.NewInt(int64(test.scalar)))
			require.NoError(t, err)

			arr := BuildArrayFromPointStruct(res)
			assert.Equal(t, test.output, arr)
		})
	}

}

func convertToBigInt(order int, a int) *FieldElement {
	o := big.NewInt(int64(order))

	b := big.NewInt(int64(a))
	return NewFieldElement(b, o)
}

func BuildPointStructFromArray(order int, parameters [4]int) *Point {
	p := [4]*FieldElement{}
	for i, parameter := range parameters {
		if i > 1 && parameter == 0 {
			p[i] = nil
		} else {
			p[i] = convertToBigInt(order, parameter)
		}
	}
	point, err := NewPoint(p[0], p[1], p[2], p[3])
	if err != nil {
		// panic is ok here as this function is just used for testing
		panic(err)
	}

	return point
}

func convertBigIntToInt(a *FieldElement) int {
	if a == nil {
		return 0
	}

	return int(a.num.Int64())
}

func BuildArrayFromPointStruct(point *Point) [4]int {
	parameters := [4]int{}

	parameters[0] = convertBigIntToInt(point.a)
	parameters[1] = convertBigIntToInt(point.b)
	parameters[2] = convertBigIntToInt(point.x)
	parameters[3] = convertBigIntToInt(point.y)

	return parameters
}
