package storage

import (
	"coinbase_feed/models"
	"context"
)

// Storager ...
type Storager interface {
	Write(context.Context, models.Ticker) error
	Close() error

	ReadLast(ctx context.Context, symbol string) (*models.Ticker, error)
}
