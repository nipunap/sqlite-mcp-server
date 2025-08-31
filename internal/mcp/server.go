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
	s.registry.RegisterTool("db/register_database", dbTools.RegisterDatabase, nil)
	s.registry.RegisterTool("db/list_databases", dbTools.ListDatabases, nil)

	// Register database operation tools
	s.registry.RegisterTool("db/get_table_schema", dbTools.GetTableSchema, nil)
	s.registry.RegisterTool("db/insert_record", dbTools.InsertRecord, nil)
	s.registry.RegisterTool("db/query", dbTools.ExecuteQuery, nil)

	// Register database query tools (previously resources, but they need parameters)
	s.registry.RegisterTool("db/get_tables", dbResources.GetTables, nil)
	s.registry.RegisterTool("db/get_schema", dbResources.GetSchema, nil)

	// Register resources (no parameters needed)
	s.registry.RegisterResource("db/databases", dbResources.GetDatabases)

	// Register prompts
	for name, content := range prompts.DBPrompts {
		s.registry.RegisterPrompt(name, content)
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
