package models

import (
	feeder "coinbase_feed"
	"strconv"
)

// Ticker for storage
type Ticker struct {
	Timestamp int64
	Symbol    string
	Bid       float64
	Ask       float64
}

func ConvertTickerFromStreamTicker(t feeder.Ticker) (Ticker, error) {
	bid, err := strconv.ParseFloat(t.BestBid, 8)
	if err != nil {
		return Ticker{}, err
	}

	ask, err := strconv.ParseFloat(t.BestAsk, 8)
	if err != nil {
		return Ticker{}, err
	}

	return Ticker{
		Timestamp: t.GetReceivedTime(),
		Symbol:    t.ProductID,
		Bid:       bid,
		Ask:       ask,
	}, nil
}
