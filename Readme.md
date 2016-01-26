
![Apex Serverless Architecture](assets/logo.png)

Apex is a small tool for deploying and managing [AWS Lambda](https://aws.amazon.com/lambda/) functions. With shims for languages not yet supported by Lambda, you can use Golang out of the box.

## Installation

Download [binaries](https://github.com/apex/apex/releases):

```
latest=$(curl -s https://api.github.com/repos/apex/apex/tags | grep name | head -n 1 | sed 's/[," ]//g' | cut -d ':' -f 2)
curl -sL https://github.com/apex/apex/releases/download/$latest/apex_darwin_amd64 -o /usr/local/bin/apex
chmod +x $_
```

Or from master:

```
go get github.com/apex/apex/cmd/apex
```

Or upgrading:

```
apex upgrade
```

## Runtimes

Currently supports:

- Nodejs
- Golang
- Python

## Features

- Supports languages Lambda does not natively support via shim, such as Go
- Binary install (useful for continuous deployment in CI etc)
- Hook support for running commands (transpile code, lint, etc)
- Project level function and resource management
- Configuration inheritance and overrides
- Command-line function invocation with JSON streams
- Transparently generates a zip for your deploy
- Ignore deploying files with .apexignore
- Function rollback support
- Tail function CloudWatchLogs
- Concurrency for quick deploys
- Dry-run to preview changes

## Example

Apex projects are made up of a project.json configuration file, and zero or more Lambda functions defined in the "functions" directory. Here's an example file structure:

```
project.json
functions
├── bar
│   ├── function.json
│   └── index.js
├── baz
│   ├── function.json
│   └── index.js
└── foo
    ├── function.json
    └── index.js
```

The project.json file defines project level configuration that applies to all functions, and defines dependencies. For this simple example the following will do:

```json
{
  "name": "example",
  "description": "Example project"
}
```

Each function uses a function.json configuration file to define function-specific properties such as the runtime, amount of memory allocated, and timeout. This file is completely optional, as you can specify defaults in your project.json file. For example:

```json
{
  "name": "bar",
  "description": "Node.js example function",
  "runtime": "nodejs",
  "memory": 128,
  "timeout": 5,
  "role": "arn:aws:iam::293503197324:role/lambda"
}
```

Now the directory structure for your project would be:

```
project.json
functions
├── bar
│   └── index.js
├── baz
│   └── index.js
└── foo
    └── index.js
```

Finally the source for the functions themselves look like this in Node.js:

```js
console.log('start bar')
exports.handle = function(e, ctx) {
  ctx.succeed({ hello: 'bar' })
}
```

Or using the Golang Lambda package, Apex supports Golang out of the box with a Node.js shim:

```go
package main

import (
  "encoding/json"

  "github.com/apex/apex"
)

type Message struct {
  Hello string `json:"hello"`
}

func main() {
  apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
    return &Message{"baz"}, nil
  })
}
```

Apex operates at the project level, but many commands allow you to specify specific functions. For example you may deploy the entire project with a single command:

```
$ apex deploy
```

Or whitelist functions to deploy:

```
$ apex deploy foo bar
```

Invoke it!

```
$ echo '{ "some": "data" }' | apex invoke foo
{ "hello": "foo" }
```

See the [Wiki](https://github.com/apex/apex/wiki) for more information.

## Credentials

Via environment variables:

- `AWS_ACCESS_KEY` AWS account access key
- `AWS_SECRET_KEY` AWS account secret key
- `AWS_REGION` AWS region

Via ~/.aws configuration:

- `AWS_PROFILE` profile name to use
- `AWS_REGION` AWS region (aws-sdk-go does not read ~/.aws/config)

## Links

- [Wiki](https://github.com/apex/apex/wiki)
- [Project Examples](_examples) with source

## Contributors

- [TJ Holowaychuk](https://github.com/tj)
- [Maciej Winnicki](https://github.com/mthenw)
- [Pilwon Huh](https://github.com/pilwon)
- [Faraz Fazli](https://github.com/farazfazli)
- [Johannes Boyne](https://github.com/johannesboyne)

## Badges

[![Build Status](https://semaphoreci.com/api/v1/projects/d27ff350-b9c5-4d99-96e5-64b1afb441c5/649392/badge.svg)](https://semaphoreci.com/tj/apex)
[![Slack Status](https://apex-dev.azurewebsites.net/badge.svg)](https://apex-dev.azurewebsites.net/)
[![GoDoc](https://godoc.org/github.com/apex/apex?status.svg)](https://godoc.org/github.com/apex/apex)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-experimental-orange.svg)

---

> [tjholowaychuk.com](http://tjholowaychuk.com) &nbsp;&middot;&nbsp;
> GitHub [@tj](https://github.com/tj) &nbsp;&middot;&nbsp;
> Twitter [@tjholowaychuk](https://twitter.com/tjholowaychuk)

