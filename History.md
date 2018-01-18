
v1.0.0-rc2 / 2018-01-17
=======================

  * add new upgrade mechanism to match Up
  * add utils from Up which are necessary for upgrade
  * add native Go support. Closes #864

v0.16.0 / 2017-10-14
====================

  * add `./vendor`
  * refactor: move `process.env` out of the shim event handlers
  * refactor shim, allowing for concurrent calls to handler (#822)
  * remove analytics. Closes #777
  * fix node runtime
  * fix deploying java and clojure functions from zipfile (#820)
  * fix #815 - clojure deployments being doubled by java plugin (#816)
  * fix for issue #380 (#814)
  * fix: move -C chdir flag check to beginning of preparation function. (#766)

v0.15.0 / 2017-06-11
====================

  * add passing of maps that describe all apex functions to terraform Closes #737 (#744)
  * added nodejs4.3-edge runtime (#741)
  * add ListFunctions to min IAM policy example (#733)
  * add documentation to set alias to another alias (#732)
  * add arbitrary python version support. Closes #723
  * Add lambda:ListAliases to min IAM policy (#729)
  * add --alias docs
  * change DefaultRetainedVersions to 25
  * Update function deploy to update the specified alias even if code has not changed. Closes #686 (#712)
  * upgrade example to use nodejs version 6.10 (#719)
  * fix tabs in deploy example

v0.14.0 / 2017-03-31
====================

  * add support for deploying existing zips. Closes #480
  * change index.js inference to use v6.10
  * fix empty DeadLetterARN config triggering deploy no matter what. Possibly Closes #701 (#713)
  * fix webpack2 example for version 2.3.1 (#709)

v0.13.1 / 2017-03-23
====================

  * add Node v6.10 support. Closes #493
  * change runtime to nodejs43 for rust lambdas (#689)
  * fix context for gh ListReleases call (#691)
  * fix alias command example by removing v prefix: v5 should be just 5 (#681)

v0.13.0 / 2017-02-15
====================

  * add clojure support (#674)
  * add variable apex_function_NAME_name to the exported variable to terraform (Closes #654) (#660)
  * add support for DLQ ARN (#649)
  * add support for KMS ARN (#648)
  * add alias command (#647)

v0.12.0 / 2017-01-06
====================

  * add `exec` sub-command for pass-through of env vars. Closes #619
  * add improved message if curl failed due to permission error (#634)
  * add initial rust support (#549)
  * add function arn to list output (#627)
  * add IAM role support (#614)
  * fix shim runtime for new functions, now "nodejs4.3"

v0.11.0 / 2016-11-21
====================

  * add comparison of environment when performing change check
  * add support for APEX_ENVIRONMENT=prod apex deploy (#586)
  * add a minimum IAM policy to the documentation. (#529)
  * remove env hack, replace with native env support (#600)
  * fix overrides of VPC in function.json. Closes #597 (#598)
  * fix .apexignore fails when folder without trailling slash (#550)

v0.10.3 / 2016-09-05
====================

* add function.json fallback for any env if it is present. Closes #471
* add aliases to apex list. Closes #430 (#457)
* fix config comparison function. Closes #528 (#530)
* fix error handling for incorrect 'handler' definition for python (fixes #498) (#501)
* fix: use executable's path as the temp path (#488)
* fix: run hooks before runtime plugins to fix (#491) (#492)
* fix ignoring of function.json when defaultEnvironment is used. Closes #467
* fix for go-github api changes. Closes #463 (#465)
* fix custom clean hook was not working on golang runtime. Closes #461 (#464)

v0.10.2 / 2016-06-16
====================

  * fix version for upgrades

v0.10.1 / 2016-06-14
====================

  * add missing --env-file for build command. Closes #452

v0.10.0 / 2016-06-14
====================

  * add getting-started docs with init command
  * add infra docs
  * add role bootstrapping to `apex init` for a smoother experience
  * add project.json config "profile" support
  * add `--env-file` support. Closes #387
  * add multiple deployment targets for different environments (#432)
  * add missing `-s, --set` environment support to the build command. Closes #447
  * add retainedVersions as a pointer, letting you zero. Closes #407
  * refactor project init
  * remove api-gateway example since it is incomplete. Closes #445
  * fix flags before infra sub-command. Closes #421
  * fix panic when 'apex invoke -L' throws an error (#431)
  * fix setting handler func in Python runtime

v0.9.0 / 2016-05-10
===================

  * add symlink docs
  * add autocomplete docs
  * add dynamic auto-completion (#403)
  * add function globbing (#397)
  * add basic Segment analytics
  * add support for symlinks. Closes #320 (#375)
  * add marking not deployed functions on list (#374)
  * add multiple funcs rollback. Closes #372 (#373)
  * add support for updating function runtimes (apex deploy). Closes #369
  * add metric price (#368)
  * add implicit APEX_FUNCTION_NAME and LAMBDA_FUNCTION_NAME env vars
  * add apex environment as tfvar `apex_environment`
  * fix symlink support. Closes #401 #288 (#404)
  * fix install script, parse Github tag JSON with no newlines fixes #400 (#402)
  * fix missing logs. Closes #393
  * fix new lines in Hooks docs. Closes #390
  * fix generated boilerplate package
  * fix panic when trying to set incomplete -s var (#379)
  * fix infra command terraform flags pass-through. Closes #325 (#386)
  * fix missing defaultEnvironment in project.json. Closes #382

v0.8.0 / 2016-04-08
===================

  * add nodejs 4.3.2 support. Closes #356
  * add `defaultEnvironment` in project.json. Closes #338
  * add test for merging envs. Closes #348
  * rename --start to --since [breaking change]
  * refactor java plugin to expect a JAR file, rather than building one based on a pom file. [breaking change]
  * fix: api-gateway example Internal server error

v0.7.3 / 2016-03-23
===================

  * add passing aws_region var to Terraform
  * add API Gateway integration example. Closes #339
  * add env switch. Closes #304
  * add indent in init command. Closes #303
  * fix prompt.Confirm on Windows

v0.7.2 / 2016-03-15
===================

  * fix checking config changes. Closes #334

v0.7.1 / 2016-03-14
===================

  * add error message to install.sh if GitHub API call fails
  * add example of following logs with no historical output
  * add deploy, rollback, and invoke --alias support. Closes #7
  * add remote Terraform state init. Closes #299
  * add openbsd. Closes #307
  * add DEBUG_SHIM to out stdout output
  * refactor logging in function. Closes #84
  * fix detecting config changes. Closes #311
  * fix AWS config precedence

v0.7.0 / 2016-03-03
===================

  * add Project.LoadFunction(name)
  * add project.LoadFunctionByPath(name, path)
  * add `init` command
  * add `infra` command
  * add zip compression. Closes #263
  * add several new examples (browserify, webpack, java)
  * add zero-ing of mtime for all files. Closes #262
  * add java support . Closes #4
  * add installation script. Closes #199
  * add VPC support. Closes #242
  * add cleaning up old versions. Closes #148
  * add a flag for passing an AWS credentials file into Apex
  * refactor rollback to use flag instead of arg. Closes #289
  * refactor Java plugin; no pom.xml generation
  * refactor env flag; moved to deploy command
  * fix EACCESS error caused by missing exec bit. Closes #281
  * fix pattern matching for .gitignore parity. Closes #228
  * fix loading region from default profile
  * fix path separator handling in windows. Closes #222.
  * fix: build all Golang source files for function

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
