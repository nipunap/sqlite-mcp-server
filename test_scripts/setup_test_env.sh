#!/bin/bash

# Setup test environment for MCP server manual testing

set -e

echo "🚀 Setting up MCP Server Test Environment"
echo "========================================"

# Build the server
echo "1. Building MCP server..."
go build -o sqlite-mcp-server cmd/server/main.go
echo "✅ Server built successfully"

# Create test database
echo "2. Creating test database..."
sqlite3 test_manual.db << 'EOF'
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    price REAL,
    category TEXT,
    in_stock BOOLEAN DEFAULT 1
);

CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    product_id INTEGER,
    quantity INTEGER,
    order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Insert sample data
INSERT INTO users (name, email) VALUES
    ('Alice Johnson', 'alice@example.com'),
    ('Bob Smith', 'bob@example.com'),
    ('Carol Davis', 'carol@example.com');

INSERT INTO products (name, price, category) VALUES
    ('Laptop', 999.99, 'Electronics'),
    ('Coffee Mug', 12.50, 'Kitchen'),
    ('Book: Go Programming', 29.99, 'Books'),
    ('Wireless Mouse', 45.00, 'Electronics');

INSERT INTO orders (user_id, product_id, quantity) VALUES
    (1, 1, 1),
    (2, 2, 2),
    (1, 3, 1),
    (3, 4, 1);
EOF

echo "✅ Test database created with sample data"

# Create a second test database
echo "3. Creating second test database..."
sqlite3 inventory.db << 'EOF'
CREATE TABLE warehouses (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    location TEXT,
    capacity INTEGER
);

CREATE TABLE stock (
    id INTEGER PRIMARY KEY,
    product_name TEXT NOT NULL,
    warehouse_id INTEGER,
    quantity INTEGER,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

INSERT INTO warehouses (name, location, capacity) VALUES
    ('Main Warehouse', 'New York', 10000),
    ('West Coast Hub', 'California', 8000);

INSERT INTO stock (product_name, warehouse_id, quantity) VALUES
    ('Laptop', 1, 50),
    ('Coffee Mug', 1, 200),
    ('Laptop', 2, 30),
    ('Wireless Mouse', 2, 75);
EOF

echo "✅ Second test database (inventory.db) created"

# Set permissions
chmod +x sqlite-mcp-server
echo "✅ Executable permissions set"

echo ""
echo "🎉 Test environment setup complete!"
echo ""
echo "Available test databases:"
echo "  - test_manual.db (users, products, orders)"
echo "  - inventory.db (warehouses, stock)"
echo ""
echo "To start testing, run:"
echo "  ./test_scripts/run_basic_tests.sh"
echo "  ./test_scripts/run_interactive_test.py"
