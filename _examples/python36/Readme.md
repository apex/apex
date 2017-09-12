
To run the example first setup your [AWS Credentials](http://apex.run/#aws-credentials), and ensure "role" in ./project.json is set to your Lambda function ARN.

Deploy the functions:

```
$ apex deploy
```

Try it out:

```
$ apex invoke simple < event.json
```
