package feeder

type (
	typeChannels string
)

const (
	subscribeActionType     = "subscribe"
	subscriptionsActionType = "subscriptions"
	errorActionType         = "error"

	tickerTypeChannel typeChannels = "ticker"
)

type SubscribeReq struct {
	Type       string           `json:"type"`
	ProductIDS []string         `json:"product_ids"`
	Channels   []ChannelElement `json:"channels"`
}

type ChannelClass struct {
	Name       string   `json:"name"`
	ProductIDS []string `json:"product_ids"`
}

type ChannelElement struct {
	ChannelClass *ChannelClass
	String       *string
}

func (x *ChannelElement) UnmarshalJSON(data []byte) error {
	x.ChannelClass = nil
	var c ChannelClass
	object, err := unmarshalUnion(data, nil, nil, nil, &x.String, false, nil, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.ChannelClass = &c
	}
	return nil
}

func (x *ChannelElement) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, false, nil, x.ChannelClass != nil, x.ChannelClass, false, nil, false, nil, false)
}

// SubscribeTickers build body for subscribe on channel ticker
func SubscribeTickers(symbols ...string) SubscribeReq {
	s := string(tickerTypeChannel)
	return SubscribeReq{
		Type:       subscribeActionType,
		ProductIDS: symbols,
		Channels: []ChannelElement{
			{
				String: &s,
			},
		},
	}
}
