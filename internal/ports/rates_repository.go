package ports

import (
	"context"

	"github.com/Danil-Ivonin/GrintexRates/internal/domain"
)

type RatesRepository interface {
	Create(ctx context.Context, rates domain.RateSnapshot) error
}
