package config

import "math/big"

type FeeAmount uint64

const (
	FeeMax FeeAmount = 1000000
)

var (
	Zero        = big.NewInt(0)
	BigAmountIn = big.NewInt(1000000000000000000)
)

var TickSpacings = map[FeeAmount]int{
	500:   10,
	3000:  60,
	10000: 200,
}
