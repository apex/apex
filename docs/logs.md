
# View function logs

Apex is integrated with CloudWatch Logs to view the output of functions. By default the logs for all functions will be displayed, unless one or more function names are passed to `apex logs`. You may also specify the duration of time in which the history is displayed (defaults to 5 minutes), as well as following and filtering results.

## Examples

View all function logs within the last 5 minutes:

```
$ apex logs
```

View logs for "uppercase" and "lowercase" functions:

```
$ apex logs uppercase lowercase
```

Follow or tail the log output for all functions:

```
$ apex logs -f
```

Follow a specific function:

```
$ apex logs -f foo
```

Follow filtered by pattern "error":

```
$ apex logs -f foo --filter error
$ apex logs -f foo -F error
```

Output the last hour of logs:

```
$ apex logs --duration 1h
$ apex logs -d 1h
```
