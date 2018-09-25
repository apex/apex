
GO ?= go

# Provide some targets for external dep management
BIN_DIR := $(GOPATH)/bin
MOCKGEN := $(BIN_DIR)/mockgen

# Build all files.
build: $(MOCKGEN)
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
	@goreleaser -p 1 --rm-dist --config .goreleaser.yml
	@echo "==> Complete"
.PHONY: release

# Clean build artifacts.
clean:
	@git clean -f
.PHONY: clean

$(MOCKGEN):
	@echo "==> Installing Mockgen"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen

local:
	go install -a -ldflags "-X github.com/apex/apex/cmd/apex/version.Version=development" github.com/apex/apex/cmd/apex/
.PHONY: local
