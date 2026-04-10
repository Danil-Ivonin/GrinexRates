package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Danil-Ivonin/GrinexRates/internal/http/dto"
	"github.com/go-resty/resty/v2"
)

type GrinexClient struct {
	httpClient *resty.Client
	url        string
}

func New(url string, timeout time.Duration) *GrinexClient {
	rc := resty.New().SetTimeout(timeout)
	return &GrinexClient{httpClient: rc, url: url}
}

func (c *GrinexClient) Fetch(ctx context.Context) (*dto.StockRates, error) {
	resp, err := c.httpClient.R().SetContext(ctx).Get(c.url)
	if err != nil {
		return nil, fmt.Errorf("client: fetch: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("client: fetch: non-2xx status %d", resp.StatusCode())
	}

	var ob dto.StockRates
	if err := json.Unmarshal(resp.Body(), &ob); err != nil {
		return nil, fmt.Errorf("client: fetch: decode response: %w", err)
	}
	return &ob, nil
}
