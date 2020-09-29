// Package feeder ...
//
// Websocket Feed
// The websocket feed provides real-time market data updates for orders and trades.
package feeder

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
)

var (
	coinbaseWsEndpoint = "wss://ws-feed.pro.coinbase.com"
	websocketTimeout   = time.Second * 60
	keepalive          = false
)

func NewCoinbaseWsOptions() *WsOptions {
	return &WsOptions{
		Endpoint:         coinbaseWsEndpoint,
		Keealive:         true,
		WebsocketTimeout: websocketTimeout,
		Debug:            false,
		Logger:           log.New(os.Stdout, "Coinbase feed: ", log.LstdFlags),
	}
}

type TickerHandler func(Ticker)

// StartFeedTicker The ticker channel provides real-time price updates every time a match happens.
// It batches updates in case of cascading matches, greatly reducing bandwidth requirements.
func StartFeedTicker(ctx context.Context, cfg *WsOptions, symbols []string, respHandler TickerHandler, errHandler ErrHandler) (doneC chan struct{}, err error) {
	if len(symbols) == 0 {
		return nil, EmptyRequest
	}

	var handler = func(data []byte) {
		if t, err := parseTicker(data); err != nil {
			errHandler(err)
		} else if t.receivedTime > 0 {
			// not send empty struct
			respHandler(t)
		}
	}

	return startServer(ctx, cfg, SubscribeTickers(symbols...), handler, errHandler)
}

// Ticker stream event
type Ticker struct {
	Base
	Sequence int64 `json:"sequence"`
	// Symbol name
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
	Open24H   string `json:"open_24h"`
	Volume24H string `json:"volume_24h"`
	Low24H    string `json:"low_24h"`
	High24H   string `json:"high_24h"`
	Volume30D string `json:"volume_30d"`
	BestBid   string `json:"best_bid"`
	BestAsk   string `json:"best_ask"`
	Side      string `json:"side"`
	Time      string `json:"time"`
	TradeID   int64  `json:"trade_id"`
	LastSize  string `json:"last_size"`
}

func parseTicker(data []byte) (Ticker, error) {
	err := isError(data)
	if err != nil {
		return Ticker{}, err
	}
	var t Ticker
	err = json.Unmarshal(data, &t)
	if err != nil {
		return Ticker{}, err
	}

	// ignore if not "ticker"
	if t.Type != string(tickerTypeChannel) {
		return Ticker{}, nil
	}
	t.receivedTime = makeTimestamp()

	return t, nil
}
