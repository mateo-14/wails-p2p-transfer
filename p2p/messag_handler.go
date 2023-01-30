package p2p

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageID string

type Message struct {
	ID      MessageID
	Payload interface{}
}

type MessageRequestHandler func(*MessageRequest)

// Handle incoming messages
type MessageHandler struct {
	s        network.Stream
	handlers map[MessageID]MessageRequestHandler
}

// Creates a new message handler
func NewMessageHandler(ctx context.Context, s network.Stream) *MessageHandler {
	msgh := &MessageHandler{
		s:        s,
		handlers: make(map[MessageID]MessageRequestHandler),
	}

	go msgh.handle(ctx)

	return msgh
}

func (m *MessageHandler) handle(ctx context.Context) {
	var msg Message
	err := messageBToMessage(m.s, &msg)
	if err != nil {
		runtime.LogErrorf(ctx, "MessageHandler: Error reading from stream: %s\n", err)
		return
	}

	handler, ok := m.handlers[msg.ID]
	if !ok {
		runtime.LogErrorf(ctx, "MessageHandler: No handler for message: %s\n", msg.ID)
		m.s.Close()
		return
	}

	runtime.LogInfof(ctx, "MessageHandler: Handler for message (%s) found. Message handled.\n", msg.ID)
	handler(&MessageRequest{
		&msg,
		m.s,
	})

	err = m.s.Close()
	if err != nil {
		runtime.LogErrorf(ctx, "MessageHandler: Error closing stream: %s\n", err)
	}

	runtime.LogInfof(ctx, "MessageHandler: Stream for message (%s) closed.\n", msg.ID)
}

func (m *MessageHandler) HandleRequest(msgID MessageID, handler MessageRequestHandler) {
	m.handlers[MessageID(msgID)] = handler
}

func messageToBytes(msg *Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(msg)
	return buf.Bytes(), err
}

func messageBToMessage(r io.Reader, msg *Message) error {
	dec := gob.NewDecoder(r)
	err := dec.Decode(msg)
	return err
}

type MessageRequest struct {
	Message *Message
	w       io.Writer
}

func (m *MessageRequest) Write(payload interface{}) error {
	res := Message{
		ID:      m.Message.ID,
		Payload: payload,
	}

	resb, err := messageToBytes(&res)
	if err != nil {
		return err
	}

	m.w.Write(resb)
	return nil
}
