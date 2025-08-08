package encryption

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecp256k1(t *testing.T) {
	// test G.n = I
	gx := new(big.Int)
	gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)

	gy := new(big.Int)
	gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	curve := NewSecp256k1()
	point, err := curve.NewPoint(gx, gy)
	require.NoError(t, err)

	res, err := curve.ScalarMul(point, curve.Order())
	require.NoError(t, err)

	x, y := res.GetCoordinates()
	assert.Nil(t, x)
	assert.Nil(t, y)
}
