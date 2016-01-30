
Apex generates a zip file for you upon deploy, however sometimes it can be useful to see exactly what's included in this file for debugging purposes. The `apex build` command outputs the zip to STDOUT for this purpose.

## Examples

Output zip to out.zip:

```sh
$ apex build foo > out.zip
```
