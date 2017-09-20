package bittrex

import (
	"encoding/json"
)

type ClientMethodCall struct {
	H string          //hub name
	M string          //method
	A json.RawMessage //method parameters
}

type MarketUpdate struct {
	MarketName string          `json:"MarketName"`
	Nonce      int             `json:"Nounce"`
	BidUpdates []BookRowUpdate `json:"Buys"`
	AskUpdates []BookRowUpdate `json:"Sells"`
	Fills      []Fill          `json:"Fills"`
}

func (c *StreamClient) clientMethodHandler(calls []ClientMethodCall) {
	for _, call := range calls {
		switch call.M {
		case "updateExchangeState":
			var updates []MarketUpdate
			err := json.Unmarshal(call.A, &updates)
			if err != nil {
				c.ErrorHandler(err)
			}
			for _, u := range updates {
				market, ok := c.Markets[u.MarketName]
				if ok {
					go market.updateHandler(u, false)
				}
			}
		default:
			//unknown client method
			continue
		}
	}
}
