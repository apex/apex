
Apex can help initialize a basic project to get you started. First specify your AWS credentials as mentioned in the previous section, then run `apex init`:

```
$ export AWS_PROFILE=sloths-stage
$ apex init
```

You'll be presented with a few prompts, the project's default Lambda IAM role & policy will be created, then you're ready to go!


```
             _    ____  _______  __
            / \  |  _ \| ____\ \/ /
           / _ \ | |_) |  _|  \  /
          / ___ \|  __/| |___ /  \
         /_/   \_\_|   |_____/_/\_\



  Enter the name of your project. It should be machine-friendly, as this
  is used to prefix your functions in Lambda.

    Project name: sloths

  Enter an optional description of your project.

    Project description: My slothy project

  [+] creating IAM sloth_lambda_function role
  [+] creating IAM sloth_lambda_logs policy
  [+] attaching policy to lambda_function role.
  [+] creating ./project.json
  [+] creating ./functions

  Setup complete, deploy those functions!

    $ apex deploy

```

Now try invoking the sample function:

```
$ apex invoke hello
{"hello":"world"}
```
