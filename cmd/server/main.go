package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nipunap/sqlite-mcp-server/internal/db"
	"github.com/nipunap/sqlite-mcp-server/internal/mcp"
)

func main() {
	// Parse flags
	registryPath := flag.String("registry", "registry.db", "Path to database registry")
	defaultDB := flag.String("db", "", "Default database to register (optional)")
	flag.Parse()

	// Create absolute path for registry
	absRegistryPath, err := filepath.Abs(*registryPath)
	if err != nil {
		log.Fatalf("Failed to resolve registry path: %v", err)
	}

	// Set up database registry
	registry, err := db.NewRegistry(absRegistryPath)
	if err != nil {
		log.Fatalf("Failed to create database registry: %v", err)
	}
	defer registry.Close()

	// Create database manager
	manager := db.NewManager(registry)
	defer manager.CloseAll()

	// Register default database if provided
	if *defaultDB != "" {
		absDefaultDB, err := filepath.Abs(*defaultDB)
		if err != nil {
			log.Fatalf("Failed to resolve default database path: %v", err)
		}

		defaultInfo := &db.DatabaseInfo{
			ID:          "default",
			Name:        "default",
			Path:        absDefaultDB,
			Description: "Default database",
			ReadOnly:    false,
			Owner:       "system",
			Status:      "active",
		}

		if err := registry.RegisterDatabase(defaultInfo); err != nil {
			log.Printf("Warning: Failed to register default database (might already exist): %v", err)
		} else {
			log.Printf("Registered default database: %s", absDefaultDB)
		}
	}

	// Create MCP server
	server, err := mcp.NewServer(manager)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Run server using STDIO transport
	if err := server.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
