
# Getting started

Step by step instructions for getting started with Apex and AWS services.

## AWS setup

Before you get started, you will need to setup some things on the AWS console.

1. Sign into the [AWS IAM Console](https://console.aws.amazon.com/iam/home#home)
2. Assign a [user](https://console.aws.amazon.com/iam/home#users) with the role `AWSLambdaFullAccess`
3. Create a new [role](https://console.aws.amazon.com/iam/home#roles) for your Lambda functions.
Note: Lambda functions don't need any policy, but adding the policy `CloudWatchLogsFullAccess` will allow log tailing.
4. Click on the Role and copy the Role ARN.


## AWS credentials

You'll need your AWS credentials set up. You have two options, set the following environment variables:

- `AWS_ACCESS_KEY` AWS account access key
- `AWS_SECRET_KEY` AWS account secret key
- `AWS_REGION` AWS region

Or set up ~/.aws/credentials via `aws configure`, and use the following environment variables to choose which profile to use:

- `AWS_PROFILE` profile name to use
- `AWS_REGION` AWS region (aws-sdk-go does not read ~/.aws/config)

## Node.js

Create a new Apex project:

```sh
mkdir myproject && cd myproject
touch project.json
mkdir -p functions/uppercase
touch functions/uppercase/index.js
```

Add the following to project.json. Give the project a name which is used to prefix the functions in Lambda, a description, and function defaults for memory, timeout and AWS role. The AWS role must be created in AWS first, and the policy will vary depending on the needs of your function(s). As a minimum you'll probably want the role to have the `CloudWatchLogsFullAccess` policy for logging.

```json
{
  "name": "myproject",
  "description": "Node.js example project",
  "memory": 128,
  "timeout": 5,
  "role": "<ROLE ARN HERE>"
}
```

Create the first function in `./functions/uppercase/index.js`:

```js
console.log('starting function')
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({ value: e.value.toUpperCase() })
}
```

Deploy it:

```
$ apex deploy
   • deploying                 function=uppercase
   • created zip (275 B)       function=uppercase
   • creating function         function=uppercase
   • creating alias            function=uppercase
   • deployed                  function=uppercase name=myproject_uppercase version=1
   • deploying config          function=uppercase
```

## Go

Create a new Apex project:

```sh
mkdir myproject && cd myproject
touch project.json
mkdir -p functions/uppercase
touch functions/uppercase/main.go
```

Add the following to project.json. Give the project a name which is used to prefix the functions in Lambda, a description, and function defaults for memory, timeout and AWS role. The AWS role must be created in AWS first, and the policy will vary depending on the needs of your function(s). As a minimum you'll probably want the role to have the `CloudWatchLogsFullAccess` policy for logging.

```json
{
  "name": "myproject",
  "description": "Go example project",
  "memory": 128,
  "timeout": 5,
  "role": "<ROLE ARN HERE>"
}
```

Create the first function in `./functions/uppercase/main.go`:
```go
package main

import (
  "encoding/json"
  "strings"

  "github.com/apex/go-apex"
)

type Message struct {
  Value string `json:"value"`
}

func main() {
  apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
    var msg Message

    if err := json.Unmarshal(event, &msg); err != nil {
      return nil, err
    }

    return &Message{strings.ToUpper(msg.Value)}, nil
  })
}

```


Deploy it:

```
$ apex deploy uppercase                                 
   • deploying                 function=uppercase
   • created zip (3.1 MB)      function=uppercase
   • creating function         function=uppercase
   • creating alias            function=uppercase
   • deploying config          function=uppercase
```

## Invoking functions

To invoke a function you'll first need some input – create ./request.json and add a dummy event:

```json
{
  "event": {
    "value": "Tobi the Ferret"
  }
}
```

All invocation is done over stdin, so pass the request.json to `apex invoke` via `<`:

```
$ apex invoke uppercase < request.json
{"value":"TOBI THE FERRET"}
```

Invoke with the `-L` or `--logs` flag to view log output along with the response:

```
$ apex invoke --logs uppercase < request.json
START RequestId: 6ddd4016-bfcf-11e5-b01b-e76f932df70e Version: 1
2016-01-20T23:42:12.726Z  6ddd4016-bfcf-11e5-b01b-e76f932df70e  processing event: {}
END RequestId: 6ddd4016-bfcf-11e5-b01b-e76f932df70e
REPORT RequestId: 6ddd4016-bfcf-11e5-b01b-e76f932df70e  Duration: 2.31 ms Billed Duration: 100 ms   Memory Size: 128 MB Max Memory Used: 27 MB  
{"value":"TOBI THE FERRET"}
```

## Deleting functions

Clean up and remove the functions created:

```
$ apex delete --force
   • deleting                  function=uppercase
```
