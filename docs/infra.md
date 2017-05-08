
Apex is integrated with [Terraform](https://www.terraform.io/) to provide infrastructure management. Apex currently only manages Lambda functions, so you'll likely want to use Terraform or CloudFormation to manage additional resources such as Lambda sources.

## Managing infrastructure

The `apex infra` command is effectively a wrapper around the `terraform` command. Apex provides several variables and helps provide structure for multiple Terraform environments.

Each environment such as "prod" or "stage" lives in the ./infrastructure directory. For reference it may look something like:

```
infrastructure/
├── prod
│   └── main.tf
├── stage
│   └── main.tf
```

For example `apex infra --env prod plan` is effectively equivalent to the following command, with many `-var`'s passed to expose information from Apex.

```
$ cd infrastructure/prod && terraform plan
```

The environment is specified via the `--env` flag, or by default falls back on the `defaultEnvironment` property of project.json.

## Terraform variables

Currently the following variables are exposed to Terraform:

- `aws_region` the AWS region name such as "us-west-2"
- `apex_environment` the environment name such as "prod" or "stage"
- `apex_function_role` the Lambda role ARN
- `apex_function_arns` A map of all lambda functions
- `apex_function_names` A map of all the names of the lambda functions

## Notes

- You'll typically need to assign `${apex_function_myfunction}:current` to specify that the "current" alias is referenced.
- The `apex_function_NAME` variables are not available until the functions have been deployed (via `apex deploy`) at least once.
