package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math"
	"math/big"
	"strings"
)

type Token struct {
	Address   string `json:"address"`
	Precision int    `json:"precision"`
}

func LastElement[T any](arr []T) (T, error) {
	if len(arr) == 0 {
		var zero T
		return zero, fmt.Errorf("array is empty")
	}
	return arr[len(arr)-1], nil
}

// FindMax finds the maximum value in a map of slices.
func FindMax[K comparable, T any](candidate map[K][]T, comparator func(T, T) int) (T, K) {
	var maxVal T
	var maxKey K
	first := true

	for key, slice := range candidate {
		for _, val := range slice {
			if first {
				maxVal = val
				maxKey = key
				first = false
			} else if comparator(val, maxVal) > 0 {
				maxVal = val
				maxKey = key
			}
		}
	}

	return maxVal, maxKey
}

// FindMin finds the minimum value in a map of slices.
func FindMin[K comparable, T any](candidate map[K][]T, comparator func(T, T) int) (T, K) {
	var minVal T
	var minKey K
	first := true

	for key, slice := range candidate {
		for _, val := range slice {
			if first {
				minVal = val
				minKey = key
				first = false
			} else if comparator(val, minVal) < 0 {
				minVal = val
				minKey = key
			}
		}
	}

	return minVal, minKey
}

// ReverseSlice reverses a slice of any type.
func ReverseSlice[T any](slice []T) []T {
	reversed := make([]T, len(slice))
	for i, j := 0, len(slice)-1; i <= j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = slice[j], slice[i]
	}
	return reversed
}

func TickToWord(tick int, tickSpacing int) int {
	compressed := math.Round(float64(tick / tickSpacing))
	if tick < 0 && tick%tickSpacing != 0 {
		compressed -= 1
	}
	return tick >> 8
}

func BigIntToReadable(number *big.Int) float64 {
	result, _ := new(big.Float).Quo(new(big.Float).SetInt(number), big.NewFloat(1e18)).Float64()
	return result
}

func StringToBigInt(amount string) (res *big.Int) {
	if strings.HasPrefix(amount, "f") {
		bigInt := new(big.Int)
		bigInt.SetString(amount, 16)

		maxValue := new(big.Int).Lsh(big.NewInt(1), uint(len(amount)*4))
		maxValue.Sub(maxValue, big.NewInt(1))

		bigInt.Sub(maxValue, bigInt)
		bigInt.Add(bigInt, big.NewInt(1))
		bigInt.Neg(bigInt)
		res = bigInt
	} else {
		res, _ = new(big.Int).SetString(amount, 16)
	}
	return
}

func IsAddressInSlice(address common.Address, slice []string) bool {
	for _, addr := range slice {
		if address.Hex() == addr {
			return true
		}
	}
	return false
}

func HasTopics(log types.Log, topics ...string) bool {
	if len(log.Topics) == 0 {
		return false
	}

	actualTopic := log.Topics[0].Hex()
	for _, t := range topics {
		if actualTopic == t {
			return true
		}
	}

	return false
}

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func PercentageDiff(a, b *big.Int) float64 {
	// Calculate the absolute difference
	diff := new(big.Int).Sub(a, b)
	diff.Abs(diff)

	// Calculate the average of the two numbers
	sum := new(big.Int).Add(a, b)
	average := new(big.Float).Quo(new(big.Float).SetInt(sum), big.NewFloat(2))

	// Calculate the percentage difference
	diffFloat := new(big.Float).SetInt(diff)
	percentageDifference := new(big.Float).Quo(diffFloat, average)
	percentageDifference.Mul(percentageDifference, big.NewFloat(100))

	res, _ := percentageDifference.Float64()
	return res
}

// LogBigInt computes the natural logarithm (ln) of a *big.Int and returns a *big.Float.
// It handles large integers by taking square roots until the value can be safely converted to float64.
func LogBigInt(n *big.Int) *big.Float {
	if n.Sign() <= 0 {
		panic("Input must be a positive integer")
	}

	// Copy n to avoid modifying the original number
	x := new(big.Int).Set(n)

	// Initialize the multiplier to adjust the logarithm
	numMul := 1

	// Attempt to convert x to float64
	for {
		// Check if x can be safely converted to float64
		f64, acc := new(big.Float).SetInt(x).Float64()
		if acc == big.Exact || acc == big.Below {
			if !math.IsInf(f64, 0) && !math.IsNaN(f64) {
				// Compute ln(x) = ln(f64) * numMul
				logFloat := math.Log(f64) * float64(numMul)
				return big.NewFloat(logFloat)
			}
		}

		// Take the square root of x
		x.Sqrt(x)
		numMul *= 2
	}
}

//todo SIMULATE FROM THE SIMCALLRESULT

func Assert(condition bool, msg interface{}) {
	var stringMsg string

	if !condition {
		switch v := msg.(type) {
		case string:
			stringMsg = v
		/*
			case simulate.SimCallResult:
				stringMsg = v.Error.Message
		*/
		case error:
			stringMsg = v.Error()
		}
		panic(stringMsg)
	}
}