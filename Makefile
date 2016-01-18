
test:
	@go test -cover ./...
.PHONY: test

build:
	@gox -os="linux darwin windows" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean
