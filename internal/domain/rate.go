package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// RateSnapshot is a single record persisted for each GetRates call
type RateSnapshot struct {
	// ID auto-generated primary key
	ID int64
	// AskPrice best ask price
	AskPrice decimal.Decimal
	// BidPrice best bid price
	BidPrice decimal.Decimal
	// TopN ask price at the N-th position
	TopN decimal.Decimal
	// AvgNM arithmetic mean of ask prices from N to M
	AvgNM decimal.Decimal
	// FetchedAt time the rates was fetched from Grinex
	FetchedAt time.Time
}
