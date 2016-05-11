
A "project" is the largest unit of abstraction in Apex. A project consists of collection of AWS Lambda functions, and
all `apex(1)` operations have access to these functions.

## Configuration

Projects must have a project.json file in the root directory. This file contains details about your project, as well as
defaults for functions, if desired. Here's an example of a project.json file declaring a default AWS IAM "role" and "memory" for all functions.

```json
{
  "name": "node",
  "description": "Node.js example project",
  "role": "arn:aws:iam::293503197324:role/lambda",
  "memory": 512
}
```

## Symlinks

It's important to note that Apex supports symlinked files and directories. Apex will read the links and pull in these files, even if the links aren't to files within your function. This enables the use of `npm link`, shared configuration and so on.

## Fields

### name

Name of the project. This field is used in the default value for "nameTemplate" to prevent collisions between multiple projects.

- type: `string`
- required

### description

Description of the project. This field is informational.

- type: `string`

### runtime

Default runtime of function(s) unless specified in their function.json configuration.

- type: `string`

### memory

Default memory allocation of function(s) unless specified in their function.json configuration.

- type: `int`

### timeout

Default timeout of function(s) unless specified in their function.json configuration.

- type: `int`

### role

Default role of function(s) unless specified in their function.json configuration.

- type: `string`

### defaultEnvironment

Default infrastructure environment.

- type: `string`

### environment

Default environment variables of function(s) unless specified in their function.json configuration.

- type: `object`

### nameTemplate

Template used to compute the function names. By default the template `{{.Project.Name}}_{{.Function.Name}}` is used, for example project "api" and `./functions/users` becomes "api_users". To disable prefixing, use `{{.Function.Name}}`, which would result in "users".

- type: `string`

### retainedVersions

Default number of retained function's versions on AWS Lambda unless specified in their function.json configuration.

- type: `int`

### vpc

Default VPC configuration of function(s) unless specified in their function.json configuration.

- type: `object`

#### vpc.securityGroups

List of security groups IDs

- type: `array`

#### vpc.subnets

List of subnets IDs

- type: `array`
