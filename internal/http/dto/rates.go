package dto

import "github.com/shopspring/decimal"

type StockRatesEntry struct {
	Price  decimal.Decimal `json:"price"`
	Volume decimal.Decimal `json:"volume"`
	Amount decimal.Decimal `json:"amount"`
}

type StockRates struct {
	Timestamp int64             `json:"timestamp"`
	Asks      []StockRatesEntry `json:"asks"`
	Bids      []StockRatesEntry `json:"bids"`
}
