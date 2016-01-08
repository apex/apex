
test:
	@go test -cover -v ./...
.PHONY: test

build:
	@gox -os="linux darwin windows" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean

completion-scripts:
	@docopt-completion apex --manual-bash > /dev/null
	@mv apex.sh ./contrib/completion/apex-completion.bash
	@docopt-completion apex --manual-zsh > /dev/null
	@mv _apex ./contrib/completion/apex-completion.zsh
.PHONY: completion-scripts