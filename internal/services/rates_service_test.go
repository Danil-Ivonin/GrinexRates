package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Danil-Ivonin/GrinexRates/internal/domain"
	"github.com/Danil-Ivonin/GrinexRates/internal/http/dto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeFetcher struct {
	rates *dto.StockRates
	err   error
}

func (f *fakeFetcher) Fetch(_ context.Context) (*dto.StockRates, error) {
	return f.rates, f.err
}

type fakeRatesRepository struct {
	saved []domain.RateSnapshot
	err   error
}

func (f *fakeRatesRepository) Create(_ context.Context, snap domain.RateSnapshot) error {
	if f.err != nil {
		return f.err
	}
	f.saved = append(f.saved, snap)
	return nil
}

func validStockRates() *dto.StockRates {
	return &dto.StockRates{
		Timestamp: time.Now().Unix(),
		Asks: []dto.StockRatesEntry{
			{Price: decimal.NewFromFloat(97.50), Volume: decimal.NewFromInt(1), Amount: decimal.NewFromInt(1)},
			{Price: decimal.NewFromFloat(98.00), Volume: decimal.NewFromFloat(2.5), Amount: decimal.NewFromFloat(2.5)},
			{Price: decimal.NewFromFloat(99.10), Volume: decimal.NewFromFloat(0.5), Amount: decimal.NewFromFloat(0.5)},
		},
		Bids: []dto.StockRatesEntry{
			{Price: decimal.NewFromFloat(97.00), Volume: decimal.NewFromInt(3), Amount: decimal.NewFromInt(3)},
			{Price: decimal.NewFromFloat(96.50), Volume: decimal.NewFromInt(1), Amount: decimal.NewFromInt(1)},
		},
	}
}

func TestRatesService_GetRates(t *testing.T) {
	t.Parallel()

	const avgNMPrecision = 6

	tests := []struct {
		name            string
		fetcher         *fakeFetcher
		repo            *fakeRatesRepository
		n, m            int32
		wantErr         bool
		errSubstr       string
		wantAskPrice    decimal.Decimal
		wantBidPrice    decimal.Decimal
		checkSavedCount int
	}{
		{
			name:            "success with valid order book",
			fetcher:         &fakeFetcher{rates: validStockRates()},
			repo:            &fakeRatesRepository{},
			n:               1,
			m:               3,
			wantErr:         false,
			wantAskPrice:    decimal.NewFromFloat(97.50),
			wantBidPrice:    decimal.NewFromFloat(97.00),
			checkSavedCount: 1,
		},
		{
			name:      "fetcher returns error",
			fetcher:   &fakeFetcher{err: errors.New("connection refused")},
			repo:      &fakeRatesRepository{},
			n:         1,
			m:         3,
			wantErr:   true,
			errSubstr: "fetch",
		},
		{
			name:      "empty asks slice",
			fetcher:   &fakeFetcher{rates: &dto.StockRates{Asks: []dto.StockRatesEntry{}, Bids: validStockRates().Bids}},
			repo:      &fakeRatesRepository{},
			n:         1,
			m:         3,
			wantErr:   true,
			errSubstr: "has no asks",
		},
		{
			name:      "empty bids slice",
			fetcher:   &fakeFetcher{rates: &dto.StockRates{Asks: validStockRates().Asks, Bids: []dto.StockRatesEntry{}}},
			repo:      &fakeRatesRepository{},
			n:         1,
			m:         3,
			wantErr:   true,
			errSubstr: "has no bids",
		},
		{
			name:            "repository save returns error",
			fetcher:         &fakeFetcher{rates: validStockRates()},
			repo:            &fakeRatesRepository{err: errors.New("db write failed")},
			n:               1,
			m:               3,
			wantErr:         true,
			errSubstr:       "save",
			checkSavedCount: 0,
		},
		{
			name:            "n > len(asks) causes calculator error",
			fetcher:         &fakeFetcher{rates: validStockRates()},
			repo:            &fakeRatesRepository{},
			n:               5,
			m:               10,
			wantErr:         true,
			errSubstr:       "topN",
			checkSavedCount: 0,
		},
		{
			name:            "m < n causes calculator error",
			fetcher:         &fakeFetcher{rates: validStockRates()},
			repo:            &fakeRatesRepository{},
			n:               3,
			m:               1,
			wantErr:         true,
			errSubstr:       "avgNM",
			checkSavedCount: 0,
		},
		{
			name:            "snapshot fields are correctly set",
			fetcher:         &fakeFetcher{rates: validStockRates()},
			repo:            &fakeRatesRepository{},
			n:               1,
			m:               3,
			wantErr:         false,
			wantAskPrice:    decimal.NewFromFloat(97.50),
			wantBidPrice:    decimal.NewFromFloat(97.00),
			checkSavedCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewRatesService(tt.fetcher, tt.repo, avgNMPrecision)

			snap, err := svc.GetRates(context.Background(), tt.n, tt.m)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errSubstr != "" {
					assert.Contains(t, err.Error(), tt.errSubstr)
				}
				return
			}

			require.NoError(t, err)

			assert.True(t, snap.AskPrice.Equal(tt.wantAskPrice), "AskPrice should be %v, got %v", tt.wantAskPrice, snap.AskPrice)
			assert.True(t, snap.BidPrice.Equal(tt.wantBidPrice), "BidPrice should be %v, got %v", tt.wantBidPrice, snap.BidPrice)
			assert.False(t, snap.TopN.IsZero(), "TopN should not be zero")
			assert.False(t, snap.AvgNM.IsZero(), "AvgNM should not be zero")
			assert.False(t, snap.FetchedAt.IsZero(), "FetchedAt should be set")
			assert.Equal(t, time.Unix(validStockRates().Timestamp, 0).UTC(), snap.FetchedAt)

			require.Len(t, tt.repo.saved, tt.checkSavedCount, "repository saved count mismatch")
			if tt.checkSavedCount > 0 {
				savedSnap := tt.repo.saved[0]
				assert.True(t, savedSnap.AskPrice.Equal(snap.AskPrice))
				assert.True(t, savedSnap.BidPrice.Equal(snap.BidPrice))
				assert.True(t, savedSnap.TopN.Equal(snap.TopN))
				assert.True(t, savedSnap.AvgNM.Equal(snap.AvgNM))
				assert.Equal(t, snap.FetchedAt, savedSnap.FetchedAt)
			}
		})
	}
}

func TestNewRatesService_MetricsInit(t *testing.T) {
	t.Parallel()
	fetcher := &fakeFetcher{}
	repo := &fakeRatesRepository{}
	svc := NewRatesService(fetcher, repo, 2)
	assert.NotNil(t, svc.fetchTotal)
	assert.NotNil(t, svc.fetchErrors)
	assert.NotNil(t, svc.fetchDuration)
}
