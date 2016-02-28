
![Apex Serverless Architecture](assets/logo.png)

Apex lets you build, deploy, and manage [AWS Lambda](https://aws.amazon.com/lambda/) functions with ease. With Apex you can use languages that are not natively supported by AWS Lambda, such as Golang, through the use of a Node.js shim injected into the build. A variety of workflow related tooling is provided for testing functions, rolling back deploys, viewing metrics, tailing logs, hooking into the build system and more.

## Installation

On OS X or Linux:

```
curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh
```

On Windows download [binary](https://github.com/apex/apex/releases).

If already installed, upgrade with:

```
apex upgrade
```

## Runtimes

Currently supports:

- Node.js
- Golang
- Python
- Java

Example projects for all supported runtimes can be found in [_examples](_examples) directory.

## Features

- Supports languages Lambda does not natively support via shim, such as Go
- Binary install (install apex quickly for continuous deployment in CI etc)
- Hook support for running commands (transpile code, lint, etc)
- Batteries included but optional (opt-in to higher level abstractions)
- Project level function and resource management
- Configuration inheritance and overrides
- Command-line function invocation with JSON streams
- Transparently generates a zip for your deploy
- Ignore deploying files with .apexignore
- Function rollback support
- Tail function logs
- Concurrency for quick deploys
- Dry-run to preview changes
- VPC support

## Example

Apex projects are made up of a project.json configuration file, and zero or more Lambda functions defined in the "functions" directory. Here's an example file structure:

```
project.json
functions
├── bar
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
└── foo
    └── index.js
```

Finally the source for the functions themselves look like this in Node.js:

```js
console.log('start bar')
exports.handle = function(e, ctx) {
  ctx.succeed({ hello: e.name })
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
$ echo '{ "name": "Tobi" }' | apex invoke bar
{ "hello": "Tobi" }
```

See the [Documentation](docs) for more information.

## Links

- [Website](http://apex.run)
- [Twitter](https://twitter.com/apexserverless)

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
