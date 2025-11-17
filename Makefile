.PHONY: build test clean help version

# Get version from git
VERSION := $(shell git describe --tags --always)
COMMIT_HASH := $(shell git rev-parse --short HEAD)

help:
	@echo "tablefy - Makefile targets:"
	@echo ""
	@echo "  make build      - Build the project with version info"
	@echo "  make build-dev  - Build the project without version info (development)"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make version    - Show current version"
	@echo ""

build:
	@echo "Building tablefy $(VERSION) (commit: $(COMMIT_HASH))..."
	go build -ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(COMMIT_HASH)" -o bin/tablefy ./cmd/tablefy
	@echo "✅ Build complete: bin/tablefy"

build-dev:
	@echo "Building tablefy (development mode)..."
	go build -o bin/tablefy ./cmd/tablefy
	@echo "✅ Build complete: bin/tablefy (version: dev)"

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -f bin/tablefy
	@echo "✅ Clean complete"

version:
	@echo "Current version: $(VERSION)"
	@echo "Current commit: $(COMMIT_HASH)"
