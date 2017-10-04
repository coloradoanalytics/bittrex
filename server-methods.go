package bittrex

import (
	"encoding/json"
	"fmt"
	"sync"
)

type ServerMethodResult struct {
	I string          //invocation Id (always present)
	R json.RawMessage //the value returned by the server method (present if the method is not void)
	E string          //error message
	H bool            //true if this is a hub error
	D json.RawMessage //an object containing additional error data (can only be present for hub errors)
	T json.RawMessage //stack trace
	S json.RawMessage //state - a dictionary containing additional custom data
}

func (c *StreamClient) methodResultHandler(msg []byte) {
	var result ServerMethodResult
	err := json.Unmarshal(msg, &result)
	if err != nil {
		c.OnError(err)
	}

	call, ok := c.MethodCalls.Get(result.I)
	if ok {
		go call.Caller.handleMethodResult(call.Method, result)
	}

}

type MethodCaller interface {
	handleMethodResult(string, ServerMethodResult)
}

type MethodCallMap struct {
	sync.RWMutex
	Callers map[string]MethodCall
}

type MethodCall struct {
	Method string
	Caller MethodCaller
}

func (m *MethodCallMap) Get(key string) (MethodCall, bool) {
	//block writes
	m.Lock()
	defer m.Unlock()
	value, ok := m.Callers[key]
	return value, ok
}

func (m *MethodCallMap) Set(key string, value MethodCall) {
	//block reads
	m.RLock()
	defer m.RUnlock()
	m.Callers[key] = value
}

func (m *MethodCallMap) Remove(key string) {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.Callers[key]
	if ok {
		delete(m.Callers, key)
	}
}

//////////////// MARKET /////////////////////////

func (m *Market) handleMethodResult(method string, result ServerMethodResult) {
	switch method {
	case "queryExchangeState":
		var update MarketUpdate
		err := json.Unmarshal(result.R, &update)
		if err != nil {
			m.OnError(err)
		}
		fmt.Println(fmt.Sprintf("%+v", update.Nonce))
		m.updateHandler(update, true)
		m.Ready = true
	default:
		//unknown method
		return
	}
}
