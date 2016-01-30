
# Installation

Download [binaries](https://github.com/apex/apex/releases) for your platform, or if you're on OSX grab the latest release with:

```
latest=$(curl -s https://api.github.com/repos/apex/apex/tags | grep name | head -n 1 | sed 's/[," ]//g' | cut -d ':' -f 2)
curl -sL https://github.com/apex/apex/releases/download/$latest/apex_darwin_amd64 -o /usr/local/bin/apex
chmod +x $_
```

Install from master with `go-get`:

```
go get github.com/apex/apex/cmd/apex
```

Upgrade to the latest release:

```
apex upgrade
```
