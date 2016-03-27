
The `apex metrics` command provides a quick glance at the overall metrics for your functions, displaying the number of invocations, total execution duration, throttling, and errors within a given time period.

## Examples

View all function metrics within the last 24 hours:

```sh
$ apex metrics

uppercase
  invocations: 1242
  duration: 65234ms
  throttles: 0
  error: 0

lowercase
  invocations: 1420
  duration: 65234ms
  throttles: 0
  error: 5

```

View metrics for the last 15 minutes:

```sh
$ apex metrics --since 15m

uppercase
  invocations: 16
  duration: 4212ms
  throttles: 0
  error: 0

lowercase
  invocations: 23
  duration: 5200ms
  throttles: 0
  error: 5

```
