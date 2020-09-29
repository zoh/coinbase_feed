package models

import feeder "coinbase_feed"

type TickerChannels struct {
	symbols []string
	chans   map[string]chan feeder.Ticker
}

func NewTickerChannel(symbols []string) *TickerChannels {
	var chans = make(map[string]chan feeder.Ticker)
	for _, c := range symbols {
		chans[c] = make(chan feeder.Ticker)
	}

	return &TickerChannels{
		symbols: symbols,
		chans:   chans,
	}
}

func (t *TickerChannels) CloseChannels() {
	for _, c := range t.chans {
		close(c)
	}
}

func (t *TickerChannels) GetSymbols() []string {
	return t.symbols
}

func (t *TickerChannels) GetChannel(symbol string) (chan feeder.Ticker, bool) {
	c, ok := t.chans[symbol]
	return c, ok
}

func (t *TickerChannels) Iter(fn func(symbol string, ch chan feeder.Ticker)) {
	for s, c := range t.chans {
		fn(s, c)
	}
}
