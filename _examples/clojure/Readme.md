
To run the example first setup your [AWS Credentials](http://apex.run/#aws-credentials), and ensure "role" in ./project.json is set to your Lambda function ARN.

Install Lein:

```
$ wget -P ~/bin https://raw.githubusercontent.com/technomancy/leiningen/stable/bin/lein
$ chmod 755 ~/bin/lein
$ ~/bin/lein
```

Deploy the functions:

```
$ apex deploy
```

Try it out:

```
$ echo '{ "say": "Hello World!" }' | apex invoke say
```
