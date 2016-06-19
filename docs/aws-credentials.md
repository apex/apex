
Before using Apex you need to first give it your account credentials so that Apex can manage resources. There are a number of ways to do that, which are outlined here.

## Via environment variables

Using environment variables only, you must specify the following:

- `AWS_ACCESS_KEY` AWS account access key
- `AWS_SECRET_KEY` AWS account secret key
- `AWS_REGION` AWS region

## Via ~/.aws files

Using the ~/.aws/credentials and ~/.aws/config files you may specify `AWS_PROFILE` to tell apex which one to reference. However, if you do not have a ~/.aws/config file, or "region" is not defined, you should set it with the `AWS_REGION` environment variable. To read more on configuring these files view [Configuring the AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

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

## Precedence

Precedence for loading the AWS credentials is:

- profile from flag
- profile from JSON config
- profile from env variables
- profile named "default"
