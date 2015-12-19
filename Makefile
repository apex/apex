
build:
	@gox -os="linux darwin windows" ./...
.PHONY: