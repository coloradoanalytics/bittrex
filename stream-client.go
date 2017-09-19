package bittrex

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

type StreamClient struct {
	ConnectionToken         string
	ConnectionId            string
	KeepAliveTimeout        float64
	DisconnectTimeout       float64
	ProtocolVersion         string
	TransportConnectTimeout float64
	Conn                    *websocket.Conn
	WriteChan               chan []byte
	Markets                 map[string]*Market
	ErrorHandler            func(error)
}

func NewStreamClient() *StreamClient {
	client := StreamClient{
		ProtocolVersion: "1.4",
		ErrorHandler:    func(err error) { panic(err) },
	}

	return &client
}

func (c *StreamClient) Connect() error {
	err := c.negotiate()
	if err != nil {
		return err
	}

	err = c.dial()
	if err != nil {
		return err
	}

	startChan := make(chan error)

	go c.reader(startChan)
	go c.writer()

	timer := time.NewTicker(10 * time.Second)

	select {
	case <-timer.C:
		//something went wrong with the start sequence
		return errors.New("timed out during start sequence")
	case err = <-startChan:
		timer.Stop()
		close(startChan)
		return err
	}

}

func (c *StreamClient) negotiate() error {
	u := url.URL{Scheme: "https", Host: "www.bittrex.com", Path: "/signalr/negotiate"}

	v := url.Values{}
	v.Set("clientProtocol", c.ProtocolVersion)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(c)
	return err
}

func (c *StreamClient) dial() error {
	u := url.URL{Scheme: "wss", Host: c.Host, Path: "/signalr/connect"}

	v := url.Values{}
	v.Set("transport", "webSockets")
	v.Set("clientProtocol", c.ProtocolVersion)
	v.Set("connectionToken", c.ConnectionToken)
	u.RawQuery = v.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.Conn = conn
	return nil
}

func (c *StreamClient) sendStart() error {
	u := url.URL{Scheme: "https", Host: "www.bittrex.com", Path: "/signalr/start"}

	v := url.Values{}
	v.Set("transport", "webSockets")
	v.Set("clientProtocol", c.ProtocolVersion)
	v.Set("connectionToken", c.ConnectionToken)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var startResponse struct {
		Response string
	}

	err = json.NewDecoder(resp.Body).Decode(&startResponse)
	if err != nil {
		return err
	}

	if startResponse.Response == "started" {
		return nil
	}

	return errors.New("Unexpected response to start request")
}

func (c *StreamClient) Close() error {
	u := url.URL{Scheme: "https", Host: "www.bittrex.com", Path: "/signalr/abort"}

	v := url.Values{}
	v.Set("transport", "webSockets")
	v.Set("clientProtocol", c.ProtocolVersion)
	v.Set("connectionToken", c.ConnectionToken)

	u.RawQuery = v.Encode()
	_, err := http.Get(u.String())

	return err
}

func (c *StreamClient) writer() {
	for {
		msg := <-c.WriteChan
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
