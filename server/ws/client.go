package ws

import (
	"context"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"gitlab.com/lattetalk/lattetalk/ws/data"
)

type Client struct {
	conns []*Conn
	ID    int64
	mu    sync.Mutex
}

func (c *Client) Broadcast(method data.ResponseMethod, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, conn := range c.conns {
		conn.Output(method, data)
	}
}

type CancelWrite func()

func CancelableWrite(conn *Conn) CancelWrite {
	ctx, cancel := context.WithCancel(context.Background())
	e := make(chan error)
	go func() {
		e <- conn.Write(ctx)
	}()
	return func() {
		cancel()
		<-e
	}
}

func (c *Client) AddConn(rwc net.Conn) (*Conn, CancelWrite, error) {
	conn := new(Conn)
	conn.rwc = rwc
	conn.subject = c.ID
	conn.w = wsutil.NewWriter(conn.rwc, ws.StateServerSide, ws.OpText)
	conn.writeChan = make(chan interface{})

	c.mu.Lock()
	conn.ID = len(c.conns)
	c.conns = append(c.conns, conn)
	c.mu.Unlock()

	// Start the Write end of Conn. When an error occurred, close the connection.
	cancel := CancelableWrite(conn)

	return conn, cancel, nil
}

func (c *Client) RemoveConn(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	t := c.conns[len(c.conns)-1]
	t.ID = id

	c.conns[id] = t
	c.conns[len(c.conns)-1] = nil
	c.conns = c.conns[:len(c.conns)-1]
}
