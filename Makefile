
test:
	@go test -cover ./...
.PHONY: test

cov:
	@go test -covermode=count -coverprofile=cov.out ./$(PKG)
	@go tool cover -html=cov.out
	@rm cov.out
.PHONY: cov

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