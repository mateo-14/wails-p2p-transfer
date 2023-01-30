package p2p

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageID string

type Message struct {
	ID      MessageID
	Payload interface{}
}

type MessageRequestHandler func(*Message, io.Writer)

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
	err := messagebToMessage(m.s, &msg)
	if err != nil {
		fmt.Println("Error reading from stream: ", err)
		return
	}

	handler, ok := m.handlers[msg.ID]
	if !ok {
		runtime.LogErrorf(ctx, "SendMessage: No handler for message: %s\n", msg.ID)
		m.s.Close()
		return
	}

	handler(&msg, m.s)
	m.s.Close()
}

func (m *MessageHandler) HandleRequest(msgID MessageID, handler MessageRequestHandler) {
	fmt.Printf("Registering handler for message: %s\n", msgID)
	m.handlers[MessageID(msgID)] = handler
}

func messageToBytes(msg *Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(msg)
	return buf.Bytes(), err
}

func messagebToMessage(r io.Reader, msg *Message) error {
	dec := gob.NewDecoder(r)
	err := dec.Decode(msg)
	return err
}
