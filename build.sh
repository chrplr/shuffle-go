#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Extract version from git tag
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
MODULE="github.com/chrplr/shuffle-go"

echo "Starting build process for Shuffle-Go (Version: $VERSION)..."

# Ensure dependencies are synchronized
echo "Updating dependencies..."
go mod tidy

LDFLAGS="-s -w -X '$MODULE/internal/version.Version=$VERSION'"

# Build the CLI version for local OS
echo "Building local shuffle-cli..."
go build -ldflags="$LDFLAGS" -o shuffle-cli ./cmd/shuffle-cli

# Build the GUI version for local OS
# Note: This requires CGO and graphics development headers (e.g., libgl1-mesa-dev on Linux)
echo "Building local shuffle-gui..."
go build -ldflags="$LDFLAGS" -o shuffle-gui ./cmd/shuffle-gui

# Cross-compilation examples using fyne-cross (if installed)
if command -v fyne-cross >/dev/null 2>&1; then
    echo "Fyne-cross detected. Building for other platforms..."
    
    # Example: Build for Linux and Windows (amd64)
    # fyne-cross linux -arch amd64 -name shuffle-gui -ldflags "$LDFLAGS" ./cmd/shuffle-gui
    # fyne-cross windows -arch amd64 -name shuffle-gui -ldflags "$LDFLAGS" ./cmd/shuffle-gui
    
    echo "Note: fyne-cross commands are commented out in build.sh by default."
    echo "Uncomment them to enable automated cross-compilation."
else
    echo "Fyne-cross not found. Skipping cross-compilation steps."
fi

echo "--------------------------------------------------"
echo "Build complete!"
echo "Binaries created:"
echo "  - ./shuffle-cli (Command Line Interface)"
echo "  - ./shuffle-gui (Graphical User Interface)"
echo "--------------------------------------------------"
