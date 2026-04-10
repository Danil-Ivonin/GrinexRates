package ports

import (
	"context"

	"github.com/Danil-Ivonin/GrinexRates/internal/domain"
)

type RatesService interface {
	GetRates(ctx context.Context, n, m int32) (domain.RateSnapshot, error)
}
