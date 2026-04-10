package storage

import (
	"context"
	"fmt"

	"github.com/Danil-Ivonin/GrintexRates/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatesRepository struct {
	pool *pgxpool.Pool
}

func NewRatesRepository(pool *pgxpool.Pool) *RatesRepository {
	return &RatesRepository{pool: pool}
}

func (r *RatesRepository) Create(ctx context.Context, rates domain.RateSnapshot) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO rate_snapshots (ask_price, bid_price, top_n, avg_nm, fetched_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		rates.AskPrice,
		rates.BidPrice,
		rates.TopN,
		rates.AvgNM,
		rates.FetchedAt,
	)
	if err != nil {
		return fmt.Errorf("storage: save snapshot: %w", err)
	}
	return nil
}
