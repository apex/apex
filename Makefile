
test:
	@go test -cover -v ./...
.PHONY: test

build:
	@gox -os="linux darwin windows" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean