package mcp

import (
	"context"

	"github.com/nipunap/sqlite-mcp-server/internal/db"
	"github.com/nipunap/sqlite-mcp-server/internal/mcp/prompts"
	"github.com/nipunap/sqlite-mcp-server/internal/mcp/resources"
	"github.com/nipunap/sqlite-mcp-server/internal/mcp/tools"
)

// Server implements the MCP server
type Server struct {
	manager   *db.Manager
	registry  *CapabilityRegistry
	transport *STDIOTransport
}

// NewServer creates a new MCP server instance
func NewServer(manager *db.Manager) (*Server, error) {
	s := &Server{
		manager:   manager,
		registry:  NewCapabilityRegistry(),
		transport: NewSTDIOTransport(),
	}

	// Initialize components
	dbTools := tools.NewDBTools(manager)
	dbResources := resources.NewDBResources(manager)

	// Register database management tools
	if err := s.registry.RegisterTool("db/register_database", dbTools.RegisterDatabase, nil); err != nil {
		return nil, err
	}
	if err := s.registry.RegisterTool("db/list_databases", dbTools.ListDatabases, nil); err != nil {
		return nil, err
	}

	// Register database operation tools
	if err := s.registry.RegisterTool("db/get_table_schema", dbTools.GetTableSchema, nil); err != nil {
		return nil, err
	}
	if err := s.registry.RegisterTool("db/insert_record", dbTools.InsertRecord, nil); err != nil {
		return nil, err
	}
	if err := s.registry.RegisterTool("db/query", dbTools.ExecuteQuery, nil); err != nil {
		return nil, err
	}

	// Register database query tools (previously resources, but they need parameters)
	if err := s.registry.RegisterTool("db/get_tables", dbResources.GetTables, nil); err != nil {
		return nil, err
	}
	if err := s.registry.RegisterTool("db/get_schema", dbResources.GetSchema, nil); err != nil {
		return nil, err
	}

	// Register resources (no parameters needed)
	if err := s.registry.RegisterResource("db/databases", dbResources.GetDatabases); err != nil {
		return nil, err
	}

	// Register prompts
	for name, content := range prompts.DBPrompts {
		if err := s.registry.RegisterPrompt(name, content); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Run starts the MCP server
func (s *Server) Run(ctx context.Context) error {
	return s.transport.HandleMessages(ctx, s.handleMessage)
}

// handleMessage processes incoming MCP messages
func (s *Server) handleMessage(msg *JSONRPCMessage) *JSONRPCMessage {
	return s.registry.HandleCapabilityRequest(msg)
}
