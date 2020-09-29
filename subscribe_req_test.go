package feeder_test

import (
	feeder "coinbase_feed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SubscribeTickers(t *testing.T) {
	data := feeder.SubscribeTickers("ETH-USD", "ETH-EUR")

	raw, err := json.Marshal(data)
	assert.NoError(t, err)

	assert.Equal(t, string(raw), `{"type":"subscribe","product_ids":["ETH-USD","ETH-EUR"],"channels":["ticker"]}`)
}
