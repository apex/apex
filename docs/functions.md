
A function is the smallest unit in Apex. A function represents an AWS Lambda function.

Functions must include at least one source file (runtime dependent), such as index.js or main.go. Optionally a function.json file may be placed in the _function's directory_, specifying details such as the memory allocated or the AWS IAM role. If one or more functions is missing a function.json file you must provide defaults for the required fields in project.json (see "Projects" for an example).

## Configuration

```json
{
  "description": "Node.js example function",
  "runtime": "nodejs",
  "memory": 128,
  "timeout": 5,
  "role": "arn:aws:iam::293503197324:role/lambda"
}
```

## Fields

Fields marked as `inherited` may be provided in the project.json file instead.

### description

Description of the function. This is used as the description in AWS Lambda.

- type: `string`
- required

### runtime

Runtime of the function. This is used as the runtime in AWS Lambda, or when required, is used to determine that the Node.js shim should be used. For example when this field is "golang", the canonical runtime used is "nodejs" and a shim is injected into the zip file.

- type: `string`
- required
- inherited

### handler

Event handler name, this is the function invoked for a given event. Defaults are:

- nodejs: `index.handle` (index.js file with `handle` exported function)
- python: `handle`
- java: `lambda.Main::handler`

### memory

Memory allocated to the function, in megabytes.

- type: `int`
- required
- inherited

### timeout

Function timeout in seconds. Note that Lambda currently restricts durations to 5 minutes (300s).

- type: `int`
- required
- inherited

### role

AWS Lambda role used.

- type: `string`
- required
- inherited

### environment

Environment variables.

- type: `object`
- inherited

### retainedVersions

Number of retained function's versions on AWS Lambda. If not specified `deploy` command will leave 10 versions.

- type: `int`
- inherited

### vpc

If your function needs to access resources in a VPC security groups and subnets have to be provided. You must provide at least one security group and one subnet.

- type: `object`
- inherited

#### vpc.securityGroups

List of security groups IDs

- type: `array`
- inherited

#### vpc.subnets

List of subnets IDs

- type: `array`
- inherited

### kms_arn

Optional ARN of the KMS key used to encrypt your function's environment variables. If empty, it means you are using the AWS Lambda default service key.

### deadletter_arn

Optional ARN of an Amazon SQS queue or Amazon SNS topic you specify as your Dead Letter Queue (DLQ).

### region

If your function needs to be deployed different region from the whole project.

- type: `string`

### edge

If your function is for Lambda@Edge. The Edge function needs to be deployed on N. Virginia region, excluding all environment variables.

- type: `boolean`
