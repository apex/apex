
release:
	@echo "[+] releasing $(VERSION)"
	@echo "[+] re-generating bin data"
	@go generate ./...
	@echo "[+] committing"
	@git release $(VERSION)
	@echo "[+] building"
	@$(MAKE) build
	@echo "[+] complete"
.PHONY: release

test:
	@go test -cover ./...
.PHONY: test

build:
	@gox -os="linux darwin windows" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean
