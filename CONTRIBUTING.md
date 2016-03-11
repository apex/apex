
# Pull Requests

When creating a pull-request you should:

- __Open an issue first:__ Confirm that the change or feature will be accepted
- __Lint your code:__ Use  `gofmt`, `golint`, and `govet` to clean up your code
- __Squash multiple commits:__ Squash multiple commits into a single commit via `git rebase -i`
- __Start message with a verb:__ Your commit message must start a lowercase verb such as "add", "fix", "refactor", "remove"
- __Reference the issue__: Ensure that your commit message references the issue with ". Closes #N"
- __Add to feature list__: If your pull-request is for a feature, make sure to add it to the Readme's feature list

# Updating dependencies

Currently Apex does not vendor deps, so you should update them with:

  $ go get -u ./...

# Running Apex in development

To run Apex in development execute:

  $ go run cmd/apex/main.go <args>

To install Apex from the repo, execute:

  $ go install ./...

Note that this will install to $GOPATH/bin, which must be in your $PATH.
