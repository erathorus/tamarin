package ws

import (
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"gitlab.com/horo-go/auth0"
	"gitlab.com/lattetalk/lattetalk/auth"
)

const (
	ioTimeOut     = time.Second
	ioMaxRead     = 1 << 20
	retryDuration = 5 * time.Millisecond
)

func handleConn(rwc net.Conn) error {
	var subject int64
	u := ws.Upgrader{
		OnHeader: func(key, value []byte) error {
			if string(key) == "Cookie" {
				var authToken string
				ok := httphead.ScanCookie(value, func(key, value []byte) bool {
					if string(key) == auth.TokenCookie {
						authToken = string(value)
						return false
					}
					return true
				})
				if ok {
					claims, err := auth.Auth0.GetClaims(authToken, auth0.HS256)
					if err == nil {
						subject, err = strconv.ParseInt(claims.Subject, 10, 64)
						if err == nil {
							return nil
						}
					}
				}
				return ws.RejectConnectionError(
					ws.RejectionStatus(400),
					ws.RejectionReason("bad cookie"),
				)
			}
			return nil
		},
	}

	safeConn := ioLimiter{rwc, ioTimeOut, ioMaxRead}
	_, err := u.Upgrade(safeConn)
	if err != nil {
		return err
	}

	client := Hub.GetClient(subject)

	conn, cancel, err := client.AddConn(safeConn)
	if err != nil {
		return err
	}

	desc := netpoll.Must(netpoll.HandleRead(rwc))

	// Subscribe to events about new data to Read.
	poller.Start(desc, func(event netpoll.Event) {
		ok := true
		defer func() {
			// If an error is occurred while reading or the connection is closed,
			// clean up.
			if !ok {
				// Stop watcher
				poller.Stop(desc)
				// Call cancel to cancel Write() goroutine
				cancel()
				// Close the connection
				rwc.Close()
				// Remove the connection from Client list
				client.RemoveConn(conn.ID)
			}
		}()

		// Handle connection is closed
		if event&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			ok = false
			return
		}

		// Read new data
		if err := conn.Read(); err != nil {
			ok = false
		}
	})

	return nil
}

// Listen listens to a tcp socket and upgrade tcp connection to WebSocket.
func Listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// Create netpoll descriptor for the listener.
	acceptDesc := netpoll.Must(netpoll.HandleListener(ln, netpoll.EventRead))

	// Subscribe to events about listener.
	poller.Start(acceptDesc, func(event netpoll.Event) {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Printf("accept error: %v; retrying in %s", err, retryDuration)
				time.Sleep(retryDuration)
				return
			}
			log.Fatalf("accept error: %v", err)
		}

		go func() {
			if err = handleConn(conn); err != nil {
				conn.Close()
			}
		}()
	})
}

type ioLimiter struct {
	net.Conn
	timeLimit time.Duration
	readLimit int64
}

func (l ioLimiter) Read(p []byte) (int, error) {
	if err := l.Conn.SetReadDeadline(time.Now().Add(l.timeLimit)); err != nil {
		return 0, err
	}
	r := io.LimitReader(l.Conn, l.readLimit)
	return r.Read(p)
}

func (l ioLimiter) Write(p []byte) (int, error) {
	if err := l.Conn.SetWriteDeadline(time.Now().Add(l.timeLimit)); err != nil {
		return 0, err
	}
	return l.Conn.Write(p)
}
