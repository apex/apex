
# Deleting functions

Apex allows you to delete functions, however it will prompt by default. Use the `-f, --force` flag to override this behaviour. You may pass zero or more function names.

## Examples

Delete all with prompt:

```
$ apex delete
The following will be deleted:

  - bar
  - foo

Are you sure? (yes/no):
```

For delete of all functions:

```
$ apex delete -f
```

For delete of specific functions:

```
$ apex delete -f foo bar
```
