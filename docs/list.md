
Apex supports listing of functions in various outputs, currently human-friendly terminal output and "tfvars" support for integration with Terraform.

## Examples

List all functions and their configuration:

```sh
$ apex list

  bar
    runtime: nodejs
    memory: 128mb
    timeout: 5s
    role: arn:aws:iam::293503197324:role/lambda
    handler: index.handle
    aliases: current@v3, foo@v4

  foo
    runtime: nodejs
    memory: 512mb
    timeout: 10s
    role: arn:aws:iam::293503197324:role/lambda
    handler: index.handle
    aliases: current@v12
```

Terraform vars output:

```sh
$ apex list --tfvars
apex_function_bar="arn:aws:lambda:us-west-2:293503197324:function:testing_bar"
apex_function_foo="arn:aws:lambda:us-west-2:293503197324:function:testing_foo"
```
