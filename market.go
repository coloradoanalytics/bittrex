package bittrex

import (
	"sync"
)

type Market struct {
	sync.Mutex
	MarketName   string
	MarketAsset  string
	BaseAsset    string
	Bids         BookRowsDescending
	Asks         BookRowsAscending
	Ready        bool
	LastNonce    int
	StreamClient *StreamClient
	OnError      func(error)
	OnUpdate     func(*Market, MarketUpdate, bool)
}

func NewMarket(baseAsset string, marketAsset string, client *StreamClient) *Market {
	market := Market{
		BaseAsset:    baseAsset,
		MarketAsset:  marketAsset,
		MarketName:   baseAsset + "-" + marketAsset,
		StreamClient: client,
		OnError:      func(err error) { panic(err) },              //default error handler
		OnUpdate:     func(m *Market, mu MarketUpdate, s bool) {}, //default after-update hook
	}

	client.Markets[market.MarketName] = &market

	return &market
}

func (m *Market) Subscribe() error {
	err := m.subscribeToDeltas()
	if err != nil {
		return err
	}

	err = m.queryBook()
	return err
}

func (m *Market) GetBid(depth int) BookRow {
	m.Lock()
	defer m.Unlock()
	row := BookRow{}
	if len(m.Bids) > depth {
		row.Price = m.Bids[depth].Price
		row.Quantity = m.Bids[depth].Quantity
	}
	return row
}

func (m *Market) GetAsk(depth int) BookRow {
	m.Lock()
	defer m.Unlock()
	row := BookRow{}
	if len(m.Asks) > depth {
		row.Price = m.Asks[depth].Price
		row.Quantity = m.Asks[depth].Quantity
	}
	return row
}
