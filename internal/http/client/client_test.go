package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Danil-Ivonin/GrinexRates/internal/http/client"
	"github.com/shopspring/decimal"
)

const validJSON = `{
	"timestamp": 1775910345,
	"asks":[
		{"price":"79.8","volume":"58219.6293","amount":"4645926.42"},
		{"price":"79.83","volume":"3135.507","amount":"250307.52"}
	],
	"bids":[
		{"price":"79.72","volume":"5550.1288","amount":"442456.27"}
	]
}`

func TestFetch_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(validJSON))
	}))
	defer srv.Close()

	c := client.New(srv.URL, time.Second)

	ob, err := c.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// timestamp
	if ob.Timestamp != 1775910345 {
		t.Errorf("expected timestamp 1775910345, got %d", ob.Timestamp)
	}

	// asks
	if len(ob.Asks) == 0 {
		t.Fatal("asks should not be empty")
	}

	expectedAsk := decimal.RequireFromString("79.8")
	if !ob.Asks[0].Price.Equal(expectedAsk) {
		t.Errorf("expected ask price %s, got %s", expectedAsk, ob.Asks[0].Price)
	}

	// bids
	if len(ob.Bids) == 0 {
		t.Fatal("bids should not be empty")
	}

	expectedBid := decimal.RequireFromString("79.72")
	if !ob.Bids[0].Price.Equal(expectedBid) {
		t.Errorf("expected bid price %s, got %s", expectedBid, ob.Bids[0].Price)
	}
}

func TestFetch_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := client.New(srv.URL, time.Second)

	_, err := c.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}

	if !strings.Contains(err.Error(), "non-2xx") {
		t.Errorf("expected error to contain 'non-2xx', got %v", err)
	}
}

func TestFetch_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"asks":`))
	}))
	defer srv.Close()

	c := client.New(srv.URL, time.Second)

	_, err := c.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected JSON decode error")
	}
}

func TestFetch_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(validJSON))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := client.New(srv.URL, time.Second)

	_, err := c.Fetch(ctx)
	if err == nil {
		t.Fatal("expected context cancelled error")
	}
}

func TestFetch_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := client.New(srv.URL, 50*time.Millisecond)

	_, err := c.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
