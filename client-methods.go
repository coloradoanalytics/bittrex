package bittrex

import (
	"encoding/json"
)

type ClientMethodCall struct {
	H string            //hub name
	M string            //method
	A []json.RawMessage //method parameters
}

type UpdateMarketStateParams struct {
	MarketName string
	Nonce      int             `json:"Nounce"`
	BidUpdates []BookRowUpdate `json:"Buys"`
	AskUpdates []BookRowUpdate `json:"Sells"`
	Fills      []Fill          `json:"Fills"`
}

func (c *StreamClient) clientMethodHandler(calls []ClientMethodCall) {
	for _, call := range calls {
		switch call.M {
		case "updateExchangeState":
			var params []UpdateMarketStateParams
			err := json.Unmarshal(call.A, &params)
			if err != nil {
				c.ErrorHandler(err)
			}
			for _, p := range params {
				market, ok := c.Markets[p.MarketName]
				if ok {
					go c.Markets[p.MarketName].updateHandler(p)
				}
			}
		default:
			//unknown client method
			continue
		}
	}
}
