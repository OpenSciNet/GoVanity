#!/bin/bash

# Create a distribution folder
mkdir -p dist

APP_NAME="govanity"

echo "🚀 Starting Cross-Platform Build..."

# 1. Build for Linux (64-bit)
echo "📦 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o dist/${APP_NAME}-linux-amd64 .

# 2. Build for Windows (64-bit)
echo "📦 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o dist/${APP_NAME}-windows-amd64.exe .

# 3. Build for macOS (Intel & Apple Silicon)
echo "📦 Building for Mac (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o dist/${APP_NAME}-darwin-amd64 .

echo "📦 Building for Mac (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o dist/${APP_NAME}-darwin-arm64 .

echo "✅ Done! Check the 'dist' folder."
ls -l dist