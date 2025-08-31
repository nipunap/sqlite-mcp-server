package mcp

import (
	"encoding/json"
	"fmt"
)

// Capability represents a server capability
type Capability struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // tool, resource, or prompt
	Description string      `json:"description"`
	Schema      interface{} `json:"schema,omitempty"`
}

// CapabilityRegistry manages server capabilities
type CapabilityRegistry struct {
	tools     map[string]ToolHandler
	resources map[string]ResourceHandler
	prompts   map[string]string
}

// ToolHandler handles tool invocations
type ToolHandler func(params json.RawMessage) (interface{}, error)

// ResourceHandler provides resource content
type ResourceHandler func() (interface{}, error)

// NewCapabilityRegistry creates a new capability registry
func NewCapabilityRegistry() *CapabilityRegistry {
	return &CapabilityRegistry{
		tools:     make(map[string]ToolHandler),
		resources: make(map[string]ResourceHandler),
		prompts:   make(map[string]string),
	}
}

// RegisterTool registers a new tool capability
func (r *CapabilityRegistry) RegisterTool(name string, handler ToolHandler, schema interface{}) error {
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}
	r.tools[name] = handler
	return nil
}

// RegisterResource registers a new resource capability
func (r *CapabilityRegistry) RegisterResource(name string, handler ResourceHandler) error {
	if _, exists := r.resources[name]; exists {
		return fmt.Errorf("resource %s already registered", name)
	}
	r.resources[name] = handler
	return nil
}

// RegisterPrompt registers a new prompt capability
func (r *CapabilityRegistry) RegisterPrompt(name string, content string) error {
	if _, exists := r.prompts[name]; exists {
		return fmt.Errorf("prompt %s already registered", name)
	}
	r.prompts[name] = content
	return nil
}

// GetCapabilities returns all registered capabilities
func (r *CapabilityRegistry) GetCapabilities() []Capability {
	var caps []Capability

	// Add tools
	for name := range r.tools {
		caps = append(caps, Capability{
			Name:        name,
			Type:        "tool",
			Description: fmt.Sprintf("Tool: %s", name),
		})
	}

	// Add resources
	for name := range r.resources {
		caps = append(caps, Capability{
			Name:        name,
			Type:        "resource",
			Description: fmt.Sprintf("Resource: %s", name),
		})
	}

	// Add prompts
	for name := range r.prompts {
		caps = append(caps, Capability{
			Name:        name,
			Type:        "prompt",
			Description: fmt.Sprintf("Prompt: %s", name),
		})
	}

	return caps
}

// HandleCapabilityRequest processes a capability request
func (r *CapabilityRegistry) HandleCapabilityRequest(msg *JSONRPCMessage) *JSONRPCMessage {
	switch msg.Method {
	case "capabilities":
		return &JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Result:  r.GetCapabilities(),
		}
	case "invoke":
		var params struct {
			Name   string          `json:"name"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal(msg.Params, &params); err != nil {
			return &JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Error: &JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
				},
			}
		}

		// Handle based on capability type
		if handler, ok := r.tools[params.Name]; ok {
			result, err := handler(params.Params)
			if err != nil {
				return &JSONRPCMessage{
					Version: "2.0",
					ID:      msg.ID,
					Error: &JSONRPCError{
						Code:    -32000,
						Message: err.Error(),
					},
				}
			}
			return &JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Result:  result,
			}
		}

		if handler, ok := r.resources[params.Name]; ok {
			result, err := handler()
			if err != nil {
				return &JSONRPCMessage{
					Version: "2.0",
					ID:      msg.ID,
					Error: &JSONRPCError{
						Code:    -32000,
						Message: err.Error(),
					},
				}
			}
			return &JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Result:  result,
			}
		}

		if content, ok := r.prompts[params.Name]; ok {
			return &JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Result:  content,
			}
		}

		return &JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Capability not found",
			},
		}
	default:
		return &JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}
