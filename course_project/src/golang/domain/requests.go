package domain

type RequestWS struct {
	Event      string   `json:"event,omitempty"`
	Feed       string   `json:"feed,omitempty"`
	ProductIDs []string `json:"product_ids,omitempty"`
}

func NewRequestWS(event string, feed string, productIDs []string) *RequestWS {
	return &RequestWS{Event: event, Feed: feed, ProductIDs: productIDs}
}

func SubscribeTickerRequestWS(productIDs []string) *RequestWS {
	return &RequestWS{
		Event:      "subscribe",
		Feed:       "ticker",
		ProductIDs: productIDs,
	}
}

func UnsubscribeTickerRequestWS(productIDs []string) *RequestWS {
	return &RequestWS{
		Event:      "unsubscribe",
		Feed:       "ticker",
		ProductIDs: productIDs,
	}
}
