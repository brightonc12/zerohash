package numbers

import (
	"fmt"
	"math/big"
)

func NewZeroBigFloat() *big.Float {
	float := new(big.Float)
	float.SetPrec(128)
	return float
}

func ParseBigFloat(str string) (*big.Float, error) {
	f := NewZeroBigFloat()
	f, _, err := f.Parse(str, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string: '%s' error: %v", str, err)
	}
	return f, nil
}

func NewBigFloat(num float64) *big.Float {
	float := new(big.Float)
	float.SetPrec(128)
	float.SetFloat64(num)
	return float
}
