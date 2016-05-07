
Apex allows you to delete functions, however it will prompt by default. Use the `-f, --force` flag to override this behaviour. You may pass zero or more function names.

## Examples

Delete all with prompt:

```sh
$ apex delete
The following will be deleted:

  - bar
  - foo

Are you sure? (yes/no):
```

For delete of all functions:

```sh
$ apex delete -f
```

For delete of specific functions:

```sh
$ apex delete -f foo bar
```

Delete all functions which name starts with "auth":

```sh
$ apex delete auth*
```
