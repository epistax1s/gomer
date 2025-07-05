#!/bin/bash

# Debug script for Gomer project
# Requires delve to be installed: go install github.com/go-delve/delve/cmd/dlv@latest

echo "Starting Gomer in debug mode..."

# Set environment variables
export GO_ENV=development

# Run with delve debugger
dlv debug ./cmd/gomer/main.go --headless --listen=:2345 --api-version=2 --accept-multiclient

echo "Debug session ended." 