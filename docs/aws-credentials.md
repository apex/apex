
Before using Apex you need to first give it your account credentials so that Apex can manage resources. There are a number of ways to do that, which are outlined here.

## Via environment variables

Using environment variables only, you must specify the following:

- `AWS_ACCESS_KEY_ID` AWS account access key
- `AWS_SECRET_ACCESS_KEY` AWS account secret key
- `AWS_REGION` AWS region

If you have multiple AWS projects you may want to consider using a tool such as [direnv](https://direnv.net/) to localize and automatically set the variables when
you're working on a project.

## Via ~/.aws files

Using the ~/.aws/credentials and ~/.aws/config files you may specify `AWS_PROFILE` to tell apex which one to reference. However, if you do not have a ~/.aws/config file, or "region" is not defined, you should set it with the `AWS_REGION` environment variable. To read more on configuring these files view [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

Here's an example of ~/.aws/credentials:

```
[example]
aws_access_key_id = xxxxxxxx
aws_secret_access_key = xxxxxxxxxxxxxxxxxxxxxxxx
```

Here's an example of ~/.aws/config:

```
[profile example]
output = json
region = us-west-2
```

## Via profile flag

If you have both ~/.aws/credentials and ~/.aws/config you may specify the profile directly with `apex --profile <name>` when issuing commands. This means you do not have to specify any environment variables, however you must provide it with each operation:

```
$ apex --profile myapp-prod deploy
```

## Via project configuration

You may store the profile name in the project.json file itself as shown in the following snippet. This is ideal since it ensures that you do not accidentally have a different environment set.

```json
{
  "profile": "myapp-prod"
}
```

## Via IAM Role

Using an IAM role can be achieved in two ways, via the __AWS_ROLE__ environment variable or via a command line flag `--iamrole`. As with other Apex credential loading, the command line flag will supersede the environment variable.

## Precedence

Precedence for loading the AWS credentials is:

- profile from flag
- profile from JSON config
- profile from env variables
- profile named "default"

## Minimum IAM Policy

Below is a policy for AWS [Identity and Access Management](https://aws.amazon.com/iam/) which provides the minimum privileges needed to use Apex to manage your Lambda functions.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iam:CreateRole",
        "iam:CreatePolicy",
        "iam:AttachRolePolicy",
        "iam:PassRole",
        "lambda:GetFunction",
        "lambda:ListFunctions",
        "lambda:CreateFunction",
        "lambda:DeleteFunction",
        "lambda:InvokeFunction",
        "lambda:GetFunctionConfiguration",
        "lambda:UpdateFunctionConfiguration",
        "lambda:UpdateFunctionCode",
        "lambda:CreateAlias",
        "lambda:UpdateAlias",
        "lambda:GetAlias",
        "lambda:ListAliases",
        "lambda:ListVersionsByFunction",
        "logs:FilterLogEvents",
        "cloudwatch:GetMetricStatistics"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

### Additional minimum IAM Policy to set VPC for Lambda

The following additional policies are needed to set VPC for your Lambda functions.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSubnets",
        "ec2:DescribeVpcs"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

Also, the role which apex made during `apex init` should have the `AWSLambdaVPCAccessExecutionRole` policy, see details in [an AWS document](https://docs.aws.amazon.com/lambda/latest/dg/vpc.html).
