#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "Starting build process for Shuffle-Go..."

# Ensure dependencies are synchronized
echo "Updating dependencies..."
go mod tidy

# Build the CLI version
echo "Building shuffle-cli..."
go build -ldflags="-s -w" -o shuffle-cli cmd/shuffle-cli/main.go

# Build the GUI version
# Note: This requires CGO and graphics development headers (e.g., libgl1-mesa-dev on Linux)
echo "Building shuffle-gui..."
go build -ldflags="-s -w" -o shuffle-gui cmd/shuffle-gui/main.go

echo "--------------------------------------------------"
echo "Build complete!"
echo "Binaries created:"
echo "  - ./shuffle-cli (Command Line Interface)"
echo "  - ./shuffle-gui (Graphical User Interface)"
echo "--------------------------------------------------"
