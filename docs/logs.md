
Apex is integrated with CloudWatch Logs to view the output of functions. By default the logs for all functions will be displayed, unless one or more function names are passed to `apex logs`. You may also specify the duration of time in which the history is displayed (defaults to 5 minutes), as well as following and filtering results.

## Examples

View all function logs within the last 5 minutes:

```sh
$ apex logs
```

View logs for "uppercase" and "lowercase" functions:

```sh
$ apex logs uppercase lowercase
```

Follow or tail the log output for all functions:

```sh
$ apex logs -f
```

Follow a specific function:

```sh
$ apex logs -f foo
```

Follow filtered by pattern "error":

```sh
$ apex logs -f foo --filter error
$ apex logs -f foo -F error
```

Output the last hour of logs:

```sh
$ apex logs --since 1h
$ apex logs -s 1h
```

Log all functions which name starts with "auth":

```sh
$ apex logs auth*
```
