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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Fetcher interface {
	Fetch(ctx context.Context) (*dto.StockRates, error)
}

type RatesService struct {
	fetcher        Fetcher
	repo           ports.RatesRepository
	avgNMPrecision int32
	fetchTotal     metric.Int64Counter
	fetchErrors    metric.Int64Counter
	fetchDuration  metric.Float64Histogram
}

func NewRatesService(fetcher Fetcher, repo ports.RatesRepository, avgNMPrecision int32) *RatesService {
	meter := otel.Meter("grinex-rates/service")

	fetchTotal, _ := meter.Int64Counter(
		"grinex.fetch.total",
		metric.WithDescription("Total number of Grinex API fetch attempts"),
		metric.WithUnit("{fetches}"),
	)
	fetchErrors, _ := meter.Int64Counter(
		"grinex.fetch.errors.total",
		metric.WithDescription("Total number of failed Grinex API fetch attempts"),
		metric.WithUnit("{fetches}"),
	)
	fetchDuration, _ := meter.Float64Histogram(
		"grinex.fetch.duration",
		metric.WithDescription("Duration of Grinex API fetch calls"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
	)

	return &RatesService{
		fetcher:        fetcher,
		repo:           repo,
		avgNMPrecision: avgNMPrecision,
		fetchTotal:     fetchTotal,
		fetchErrors:    fetchErrors,
		fetchDuration:  fetchDuration,
	}
}

// GetRates fetches rates, computes topN and avgNM for the given
// positions, persists the result, and returns the populated RateSnapshot
func (s *RatesService) GetRates(ctx context.Context, n, m int32) (domain.RateSnapshot, error) {
	s.fetchTotal.Add(ctx, 1)

	start := time.Now()
	rates, err := s.fetcher.Fetch(ctx)
	elapsed := time.Since(start).Seconds()
	s.fetchDuration.Record(ctx, elapsed)

	if err != nil {
		s.fetchErrors.Add(ctx, 1)
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

	avgNM, err := calculator.AvgNM(askPrices, int(n), int(m), s.avgNMPrecision)
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
