package p2p

import (
	"context"
	"encoding/gob"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RequestID int64

type ResponseData struct {
	ID RequestID
}

type RequestData struct {
	ID RequestID
}

type Response struct {
	ResponseData
	Body io.ReadCloser
}

type Request struct {
	RequestData
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
	var reqd RequestData
	err := bytesToStruct(m.s, &reqd)
	if err != nil {
		runtime.LogErrorf(ctx, "MessageHandler: Error reading from stream: %s\n", err)
		return
	}

	defer m.s.Close()

	handler, ok := m.handlers[reqd.ID]

	if !ok {
		runtime.LogErrorf(ctx, "MessageHandler: No handler for message: %s\n", reqd.ID)
		return
	}

	runtime.LogInfof(ctx, "MessageHandler: Handler for message (%s) found. Message handled.\n", reqd.ID)

	req := Request{
		RequestData: reqd,
	}
	req.WriteCloser = m.s
	req.Body = m.s
	handler(&req)
}

func (m *MessageHandler) HandleRequest(msgID RequestID, handler MessageRequestHandler) {
	m.handlers[RequestID(msgID)] = handler
}

func structToBytes[T any](req *T, w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(req)
}

func bytesToStruct[T any](r io.Reader, res *T) error {
	dec := gob.NewDecoder(r)
	err := dec.Decode(res)
	return err
}

func (m *Request) Write(r io.Reader) (int64, error) {
	res := ResponseData{
		ID: m.ID,
	}

	err := structToBytes(&res, m.WriteCloser)
	if err != nil {
		return 0, err
	}

	return io.Copy(m.WriteCloser, r)
}
