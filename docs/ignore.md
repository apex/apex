
## Ignoring files with .apexignore

Optional `.apexignore` files may be placed in the project's root, or within specific function directories. It uses [shell pattern matching](https://www.gnu.org/software/findutils/manual/html_node/find_html/Shell-Pattern-Matching.html); all patterns defined, are relative to function directories and not the project itself.

## Example

Here's an example ignoring go source files and the function.json itself:

```
*.go
function.json
```
