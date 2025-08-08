package encryption

import (
	"math/big"
	"testing"
)

func TestFieldElementEqual(t *testing.T) {
	fe1 := NewFieldElement(big.NewInt(int64(4)), big.NewInt(int64(19)))
	fe2 := NewFieldElement(big.NewInt(int64(10)), big.NewInt(int64(19)))
	fe3 := NewFieldElement(big.NewInt(int64(10)), big.NewInt(int64(20)))

	if !fe1.Equal(fe1) {
		t.Errorf("%v is equal to %v", fe1, fe1)
	}

	if fe1.Equal(fe2) {
		t.Errorf("%v is not equal to %v", fe1, fe2)
	}

	if fe2.Equal(fe3) {
		t.Errorf("%v has different field than %v", fe1, fe2)
	}
}

func TestFieldElementOperations(t *testing.T) {
	testCases := []struct {
		name   string
		op     string
		fe1    [2]int
		fe2    [2]int
		output [2]int
		err    error
	}{
		{"less than p after add", "ADD", [2]int{4, 19}, [2]int{10, 19}, [2]int{14, 19}, nil},
		{"more than p after add", "ADD", [2]int{14, 19}, [2]int{10, 19}, [2]int{5, 19}, nil},
		{"different order", "ADD", [2]int{14, 19}, [2]int{14, 15}, [2]int{0, 0}, ErrOpInDifferentField},

		{"positive after sub", "SUB", [2]int{14, 19}, [2]int{10, 19}, [2]int{4, 19}, nil},
		{"negative after sub", "SUB", [2]int{10, 19}, [2]int{14, 19}, [2]int{15, 19}, nil},
		{"order different", "SUB", [2]int{14, 19}, [2]int{14, 15}, [2]int{0, 0}, ErrOpInDifferentField},

		{"normal multiplication", "MUL", [2]int{3, 13}, [2]int{12, 13}, [2]int{10, 13}, nil},
		{"different order", "MUL", [2]int{3, 13}, [2]int{3, 10}, [2]int{0, 0}, ErrOpInDifferentField},

		// fe2 first parameter is the exponent value
		{"normal exponentation", "EXP", [2]int{3, 13}, [2]int{3, 0}, [2]int{1, 13}, nil},
		{"negative power", "EXP", [2]int{3, 13}, [2]int{-16, 0}, [2]int{9, 13}, nil},

		{"division", "DIV", [2]int{2, 19}, [2]int{7, 19}, [2]int{3, 19}, nil},
		{"different order", "DIV", [2]int{2, 19}, [2]int{7, 10}, [2]int{0, 0}, ErrOpInDifferentField},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			fe1 := NewFieldElement(big.NewInt(int64(test.fe1[0])), big.NewInt(int64(test.fe1[1])))
			fe2 := NewFieldElement(big.NewInt(int64(test.fe2[0])), big.NewInt(int64(test.fe2[1])))
			output := NewFieldElement(big.NewInt(int64(test.output[0])), big.NewInt(int64(test.output[1])))
			if test.output[1] == 0 {
				output = nil
			}

			var res *FieldElement
			var err error

			switch test.op {
			case "ADD":
				res, err = fe1.Add(fe2)
			case "SUB":
				res, err = fe1.Sub(fe2)
			case "MUL":
				res, err = fe1.Mul(fe2)
			case "EXP":
				res, err = fe1.Exponent(big.NewInt(int64(test.fe2[0])))
			case "DIV":
				res, err = fe1.Division(fe2)
			}

			if err != test.err {
				t.Errorf("Error should have been %v", test.err)
			}

			if err == nil && !res.Equal(output) {
				t.Errorf("%v and %v are equal", res, output)
			}
		})
	}
}
