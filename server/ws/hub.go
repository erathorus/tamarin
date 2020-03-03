package ws

import (
	"sync"
)

type ClientHub struct {
	clients map[int64]*Client
	mu      sync.Mutex
}

func (h *ClientHub) GetClient(id int64) *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	c, ok := h.clients[id]

	if !ok {
		c = &Client{ID: id}
		h.clients[id] = c
	}

	return c
}

var Hub *ClientHub

func init() {
	Hub = &ClientHub{
		clients: make(map[int64]*Client),
	}
}
