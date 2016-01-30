
# Rolling back function deploys

Apex allows you to roll back to the previous, or specified version of a function.

## Examples

Rollback to the previous release:

```
$ apex rollback foo
```

Rollback to specific version:

```
$ apex rollback foo 5
```

Preview rollback with `--dry-run`:

```
$ apex rollback --dry-run lowercase

~ alias testing_lowercase
 alias: current
 version: 2

$ apex rollback --dry-run uppercase 1

~ alias testing_uppercase
 version: 1
 alias: current
```
