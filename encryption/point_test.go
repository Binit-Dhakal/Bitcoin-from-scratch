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
		parameterA [4]int //a,b,x,y // 0 in x, y mean nil value
		parameterB [4]int
		output     [4]int
	}{
		{"Identity Case: A+I", [4]int{5, 7, -1, -1}, [4]int{5, 7, 0, 0}, [4]int{5, 7, -1, -1}},
		{"Identity Case: I+A", [4]int{5, 7, 0, 0}, [4]int{5, 7, -1, -1}, [4]int{5, 7, -1, -1}},
		{"Vertical Line: slope inf", [4]int{5, 7, -1, -1}, [4]int{5, 7, -1, 1}, [4]int{5, 7, 0, 0}},
		{"different point: x1 != x2", [4]int{5, 7, 2, 5}, [4]int{5, 7, -1, -1}, [4]int{5, 7, 3, -7}},
		{"same point: x1==x2", [4]int{5, 7, -1, -1}, [4]int{5, 7, -1, -1}, [4]int{5, 7, 18, 77}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			p1 := BuildPointStructFromArray(test.parameterA)
			p2 := BuildPointStructFromArray(test.parameterB)

			res, err := p1.Add(p2)

			require.NoError(t, err)

			arr := BuildArrayFromPointStruct(res)

			assert.Equal(t, test.output, arr)
		})
	}
}

func convertToBigInt(a int) *big.Int {
	return big.NewInt(int64(a))
}

func BuildPointStructFromArray(parameters [4]int) *Point {
	p := [4]*big.Int{}
	for i, parameter := range parameters {
		if parameter == 0 {
			p[i] = nil
		} else {
			p[i] = convertToBigInt(parameter)
		}
	}
	point, err := NewPoint(p[0], p[1], p[2], p[3])
	if err != nil {
		// panic is ok here as this function is just used for testing
		panic(err)
	}

	return point
}

func convertBigIntToInt(a *big.Int) int {
	if a == nil {
		return 0
	}

	return int(a.Int64())
}

func BuildArrayFromPointStruct(point *Point) [4]int {
	parameters := [4]int{}

	parameters[0] = convertBigIntToInt(point.a)
	parameters[1] = convertBigIntToInt(point.b)
	parameters[2] = convertBigIntToInt(point.x)
	parameters[3] = convertBigIntToInt(point.y)

	return parameters
}
