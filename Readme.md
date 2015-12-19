
# Apex

Apex is a small tool for deploying and managing [AWS Lambda](https://aws.amazon.com/lambda/) functions. With shims for languages not yet supported by Lambda, you can use Golang, Ruby, and others out of the box.

## Installation

```
$ go get github.com/apex/apex/cmd/apex
```

## Example

This example shows how you can use Apex to launch a simple Node.js echo function.

First create the function implementation in "index.js".

```js
exports.handle = function(e, ctx) {
  ctx.succeed(e)
}
```

Next create a "package.json" with the function name a configuration:

```json
{
  "name": "echo",
  "description": "Echo request example",
  "runtime": "nodejs",
  "memory": 128,
  "timeout": 5
}
```

Create and deploy the function:

```
$ apex deploy
```

Create a file with a sample request "request.js":

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

## Environment variables

Currently the following environment variables are used:

- `AWS_ACCESS_KEY` AWS account access key
- `AWS_SECRET_KEY` AWS account secret key
- `AWS_REGION` AWS region

## Shimming

Apex uses a Node.js shim for non-native language support. Because of this you must use stderr for logging, as stdout is used for the reply JSON.

## FAQ

### How do you manage multiple environments?

It's highly recommended to create separate AWS accounts for staging and production environments. This provides complete isolation of resources, and allows you to easily provide granular access to environments as required.

### How do you structure projects with multiple Lambda functions?

Since apex(1) is function-oriented, you can structure projects however you like. You may use a single repository with multiple functions, or a repository per function. Use the `-C` flag to invoke commands against a specific function's directory.

### How is this different than serverless.com?

The Serverless implementation is written in Node.js, and has a larger scope than Apex. If you're looking for full project management Serverless may be a better option for you!

Serverless uses [CloudFormation](https://aws.amazon.com/cloudformation/) to bootstrap resources, which can be great for getting started, but less robust than [Terraform](https://terraform.io/) for managing infrastructure throughout its lifetime. For this reason Apex does not currently provide resource management.

At the time of writting Serverless does not support shimming for languages which are not supported natively by Lambda, such as Golang. Apex does this for you out of the box.

# License

MIT