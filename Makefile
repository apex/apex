
# Release the given VERSION.
release:
	@echo "[+] releasing $(VERSION)"
	@echo "[+] re-generating"
	@go generate ./...
	@echo "[+] building"
	@$(MAKE) build
	@echo "[+] comitting"
	@git release $(VERSION)
	@echo "[+] complete"
.PHONY: release

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Build release binaries.
build:
	@gox -os="linux darwin windows openbsd" ./...
.PHONY: build

# Clean build artifacts.
clean:
	@git clean -f
.PHONY: clean
