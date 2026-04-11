package calculator_test

import (
	"testing"

	"github.com/Danil-Ivonin/GrinexRates/internal/calculator"
	"github.com/shopspring/decimal"
)

func TestTopN(t *testing.T) {
	t.Parallel()
	asksExmpl := []decimal.Decimal{
		decimal.NewFromFloat(100.5),
		decimal.NewFromFloat(101.0),
		decimal.NewFromFloat(102.0),
	}
	tests := []struct {
		name    string
		asks    []decimal.Decimal
		n       int
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name:    "first element",
			asks:    asksExmpl,
			n:       1,
			want:    decimal.NewFromFloat(100.5),
			wantErr: false,
		},
		{
			name:    "last element",
			asks:    asksExmpl,
			n:       3,
			want:    decimal.NewFromFloat(102.0),
			wantErr: false,
		},
		{
			name:    "n less than 1",
			asks:    asksExmpl,
			n:       0,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "empty slice",
			asks:    []decimal.Decimal{},
			n:       1,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "n exceeds length",
			asks:    asksExmpl,
			n:       4,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "negative n",
			asks:    asksExmpl,
			n:       -1,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "single element slice",
			asks:    []decimal.Decimal{decimal.NewFromFloat(100.5)},
			n:       1,
			want:    decimal.NewFromFloat(100.5),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := calculator.TopN(tt.asks, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("TopN(%v, %d) error = %v, wantErr %v", tt.asks, tt.n, err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("TopN(%v, %d) = %v, want %v", tt.asks, tt.n, got, tt.want)
			}
		})
	}
}

func TestAvgNM(t *testing.T) {
	t.Parallel()
	asksExmpl := []decimal.Decimal{
		decimal.NewFromInt(100),
		decimal.NewFromFloat(101),
		decimal.NewFromFloat(102),
	}
	tests := []struct {
		name    string
		asks    []decimal.Decimal
		n       int
		m       int
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name:    "mean of three elements",
			asks:    asksExmpl,
			n:       1,
			m:       3,
			want:    decimal.NewFromInt(101),
			wantErr: false,
		},
		{
			name:    "single element mean",
			asks:    []decimal.Decimal{decimal.NewFromInt(101)},
			n:       1,
			m:       1,
			want:    decimal.NewFromInt(101),
			wantErr: false,
		},
		{
			name:    "mean of subset",
			asks:    asksExmpl,
			n:       1,
			m:       2,
			want:    decimal.NewFromFloat(100.5),
			wantErr: false,
		},
		{
			name:    "n less than 1",
			asks:    asksExmpl,
			n:       0,
			m:       3,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "m less than n",
			asks:    asksExmpl,
			n:       3,
			m:       1,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "m exceeds length",
			asks:    asksExmpl,
			n:       1,
			m:       4,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "empty slice",
			asks:    []decimal.Decimal{},
			n:       1,
			m:       1,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "negative n",
			asks:    asksExmpl,
			n:       -1,
			m:       2,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:    "n equals m middle of slice",
			asks:    asksExmpl,
			n:       2,
			m:       2,
			want:    decimal.NewFromInt(101),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := calculator.AvgNM(tt.asks, tt.n, tt.m, 5)
			if (err != nil) != tt.wantErr {
				t.Errorf("AvgNM(%v, %d, %d) error = %v, wantErr %v", tt.asks, tt.n, tt.m, err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("AvgNM(%v, %d, %d) = %v, want %v", tt.asks, tt.n, tt.m, got, tt.want)
			}
		})
	}
}
