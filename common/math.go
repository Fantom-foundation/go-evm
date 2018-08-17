package common

import (
	"math"
	"math/big"
)

// https://github.com/google/skylark/blob/0a5e39a3470a7db0846fb24f517d5a41cd344f64/int.go#L98
func bigintToUint64(i *big.Int) (uint64, big.Accuracy) {
	sign := i.Sign()
	if sign > 0 {
		if i.BitLen() > 64 {
			return math.MaxUint64, big.Below
		}
	} else if sign < 0 {
		return 0, big.Above
	}
	return i.Uint64(), big.Exact
}

func BigintToUInt64(i *big.Int) uint64 {
	num, _ := bigintToUint64(i)
	// TODO: Deal with cases where it's not `big.Exact`
	return num
}
