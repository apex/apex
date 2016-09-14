
Apex allows you to invoke functions from the command-line, optionally passing a JSON event or stream to STDIN. It's important to note that `invoke` will execute the remote Lambda function and not locally execute your function. It will execute the $LATEST Lambda function available.

## Examples

Invoke without an event:

```sh
$ apex invoke collect-stats
```

Invoke with JSON event:

```sh
$ echo -n '{ "value": "Tobi the ferret" }' | apex invoke uppercase
{ "value": "TOBI THE FERRET" }
```

Invoke from a file:

```sh
$ apex invoke uppercase < event.json
```

Invoke a with stdin clipboard data:

```sh
$ pbpaste | apex invoke auth
```

Invoke function in a different project:

```sh
$ pbpaste | apex -C path/to/project invoke auth
```

Streaming invokes making multiple requests, generating data with [phony][1]:

```sh
$ echo -n '{ "user": "{{name}}" }' | phony | apex invoke uppercase
{"user":"DELMER MALONE"}
{"user":"JC REEVES"}
{"user":"LUNA FLETCHER"}
...
```

Streaming invokes making multiple requests and outputting the response logs:

```sh
$ echo -n '{ "user": "{{name}}" }' | phony | apex invoke uppercase --logs
START RequestId: 30e826a4-a6b5-11e5-9257-c1543e9b73ac Version: $LATEST
END RequestId: 30e826a4-a6b5-11e5-9257-c1543e9b73ac
REPORT RequestId: 30e826a4-a6b5-11e5-9257-c1543e9b73ac	Duration: 0.73 ms	Billed Duration: 100 ms 	Memory Size: 128 MB	Max Memory Used: 10 MB
{"user":"COLTON RHODES"}
START RequestId: 30f0b23c-a6b5-11e5-a034-ad63d48ca53a Version: $LATEST
END RequestId: 30f0b23c-a6b5-11e5-a034-ad63d48ca53a
REPORT RequestId: 30f0b23c-a6b5-11e5-a034-ad63d48ca53a	Duration: 2.56 ms	Billed Duration: 100 ms 	Memory Size: 128 MB	Max Memory Used: 9 MB
{"user":"CAROLA BECK"}
START RequestId: 30f51e67-a6b5-11e5-8929-f53378ef0f47 Version: $LATEST
END RequestId: 30f51e67-a6b5-11e5-8929-f53378ef0f47
REPORT RequestId: 30f51e67-a6b5-11e5-8929-f53378ef0f47	Duration: 0.22 ms	Billed Duration: 100 ms 	Memory Size: 128 MB	Max Memory Used: 9 MB
{"user":"TOBI FERRET"}
...
```

[1]: https://github.com/yields/phony
