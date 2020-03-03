package ws

import (
	"context"
	"encoding/json"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"gitlab.com/lattetalk/lattetalk/ws/data"
)

type Conn struct {
	ID        int
	rwc       net.Conn
	subject   int64
	w         *wsutil.Writer
	writeChan chan interface{}
}

func (c *Conn) readRequest() (*data.Request, error) {
	h, r, err := wsutil.NextReader(c.rwc, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlHandler{
			Src:   r,
			Dst:   c.rwc,
			State: ws.StateServerSide}.Handle(h)
	}

	req := new(data.Request)
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(req); err != nil {
		return nil, err
	}
	req.SetSubject(c.subject)

	return req, nil
}

func (c *Conn) Read() error {
	req, err := c.readRequest()
	if err != nil {
		return err
	}

	// Handle control message.
	if req == nil {
		return nil
	}

	switch req.Method {
	case data.RequestNewMessage:
		return handleNewMessage(req)
	default:
		c.Output(data.ResponseError, "method is not supported")
	}

	return nil
}

func (c *Conn) Output(method data.ResponseMethod, obj interface{}) {
	go func() {
		c.writeChan <- &data.Response{
			Method: method,
			Data:   obj,
		}
	}()
}

func (c *Conn) Write(ctx context.Context) error {
	for {
		var obj interface{}
		select {
		case <-ctx.Done():
			return nil
		case obj = <-c.writeChan:
		}
		encoder := json.NewEncoder(c.w)
		if err := encoder.Encode(obj); err != nil {
			return err
		}
		if err := c.w.Flush(); err != nil {
			return err
		}
	}
}
