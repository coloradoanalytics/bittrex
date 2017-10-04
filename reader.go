package bittrex

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

func (c *StreamClient) reader(startChan chan error) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("error: %v", err)
			}
			break
		}

		//fmt.Println(string(msg))

		var message struct {
			C string `json:"C,omitempty"` // present for persistent connection messages
			I string `json:"I,omitempty"` // present for method result messages
		}

		err = json.Unmarshal(msg, &message)
		if err != nil {
			c.OnError(err)
		}

		if message.C != "" {
			var pcm PersistentConnectionMessage
			err = json.Unmarshal(msg, &pcm)
			if err != nil {
				c.OnError(err)
			}

			if pcm.S == 1 {
				startChan <- c.sendStart()
			} else {
				//assume any message that makes it this far is an array of client method calls
				var calls []ClientMethodCall
				err = json.Unmarshal(pcm.M, &calls)
				if err != nil {
					c.OnError(err)
				}

				c.clientMethodHandler(calls)
			}
		} else if message.I != "" {
			c.methodResultHandler(msg)
		}
	}
}

type PersistentConnectionMessage struct {
	C string          //message id, present for all non-KeepAlive messages
	M json.RawMessage //an array containing actual data
	S int             //indicates that the transport was initialized (a.k.a. init message)
	G string          //groups token – an encrypted string representing group membership
}
