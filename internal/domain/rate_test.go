package domain_test

import (
	"testing"
	"time"

	"github.com/Danil-Ivonin/GrinexRates/internal/domain"
	"github.com/shopspring/decimal"
)

func TestRateSnapshotFieldsExist(t *testing.T) {
	fetchedTime := time.Now()
	rs := domain.RateSnapshot{
		ID:        42,
		AskPrice:  decimal.RequireFromString("100.50"),
		BidPrice:  decimal.RequireFromString("100.25"),
		TopN:      decimal.RequireFromString("100.75"),
		AvgNM:     decimal.RequireFromString("100.40"),
		FetchedAt: fetchedTime,
	}

	if rs.ID != 42 {
		t.Errorf("ID: got %d, want %d", rs.ID, 42)
	}
	if !rs.AskPrice.Equal(decimal.RequireFromString("100.50")) {
		t.Errorf("AskPrice: got %v, want 100.50", rs.AskPrice)
	}
	if !rs.BidPrice.Equal(decimal.RequireFromString("100.25")) {
		t.Errorf("BidPrice: got %v, want 100.25", rs.BidPrice)
	}
	if !rs.TopN.Equal(decimal.RequireFromString("100.75")) {
		t.Errorf("TopN: got %v, want 100.75", rs.TopN)
	}
	if !rs.AvgNM.Equal(decimal.RequireFromString("100.40")) {
		t.Errorf("AvgNM: got %v, want 100.40", rs.AvgNM)
	}
	if rs.FetchedAt != fetchedTime {
		t.Errorf("FetchedAt: got %v, want %v", rs.FetchedAt, fetchedTime)
	}
}

func TestRateSnapshotZeroValue(t *testing.T) {
	rs := domain.RateSnapshot{}

	if rs.ID != 0 {
		t.Errorf("ID zero value: got %d, want 0", rs.ID)
	}
	if !rs.AskPrice.IsZero() {
		t.Errorf("AskPrice zero value: got %v, want zero decimal", rs.AskPrice)
	}
	if !rs.BidPrice.IsZero() {
		t.Errorf("BidPrice zero value: got %v, want zero decimal", rs.BidPrice)
	}
	if !rs.TopN.IsZero() {
		t.Errorf("TopN zero value: got %v, want zero decimal", rs.TopN)
	}
	if !rs.AvgNM.IsZero() {
		t.Errorf("AvgNM zero value: got %v, want zero decimal", rs.AvgNM)
	}
	if !rs.FetchedAt.IsZero() {
		t.Errorf("FetchedAt zero value: got %v, want zero time", rs.FetchedAt)
	}
}
