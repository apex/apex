
Apex lets you perform a "dry run" of any operation with the `--dry-run` flag, and no destructive AWS changes will be made.

## Notation

Dry runs use the following symbols:

- `+` resource will be created
- `-` resource will be removed
- `~` resource will be updated

## Examples

For example if you have the functions "foo" and "bar" which have never been deployed, you'll see the following output. This output represents the final requests made to AWS; notice how the function names are prefixed with the project's ("testing") to prevent collisions, and aliases are made to maintain the "current" release alias.

```sh
$ apex deploy --dry-run

  + function testing_bar
    handler: _apex_index.handle
    runtime: nodejs
    memory: 128
    timeout: 5

  + alias testing_bar
    alias: current
    version: 1

  + function testing_foo
    memory: 128
    timeout: 5
    handler: _apex_index.handle
    runtime: nodejs

  + alias testing_foo
    alias: current
    version: 1
```

If you were to run `apex deploy foo`, then run `apex deploy --dry-run` again, you'll see that only "bar" needs deploying:

```sh
$ apex deploy --dry-run

  + function testing_bar
    runtime: nodejs
    memory: 128
    timeout: 5
    handler: index.handle

  + alias testing_bar
    alias: current
    version: 1
```

Similarly this can be used to preview configuration changes:

```sh
$ apex deploy --dry-run

  ~ alias testing_foo
    alias: current
    version: $LATEST

  ~ config testing_foo
    memory: 128 -> 512
    timeout: 5 -> 10
```

As mentioned this works for all AWS operations, here's a delete preview:

```sh
$ apex delete --dry-run -f

  - function testing_bar

  - function testing_foo
```
