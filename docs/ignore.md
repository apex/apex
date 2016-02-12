
Optional `.apexignore` files may be placed in the project's root, or within specific function directories. It uses [.gitignore pattern format](https://git-scm.com/docs/gitignore#_pattern_format); all patterns defined, are relative to function directories and not the project itself. By default both `.apexignore` and `function.json` are ignored.

## Example

Here's an example ignoring go source files and the function.json itself:

```
*.go
```
