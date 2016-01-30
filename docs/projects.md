
# Projects

A "project" is the largest unit of abstraction in Apex. A project consists of collection of AWS Lambda functions, and
all `apex(1)` operations have access to these functions.

## Configuration

Projects __MUST__ have a project.json file in the root directory. This file contains details about your project, as well as
defaults for functions, if desired. Here's an example of a project.json file declaring a default AWS IAM "role" and "memory" for _all_ functions.

```json
{
  "name": "node",
  "description": "Node.js example project",
  "role": "arn:aws:iam::293503197324:role/lambda",
  "memory": 512
}
```

## Fields

### name (string, required)

Name of the project. This field is used in the default value for "nameTemplate" to prevent collisions between multiple projects.

### description (string, optional)

Description of the project. This field is informational.

### runtime (string, optional)

Default runtime of function(s) unless specified in their [[function.json]] configuration.

### memory (int, optional)

Default memory allocation of function(s) unless specified in their [[function.json]] configuration.

### timeout (int, optional)

Default timeout of function(s) unless specified in their [[function.json]] configuration.

### role (string, optional)

Default role of function(s) unless specified in their [[function.json]] configuration.

### nameTemplate (string, optional)

Template used to compute the function names. By default the template `{{.Project.Name}}_{{.Function.Name}}` is used, for example project "api" and `./functions/users` becomes "api_users". To disable prefixing, use `{{.Function.Name}}`, which would result in "users".
