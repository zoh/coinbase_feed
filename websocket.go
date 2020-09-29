package feeder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
)

type WsHandler func(message []byte)

type ErrHandler func(err error)

type Base struct {
	Type         string `json:"type"`
	receivedTime int64  `json:"-"`
}

func (b Base) GetReceivedTime() int64 {
	return b.receivedTime
}

// ErrorResponse error response message from websocket
type ErrorResponse struct {
	Base
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Reason)
}

func isError(data []byte) error {
	v, err := fastjson.ParseBytes(data)
	if err != nil {
		return err
	}

	if string(v.Get("type").GetStringBytes()) == errorActionType {
		return ErrorResponse{
			Base: Base{
				Type:         errorActionType,
				receivedTime: makeTimestamp(),
			},
			Message: string(v.Get("message").GetStringBytes()),
			Reason:  string(v.Get("reason").GetStringBytes()),
		}
	}
	return nil
}

type WsOptions struct {
	Endpoint         string
	Keealive         bool
	WebsocketTimeout time.Duration
	Debug            bool
	// Use with debug flag
	Logger *log.Logger
}

func (opt *WsOptions) debug(v ...interface{}) {
	if opt.Debug {
		opt.Logger.Println(v...)
	}
}

func startServer(ctx context.Context, opt *WsOptions, body interface{}, handler WsHandler, errHandler ErrHandler) (doneC chan struct{}, err error) {
	c, _, err := websocket.DefaultDialer.DialContext(ctx, opt.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	doneC = make(chan struct{})

	if body != nil {
		req, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		opt.debug("send to ws:", string(req))
		err = c.WriteMessage(websocket.TextMessage, req)
		if err != nil {
			return nil, err
		}
	}

	var closeServer = func() {
		opt.debug("close ws connection")
		c.Close()
	}

	go func() {
		select {
		case <-ctx.Done():
			closeServer()
		case <-doneC:
			closeServer()
		}
	}()

	go func() {
		defer close(doneC)
		if opt.Keealive {
			keepAlive(c, opt.WebsocketTimeout)
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					errHandler(err)
					return
				}
				opt.debug("input message", string(message))
				handler(message)
			}
		}
	}()
	return
}

// ping ws connection
func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Now().Sub(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
