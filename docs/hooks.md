
Apex supports the notion of hooks, which allow you to execute shell commands throughout a function's lifecycle. For example you may use these hooks to run tests or linting before a deploy, or to transpile source code using Babel, CoffeeScript, or similar.

Hooks may be specified in project.json or function.json. Hooks are executed in the function's directory, not the project's directory. If a hook exits > 0 then Apex will halt and display the error.

Internally Apex uses these hooks to implement out-of-the-box support for Golang and other compiled languages.

## Supported hooks

- `build` run before a function zip is built (use this to compile binaries or transform source)
- `deploy` run before a function is deployed (useful for testing, linting)
- `clean` run after a function is deployed (useful for cleaning up build artifacts)

## Examples

Here's the hooks used internally for Golang support.

```json
{
  "hooks": {
    "build": "GOOS=linux GOARCH=amd64 go build -o main main.go",
    "clean": "rm -f main"
  }
}
```
