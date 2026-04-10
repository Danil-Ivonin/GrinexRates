package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Danil-Ivonin/GrinexRates/internal/calculator"
	"github.com/Danil-Ivonin/GrinexRates/internal/domain"
	"github.com/Danil-Ivonin/GrinexRates/internal/http/dto"
	"github.com/Danil-Ivonin/GrinexRates/internal/ports"
	"github.com/shopspring/decimal"
)

type Fetcher interface {
	Fetch(ctx context.Context) (*dto.StockRates, error)
}

type RatesService struct {
	fetcher Fetcher
	repo    ports.RatesRepository
}

func NewRatesService(fetcher Fetcher, repo ports.RatesRepository) *RatesService {
	return &RatesService{fetcher: fetcher, repo: repo}
}

// GetRates fetches rates, computes topN and avgNM for the given
// positions, persists the result, and returns the populated RateSnapshot
func (s *RatesService) GetRates(ctx context.Context, n, m int32) (domain.RateSnapshot, error) {
	rates, err := s.fetcher.Fetch(ctx)
	if err != nil {
		return domain.RateSnapshot{}, fmt.Errorf("service: get rates: fetch: %w", err)
	}

	if len(rates.Asks) == 0 {
		return domain.RateSnapshot{}, fmt.Errorf("service: get rates: rates has no asks")
	}

	if len(rates.Bids) == 0 {
		return domain.RateSnapshot{}, fmt.Errorf("service: get rates: rates has no bids")
	}

	// Extract price from rate entries
	askPrices := make([]decimal.Decimal, len(rates.Asks))
	for i, entry := range rates.Asks {
		askPrices[i] = entry.Price
	}

	topN, err := calculator.TopN(askPrices, int(n))
	if err != nil {
		return domain.RateSnapshot{}, fmt.Errorf("service: get rates: topN: %w", err)
	}

	avgNM, err := calculator.AvgNM(askPrices, int(n), int(m))
	if err != nil {
		return domain.RateSnapshot{}, fmt.Errorf("service: get rates: avgNM: %w", err)
	}

	// Best ask = asks[0].Price, best bid = bids[0].Price
	askPrice, bidPrice := rates.Asks[0].Price, rates.Bids[0].Price

	snap := domain.RateSnapshot{
		AskPrice:  askPrice,
		BidPrice:  bidPrice,
		TopN:      topN,
		AvgNM:     avgNM,
		FetchedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, snap); err != nil {
		return domain.RateSnapshot{}, fmt.Errorf("service: get rates: save: %w", err)
	}

	return snap, nil
}
