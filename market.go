package bittrex

import (
	"sync"
)

type Market struct {
	sync.Mutex
	MarketName  string
	MarketAsset string
	BaseAsset   string
	Bids        BookRowsDescending
	Asks        BookRowsAscending
	Ready       bool
	LastNonce   int
	OnError     func(error)
	OnUpdate    func(*Market)
}

func NewMarket(baseAsset string, marketAsset string) *Market {
	market := Market{
		BaseAsset:   baseAsset,
		MarketAsset: marketAsset,
		MarketName:  baseAsset + "-" + marketAsset,
		OnError:     func(err error) { panic(err) }, //default error handler
		OnUpdate:    func(m *Market) {},             //default after-update hook
	}

	return &market
}

func (m *Market) fetchBook() error {
	return nil
}
