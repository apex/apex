
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

test:
	@go test -cover ./...
.PHONY: test

build:
	@gox -os="linux darwin windows openbsd" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean
