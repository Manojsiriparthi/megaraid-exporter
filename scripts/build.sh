#!/bin/bash

set -e

echo "Building MegaRAID Exporter..."

# Build for current platform
go build -o bin/megaraid-exporter ./cmd/exporter

# Build for Linux (if cross-compiling)
GOOS=linux GOARCH=amd64 go build -o bin/megaraid-exporter-linux ./cmd/exporter

echo "Build completed successfully!"
echo "Binary location: bin/megaraid-exporter"
