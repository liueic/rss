#!/bin/bash

# Local testing script for RSS Watcher
# Usage: ./test-local.sh [your-bark-device-key]

set -e

echo "RSS Watcher - Local Test Script"
echo "================================"
echo ""

# Check if Bark device key is provided
if [ -z "$1" ]; then
    if [ -z "$BARK_DEVICE_KEY" ]; then
        echo "Error: BARK_DEVICE_KEY not provided"
        echo ""
        echo "Usage: ./test-local.sh YOUR_DEVICE_KEY"
        echo "   or: export BARK_DEVICE_KEY=YOUR_KEY && ./test-local.sh"
        exit 1
    fi
else
    export BARK_DEVICE_KEY="$1"
fi

echo "✓ Bark device key configured"

# Build the application
echo "Building rsswatcher..."
go build -o rsswatcher ./cmd/rsswatcher
echo "✓ Build successful"

# Check if test config exists
if [ ! -f "feeds.yaml.test" ]; then
    echo "Creating test configuration..."
    cat > feeds.yaml.test << 'EOF'
feeds:
  - id: go-blog
    name: Go Blog
    url: https://go.dev/blog/feed.atom
    notify: true
    dedupe_key: guid
    aggregate: false
EOF
    echo "✓ Test configuration created"
fi

# Create test state directory
mkdir -p state-test
echo "✓ Test state directory ready"

# Run the watcher
echo ""
echo "Running RSS Watcher with test configuration..."
echo "----------------------------------------------"
./rsswatcher --config feeds.yaml.test --state state-test/test_states.json

echo ""
echo "✓ Test completed successfully!"
echo ""
echo "Check your iOS device for notifications."
echo "State saved to: state-test/test_states.json"
