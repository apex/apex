
GO ?= go

# Build all files.
build:
	@echo "==> Building"
	@$(GO) generate ./...
.PHONY: build

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Release binaries to GitHub.
release:
	@echo "==> Releasing"
	@goreleaser -p 1 --rm-dist -config .goreleaser.yml
	@echo "==> Complete"
.PHONY: release

# Clean build artifacts.
clean:
	@git clean -f
.PHONY: clean
