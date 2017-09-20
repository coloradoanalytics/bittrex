package bittrex

import (
	"encoding/json"
	//"fmt"
	"sort"
	"strconv"
	"time"
)

func (m *Market) subscribeToDeltas() error {
	i := strconv.Itoa(time.Now().Nanosecond())

	message := CallMethodMessage{
		M: "subscribeToExchangeDeltas",
		A: []string{m.MarketName},
		I: i,
		H: "coreHub",
	}

	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	m.StreamClient.WriteChan <- msg
	return nil
}

func (m *Market) queryBook() error {
	i := strconv.Itoa(time.Now().Nanosecond())

	message := CallMethodMessage{
		M: "queryExchangeState",
		A: []string{m.MarketName},
		I: i,
		H: "coreHub",
	}

	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	call := MethodCall{
		Method: "queryExchangeState",
		Caller: m,
	}

	m.StreamClient.MethodCalls.Set(i, call)
	m.StreamClient.WriteChan <- msg
	return nil
}

func (m *Market) updateHandler(update MarketUpdate, setState bool) {
	//fmt.Println("updateHandler")
	if m.Ready != true && !setState {
		//currently out of sync, so record nonce and ignore update
		m.LastNonce = update.Nonce
		return
	}

	if m.LastNonce == update.Nonce && !setState {
		//repeated record, ignore
		return
	}

	if m.LastNonce+1 != update.Nonce && !setState {
		//market is out of sync
		m.Ready = false
		m.Bids = nil
		m.Asks = nil
		err := m.queryBook()
		if err != nil {
			m.OnError(err)
		}
		return
	}

	if len(update.BidUpdates) > 0 {
		bidAdd, bidRemove, bidUpdate := separateRowUpdates(update.BidUpdates)
		if len(bidAdd) > 0 {
			m.addBids(bidAdd)
		}

		if len(bidRemove) > 0 {
			m.removeBids(bidRemove)
		}

		if len(bidUpdate) > 0 {
			m.updateBids(bidUpdate)
		}
	}

	if len(update.AskUpdates) > 0 {
		askAdd, askRemove, askUpdate := separateRowUpdates(update.AskUpdates)
		if len(askAdd) > 0 {
			m.addAsks(askAdd)
		}

		if len(askRemove) > 0 {
			m.removeAsks(askRemove)
		}

		if len(askUpdate) > 0 {
			m.updateAsks(askUpdate)
		}
	}

	if len(update.Fills) > 0 {
		//handle fills hook
	}

	m.LastNonce = update.Nonce

	m.OnUpdate(m)

}

func (m *Market) addBids(updates []BookRowUpdate) {
	m.Lock()
	defer m.Unlock()
	for _, u := range updates {
		m.Bids = append(m.Bids, BookRow{Price: u.Price, Quantity: u.Quantity})
	}
	sort.Sort(m.Bids)
}

func (m *Market) removeBids(updates []BookRowUpdate) {
	m.Lock()
	defer m.Unlock()
	for _, u := range updates {
		for i, b := range m.Bids {
			if b.Price == u.Price {
				m.Bids = append(m.Bids[:i], m.Bids[i+1:]...)
				break
			}
		}
	}
}

func (m *Market) updateBids(updates []BookRowUpdate) {
	m.Lock()
	defer m.Unlock()
	for _, u := range updates {
		for i, _ := range m.Bids {
			if m.Bids[i].Price == u.Price {
				m.Bids[i].Quantity = u.Quantity
				break
			}
		}
	}
}

func (m *Market) addAsks(updates []BookRowUpdate) {
	m.Lock()
	defer m.Unlock()
	for _, u := range updates {
		m.Asks = append(m.Asks, BookRow{Price: u.Price, Quantity: u.Quantity})
	}
	sort.Sort(m.Asks)
}

func (m *Market) removeAsks(updates []BookRowUpdate) {
	m.Lock()
	defer m.Unlock()
	for _, u := range updates {
		for i, b := range m.Asks {
			if b.Price == u.Price {
				m.Asks = append(m.Asks[:i], m.Asks[i+1:]...)
				break
			}
		}
	}
}

func (m *Market) updateAsks(updates []BookRowUpdate) {
	m.Lock()
	defer m.Unlock()
	for _, u := range updates {
		for i, _ := range m.Asks {
			if m.Asks[i].Price == u.Price {
				m.Asks[i].Quantity = u.Quantity
				break
			}
		}
	}
}
