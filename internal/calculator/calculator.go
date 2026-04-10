// Package calculator implements TopN and AvgNM computations on order book price strings.
package calculator

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// TopN returns the ask price at position n-1 from the asks slice.
func TopN(asks []decimal.Decimal, n int) (decimal.Decimal, error) {
	if n < 1 {
		return decimal.Zero, fmt.Errorf("calculator: topN: n=%d must be >= 1", n)
	}
	if n > len(asks) {
		return decimal.Zero, fmt.Errorf("calculator: topN: n=%d out of bounds (len=%d)", n, len(asks))
	}
	return asks[n-1], nil
}

// AvgNM returns the arithmetic mean of ask prices from n to m (1-indexed, inclusive)
func AvgNM(asks []decimal.Decimal, n, m int, precision int32) (decimal.Decimal, error) {
	if n < 1 {
		return decimal.Zero, fmt.Errorf("calculator: avgNM: n=%d must be >= 1", n)
	}
	if m < n {
		return decimal.Zero, fmt.Errorf("calculator: avgNM: m=%d must be >= n=%d", m, n)
	}
	if m > len(asks) {
		return decimal.Zero, fmt.Errorf("calculator: avgNM: m=%d out of bounds (len=%d)", m, len(asks))
	}

	sum := decimal.Zero
	count := decimal.NewFromInt(int64(m - n + 1))

	for i := n - 1; i < m; i++ {
		sum = sum.Add(asks[i])
	}

	return sum.Div(count).Round(precision), nil
}
