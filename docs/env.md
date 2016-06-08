
AWS Lambda does not support environment variables out of the box, so this is a feature provided by Apex.

The `-s, --set` flag allows you to set environment variables which are exposed to the function at runtime. For example in Node.js using `process.env.NAME` or in Go using `os.Getenv("NAME")`. Behind the scenes this generates a `.env.json` file which is injected into your function's zip file upon deploy. You may use this flag multiple times.

## Examples

For example suppose you had a Loggly log collector and it needs an API token, you might deploy with:

```
$ apex deploy -s LOGGLY_TOKEN=token log-collector
```

Or suppose you have multiple functions using the GitHub API, you may want to expose a token to all of them:

```
$ apex deploy -s GITHUB_TOKEN=token
```

Specify environment variables in project.json or function.json:

```json
{
  "environment": {
    "LOGGLY_TOKEN": "12314212213123"
  }
}
```

Environment variables in json configuration are must be string.
