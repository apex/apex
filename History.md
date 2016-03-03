
v0.6.1 / 2016-02-08
===================

  * fix adding generated files to archive. Closes #221

v0.6.0 / 2016-02-06
===================

  * add babel example
  * add hook documentation. Closes #212
  * add FileInfo wrapper to zero mtime. Closes #152
  * add inline markdown docs. Closes #194
  * add env populating for python runtime. Closes #202
  * add multiple functions to metrics. Closes #182
  * add support for AWS_PROFILE region from ~/.aws/config
  * add function.json ignoring by default. Closes #196
  * add support for pulling region from aws-cli config. Closes #90
  * add ./docs moved from wiki and additional content
  * add support for invoke without event. Closes #173
  * add pager to docs command when stdout is a tty
  * add function name to error output. Closes #160
  * add multi-function log support. Closes #159
  * add tfvars output to list command. Closes #155
  * add nodejs prelude script injection. Closes #140
  * add support for aws profile switching via aws shared credentials
  * add support for overriding the build hook when using golang runtime
  * rename --duration to --start. Closes #204
  * remove old --verbose flag
  * remove Go runtime, use apex/go-apex, updates examples. Closes #156
  * refactor shim plugin to access zip directly. Closes #190
  * refactor commands into packages. Closes #107
  * refactor functions loading. Closes #132
  * fix updating configuration. Closes #206
  * fix list command when function does not exist (ignore remote config)

v0.5.1 / 2016-01-25
===================

  * fix env variable precedence when set via flag
  * fix open file limit bug

v0.5.0 / 2016-01-24
===================

  * add metrics command
  * add hook support. Closes #68
  * add plugins, replacing runtimes. Closes #130
  * add coffeescript example using hooks
  * add invoke support to omit .event and .context. Closes #13
  * add .apexignore support. Closes #69
  * add simple CONTRIBUTING.md. Closes #121
  * add deploy -c, --concurrency. Closes #46
  * add "apex version" and dropped "-v" and "--version" global flag
  * add -f flag support to "apex logs". Do not follow by default
  * add wiki multi-arg support. Closes #117
  * change wiki code to bold instead of gray
  * remove deferring of file Close() for builds (keep fd count low)
  * rename {Project,Function}.SetEnv to Setenv
  * rename wiki command to docs

v0.4.1 / 2016-01-17
===================

  * add Project.name(fn) to compute nameTemplate
  * remove Function Config.Name support, fixing name reference bug. Closes #81

v0.4.0 / 2016-01-17
===================

  * add help command, pulling data from wiki. Closes #74
  * add project nameTemplate support. Closes #73
  * add python example
  * add --env flag back
  * add --dry-run. Closes #47
  * add function name inference
  * add runtime inference
  * add function config inheritance from project config
  * change logger to use cli handler

v0.3.0 / 2016-01-09
===================

  * add Project.Name function prefixes to prevent collisions
  * add Function.Prefix support
  * add initial log tailing
  * add basic config validation
  * remove old Function.Verbose field

v0.2.0 / 2016-01-06
===================

  * add concurrent deploys
  * add --log-level. Closes #42
  * add multi-function and project level management
  * add rollback support

v0.1.0 / 2016-01-04
===================

  * add test target to Makefile
  * add updating of configuration on deploy. Closes #25
  * add go generate directive for mockgen
  * add CI badge
  * add Function.Delete unit tests
  * add --verbose support to deploys. Closes #26
  * add delete command
  * add removal of build artifacts. Closes #16
  * add Getenv(), deploy --env and shim json file
  * add runtime target and main field support
  * add Kinesis handler
  * add invoke --async support
  * rename zip command to build
  * remove Getenv(), prime via os.Setenv() instead
  * change to use lambda.json

v0.0.2 / 2015-12-19
===================

  * add newline to invoke output
  * add Python support and example. Closes #2
  * remove invoke stderr newline
  * rename ./node to ./shim

v0.0.1 / 2015-12-19
===================

  * initial release
