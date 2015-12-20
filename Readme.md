
# Apex

Apex is a small tool for deploying and managing [AWS Lambda](https://aws.amazon.com/lambda/) functions. With shims for languages not yet supported by Lambda, you can use Golang out of the box.

## Installation

Download [binaries](https://github.com/apex/apex/releases) or:

```
$ go get github.com/apex/apex/...
```

## Runtimes

Currently supports:

- Nodejs
- Golang
- Python

## Features

- Supports languages Lambda does not natively support via shim
- Binary install (useful for continuous deployment in CI etc)
- Command-line function invocation with JSON streams
- Transparently generates a zip for your deploy

## Example

This example shows how you can use Apex to launch a simple Node.js echo function.

First create the function implementation in "index.js":

```js
exports.handle = function(e, ctx) {
  ctx.succeed(e)
}
```

Next create a "lambda.json" with the function name a configuration:

```json
{
  "name": "echo",
  "description": "Echo request example",
  "runtime": "nodejs",
  "memory": 128,
  "timeout": 5
}
```

Deploy the function:

```
$ apex deploy
```

Create a file with a sample request in "request.js":

```js
{
  "event": {
    "hello": "world"
  },
  "context": {
    "user": "Tobi"
  }
}
```

Test out your new function:

```
$ apex invoke < request.json
{"hello":"world"}
```

## Streaming input

The `invoke` sub-command allows you to stream input over stdin:

```
$ apex invoke < request.json
```

This not only works for single requests, but for multiple, as shown in the following example using [phony(1)](https://github.com/yields/phony):

```
$ echo -n '{ "event": { "user": "{{name}}" } }' | phony | apex invoke
{"user":"Delmer Malone"}
{"user":"Jc Reeves"}
{"user":"Luna Fletcher"}
...
```

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

# License

MIT
