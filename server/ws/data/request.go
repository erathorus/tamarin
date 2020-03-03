package data

import (
	"encoding/json"
)

type RequestMethod string

const (
	// Receive a new message from a client.
	RequestNewMessage = RequestMethod("new_message")
)

type Request struct {
	Method  RequestMethod   `json:"method"`
	Data    json.RawMessage `json:"data"`
	subject int64
}

func (r *Request) SetSubject(subject int64) {
	r.subject = subject
}

func (r *Request) Subject() int64 {
	return r.subject
}
