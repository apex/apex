
Apex now supports the new native AWS Lambda environment variables, which are encrypted via KMS. There are several ways to set these variables, let's take a look!

## Flag --set

The `-s, --set` flag allows you to set environment variables which are exposed to the function at runtime. For example in Node.js using `process.env.NAME` or in Go using `os.Getenv("NAME")`. You may use this flag multiple times.

For example suppose you had a Loggly log collector and it needs an API token, you might deploy with:

```
$ apex deploy -s LOGGLY_TOKEN=token log-collector
```

## Flag --env-file

The `-E, --env-file` flag allows you to set multiple environment variables using a JSON file.

```
$ apex deploy --env-file /path/to/env.json
```

Sample env.json:

```json
{
  "LOGGLY_TOKEN": "12314212213123"
}
```

## Config (project.json or function.json)

Specify environment variables in project.json or function.json, note that the values _must_ be strings.

```json
{
  "environment": {
    "LOGGLY_TOKEN": "12314212213123"
  }
}
```

## Precedence

The precedence is currently as follows:

- `-s, --set` flag values
- `-E, --env-file` file values
- environment variables specified in project.json or function.json
