package ws

import (
	"github.com/mailru/easygo/netpoll"
)

var poller netpoll.Poller

func init() {
	var err error
	poller, err = netpoll.New(nil)
	if err != nil {
		panic(err)
	}
}
