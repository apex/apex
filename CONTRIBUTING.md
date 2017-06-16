
# Pull Requests

When creating a pull-request you should:

- __Open an issue first:__ Confirm that the change or feature will be accepted
- __Lint your code:__ Use  `gofmt`, `golint`, and `govet` to clean up your code
- __Squash multiple commits:__ Squash multiple commits into a single commit via `git rebase -i`
- __Start message with a verb:__ Your commit message must start a lowercase verb such as "add", "fix", "refactor", "remove"
- __Reference the issue__: Ensure that your commit message references the issue with ". Closes #N"
- __Add to feature list__: If your pull-request is for a feature, make sure to add it to the Readme's feature list
- __Add a GIF__

# Updating dependencies

Currently Apex does not vendor deps, so you should update them with:

    $ go get -u ./...

# Running Apex in development

Apex requires Go 1.6 or higher.

To run Apex in development execute:

    $ go run cmd/apex/main.go <args>

To install Apex from the repo, execute:

    $ go install ./...

Note that this will install to $GOPATH/bin, which must be in your $PATH.

# Running Tests
To run the test for the project:

    $ go test -v ./...

# Running the Linter
To run the linter for the project:

    $ go vet ./...


# Testing Bash Auto-completion

Some examples for testing bash auto-completion. The value after `--` is
the current word. Or run `apex autocomplete` on its own to view the valid list
candidate words.

```
$ compgen -W "$(apex autocomplete)" --
$ compgen -W "$(apex autocomplete)" -- depl
$ compgen -W "$(apex autocomplete deploy)" --
$ compgen -W "$(apex autocomplete deploy)" -- api_
```
