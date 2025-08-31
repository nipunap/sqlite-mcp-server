package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JSONRPCMessage represents a JSON-RPC 2.0 message
type JSONRPCMessage struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method,omitempty"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *JSONRPCError    `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// STDIOTransport handles STDIO-based JSON-RPC communication
type STDIOTransport struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewSTDIOTransport creates a new STDIO transport
func NewSTDIOTransport() *STDIOTransport {
	return &STDIOTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

// ReadMessage reads a JSON-RPC message from STDIN
func (t *STDIOTransport) ReadMessage() (*JSONRPCMessage, error) {
	line, err := t.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read message: %v", err)
	}

	var msg JSONRPCMessage
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %v", err)
	}

	return &msg, nil
}

// WriteMessage writes a JSON-RPC message to STDOUT
func (t *STDIOTransport) WriteMessage(msg *JSONRPCMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	if _, err := t.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	if err := t.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("failed to write newline: %v", err)
	}

	return t.writer.Flush()
}

// HandleMessages processes incoming messages until context is cancelled
func (t *STDIOTransport) HandleMessages(ctx context.Context, handler func(*JSONRPCMessage) *JSONRPCMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := t.ReadMessage()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			response := handler(msg)
			if response != nil {
				if err := t.WriteMessage(response); err != nil {
					return err
				}
			}
		}
	}
}
