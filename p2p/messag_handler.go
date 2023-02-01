package p2p

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RequestID string

const (
	GetFiles RequestID = "get/files"
)

type Response struct {
	ID   RequestID
	Body io.ReadCloser
}

type Request struct {
	ID   RequestID
	Body io.ReadCloser
	io.WriteCloser
}

type MessageRequestHandler func(*Request)

// Handle incoming messages
type MessageHandler struct {
	s        network.Stream
	handlers map[RequestID]MessageRequestHandler
}

// Creates a new message handler
func NewMessageHandler(ctx context.Context, s network.Stream) *MessageHandler {
	msgh := &MessageHandler{
		s:        s,
		handlers: make(map[RequestID]MessageRequestHandler),
	}

	go msgh.handle(ctx)

	return msgh
}

func (m *MessageHandler) handle(ctx context.Context) {
	var req Request
	err := bytesToStruct(m.s, &req)
	if err != nil {
		runtime.LogErrorf(ctx, "MessageHandler: Error reading from stream: %s\n", err)
		return
	}

	handler, ok := m.handlers[req.ID]
	if !ok {
		runtime.LogErrorf(ctx, "MessageHandler: No handler for message: %s\n", req.ID)
		m.s.Close()
		return
	}

	req.WriteCloser = m.s
	req.Body = m.s
	runtime.LogInfof(ctx, "MessageHandler: Handler for message (%s) found. Message handled.\n", req.ID)
	handler(&req)

	err = m.s.Close()
	if err != nil {
		runtime.LogErrorf(ctx, "MessageHandler: Error closing stream: %s\n", err)
	}

	runtime.LogInfof(ctx, "MessageHandler: Stream for message (%s) closed.\n", req.ID)
}

func (m *MessageHandler) HandleRequest(msgID RequestID, handler MessageRequestHandler) {
	m.handlers[RequestID(msgID)] = handler
}

func structToBytes[T any](req *T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(req)
	return buf.Bytes(), err
}

func bytesToStruct[T any](r io.Reader, res *T) error {
	dec := gob.NewDecoder(r)
	err := dec.Decode(res)
	return err
}

func (m *Request) Write(body []byte) error {
	res := Response{
		ID: m.ID,
	}

	resb, err := structToBytes(&res)
	if err != nil {
		return err
	}

	m.Write(resb)
	m.Write(body)
	return nil
}
