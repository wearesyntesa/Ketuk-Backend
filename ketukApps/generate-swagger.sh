#!/bin/bash

# Swagger Documentation Generator Script
# This script regenerates Swagger documentation for KetukApps API

echo "🔄 Regenerating Swagger documentation..."

# Navigate to the project directory
cd "$(dirname "$0")"

# Generate Swagger docs
if [ -f "$HOME/go/bin/swag" ]; then
    $HOME/go/bin/swag init
    echo "✅ Swagger documentation generated successfully!"
    echo "📚 Documentation files created in ./docs/"
    echo ""
    echo "🚀 To view the documentation:"
    echo "   1. Start the server: go run main.go"
    echo "   2. Open: http://localhost:8081/swagger/index.html"
else
    echo "❌ Error: swag command not found"
    echo "📦 Please install swag first:"
    echo "   go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi
