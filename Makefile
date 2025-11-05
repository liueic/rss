.PHONY: build test clean run fmt vet help

# Default target
help:
	@echo "RSS Watcher - Available targets:"
	@echo "  make build    - Build the rsswatcher binary"
	@echo "  make test     - Run tests"
	@echo "  make run      - Run rsswatcher with default config"
	@echo "  make fmt      - Format Go code"
	@echo "  make vet      - Run go vet"
	@echo "  make clean    - Remove built binaries"
	@echo "  make all      - Format, vet, test, and build"

# Build the binary
build:
	go build -v -o rsswatcher ./cmd/rsswatcher

# Run tests
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run the application
run: build
	./rsswatcher --config feeds.yaml --state state/last_states.json

# Format code
fmt:
	gofmt -w .

# Run go vet
vet:
	go vet ./...

# Clean built binaries
clean:
	rm -f rsswatcher
	rm -f coverage.txt

# Run all checks and build
all: fmt vet test build
	@echo "All checks passed!"
