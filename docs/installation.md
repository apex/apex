
On macOS, Linux, or OpenBSD run the following:

```
curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh
```

This command will install `apex` binary as `/usr/local/bin/apex` and
you may need to run the `sudo` version below, or alternatively chown `/usr/local`:
```
curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sudo sh
```

You can specify a destination as below

```
curl https://raw.githubusercontent.com/apex/apex/master/install.sh | DEST=$HOME/bin/apex sh
```

this command will install `apex` binary as `$HOME/bin/apex` and you may not need to use `sudo`.

On Windows download [binary](https://github.com/apex/apex/releases).

If already installed, upgrade with:

```
apex upgrade
```
