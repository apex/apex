
# Functions

A function is the smallest unit in Apex. A function represents an AWS Lambda function.

Functions __MUST__ include at least one source file (runtime dependent), such as index.js or main.go. Optionally a function.json file may be placed in the _function's directory_, specifying details such as the memory allocated or the AWS IAM role. If one or more functions is missing a function.json file you __MUST__ provide defaults for the required fields in project.json (see "Projects" for an example).

# Configuration

```json
{
  "description": "Node.js example function",
  "runtime": "nodejs",
  "memory": 128,
  "timeout": 5,
  "role": "arn:aws:iam::293503197324:role/lambda"
}
```

### Fields

Fields marked as `inherited` may be provided in the [[project.json]] file instead.

## description (string, optional)

Description of the function. This is used as the description in AWS Lambda.

## runtime (string, required, inherited)

Runtime of the function. This is used as the runtime in AWS Lambda, or when required, is used to determine that the Node.js shim should be used. For example when this field is "golang", the canonical runtime used is "nodejs" and a shim is injected into the zip file.

## memory (int, required, inherited)

Memory allocated to the function, in megabytes.

## timeout (int, required, inherited)

Function timeout in sections.

## role (string, required, inherited)

AWS Lambda role used.
