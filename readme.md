# Bittrex SignalR/Websockets client

Roughly written by an amateur. Many improvements are possible and help is appreciated.

For now this library just connects to the Bittrex SignalR/Websockets service, subscribes to markets, and
keeps an up-to-date order book. After Bittrex's v2 API is available, I hope to add that functionality.
For now, it's possible to combine this library with `github.com/toorop/go-bittrex` if you want to use the
REST API

```
package main

import (
  "github.com/coloradoanalytics/bittrex"
  "fmt"
  "time"
)

func main() {
  client := bittrex.NewStreamClient()
  
  market := bittrex.NewMarket("BTC", "ETH", client)
  market.onUpdate = makeUpdateHandler()
  
  client.Connect()
  market.Subscribe()
  
  //run for 1 minute and then disconnect from Bittrex
  time.Sleep(1 * time.Minute)
  
  
}

func makeUpdateHandler() func(*bittrex.Market, bittrex.MarketUpdate, bool) {

  return func(m *bittrex.Market, update bittrex.MarketUpdate, setState bool) {
    //do something in response to an order book update
    
    //update is the message that came from Bittrex and was used to update the book
    //update contains the fills if you want to do something with them
    //setState is true if the book was set all at once vs false if it was updated with a few changes
    //m is a pointer to the market that was updated
    
    fmt.Println("Bid in", m.MarketName, "is now", m.GetBid(0).Price)
  }
}
```
Bittrex doesn't provide a way to stop

To retrieve any level of the order book, use market.GetBid(depth) or market.GetAsk(depth)
