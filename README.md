# protolint
![Action](https://github.com/yoheimuta/protolint/workflows/Go/badge.svg)
[![Release](https://img.shields.io/github/v/release/yoheimuta/protolint?include_prereleases)](https://github.com/yoheimuta/protolint/releases)[
![Go Report Card](https://goreportcard.com/badge/github.com/yoheimuta/protolint)](https://goreportcard.com/report/github.com/yoheimuta/protolint)
[![License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/yoheimuta/protolint/blob/master/LICENSE)
[![Docker](https://img.shields.io/docker/pulls/yoheimuta/protolint)](https://hub.docker.com/r/yoheimuta/protolint)

protolint is the pluggable linting/fixing utility for Protocol Buffer files (proto2+proto3):

- Runs fast because this works without compiler.
- Easy to follow the official style guide. The rules and the style guide correspond to each other exactly.
  - Fixer automatically fixes all the possible official style guide violations.
- Allows to disable rules with a comment in a Protocol Buffer file.
  - It is useful for projects which must keep API compatibility while enforce the style guide as much as possible.
  - Some rules can be automatically disabled by inserting comments to the spotted violations.
- Loads plugins to contain your custom lint rules.
- Undergone testing for all rules.
- Many integration supports.
  - protoc plugin
  - Editor integration
  - GitHub Action
  - CI Integration

## Demo

For example, vim-protolint works like the following.

<img src="_doc/demo-v2.gif" alt="demo" width="600"/>

## Installation

### Via Homebrew

protolint can be installed for Mac or Linux using Homebrew via the [yoheimuta/protolint](https://github.com/yoheimuta/homebrew-protolint) tap.

```sh
brew tap yoheimuta/protolint
brew install protolint
```

Since [homebrew-core](https://github.com/Homebrew/homebrew-core/pkgs/container/core%2Fprotolint) includes `protolint,` you can also install it by just `brew install protolint.` This is the default tap that is installed by default. It's easier, but not maintained by the same author. To keep it updated, I recommend you run `brew tap yoheimuta/protolint` first.


### Via GitHub Releases

You can also download a pre-built binary from this release page:

- https://github.com/yoheimuta/protolint/releases

In the downloads section of each release, you can find pre-built binaries in .tar.gz packages.

### Use the maintained Docker image

protolint ships a Docker image [yoheimuta/protolint](https://hub.docker.com/r/yoheimuta/protolint) that allows you to use protolint as part of your Docker workflow.

```
❯❯❯ docker run --volume "$(pwd):/workspace" --workdir /workspace yoheimuta/protolint lint _example/proto
[_example/proto/invalidFileName.proto:1:1] File name should be lower_snake_case.proto.
[_example/proto/issue_88/oneof_options.proto:11:5] Found an incorrect indentation style "    ". "  " is correct.
[_example/proto/issue_88/oneof_options.proto:12:5] Found an incorrect indentation style "    ". "  " is correct.
```

### From Source

The binary can be installed from source if Go is available.
However, I recommend using one of the pre-built binaries instead because it doesn't include the version info.

```sh
go install github.com/yoheimuta/protolint/cmd/protolint@latest
```

### Within JavaScript / TypeScript

You can use `protolint` using your nodejs package manager like `npm` or `yarn`.

```sh
$ npm install protolint --save-dev
```

This will add a reference to a development dependency to your local `package.json`.

During install, the [install.mjs](bdist/js/install.mjs) script will be called. It will download the matching `protolint` from github. Just like [@electron/get](https://github.com/electron/get/), you can bypass the download using the following environment variables:

| Environment Variable          | Default value                         | Description                                   |
|-------------------------------|---------------------------------------|-----------------------------------------------|
| PROTOLINT_MIRROR_HOST         | https://github.com                    | HTTP/Web server base url hosting the binaries |
| PROTOLINT_MIRROR_REMOTE_PATH  | yoheimuta/protolint/download/releases | Path to the archives on the remote host       |
| PROTOLINT_MIRROR_USERNAME     |                                       | HTTP Basic auth user name                     |
| PROTOLINT_MIRROR_PASSWORD     |                                       | HTTP Basic auth password                      |
| PROTOLINT_PROXY               |                                       | HTTP(S) Proxy with optional auth data         |

Within the remote path, the archives from the [releases](https://github.com/yoheimuta/protolint/releases/latest/) page must be
mirrored.

After that, you can use `npx protolint` (with all supplied protolint arguments) within your dev-scripts.

```json
{
  ...
  "scripts": {
    "protoc": "....",
    "preprotoc": "npx protolint"
  },
  ...
}
```

You can add a `protolint` node to your `package.json` which may contain the content of `protolint.yml` below the `lint` node, i.e. the root element of the configuration will be `protolint`.

If you want to get an output that matches the TSC compiler, use reporter `tsc`.

### Within Python projects

You can use `protolint` as a linter within your python projects, the wheel `protolint-bin` on [pypi](https://pypi.org) contains the pre-compiled binaries for various platforms. Just add the desired version to
your `pyproject.toml` or `requirements.txt`.

The wheels downloaded will contain the compiled go binaries for `protolint` and `protoc-gen-protolint`. Your platform must
be compatible with the supported binary platforms.

You can add the linter configuration to the `tools.protolint` package in `pyproject.toml`.

## Usage

```sh
protolint lint example.proto example2.proto # file mode, specify multiple specific files
protolint lint .                            # directory mode, search for all .proto files recursively
protolint .                                 # same as "protolint lint ."
protolint lint -config_path=path/to/your_protolint.yaml . # use path/to/your_protolint.yaml
protolint lint -config_dir_path=path/to .   # search path/to for .protolint.yaml
protolint lint -fix .                       # automatically fix some of the problems reported by some rules
protolint lint -fix -auto_disable=next .    # this is preferable when you want to fix problems while maintaining the compatibility. Automatically fix some problems and insert disable comments to the other problems. The available values are next and this.
protolint lint -auto_disable=next .         # automatically insert disable comments to the other problems. 
protolint lint -v .                         # with verbose output to investigate the parsing error
protolint lint -no-error-on-unmatched-pattern . # exits with success code even if no file is found (file & directory mode)
protolint lint -reporter junit .            # output results in JUnit XML format
protolint lint -output_file=path/to/out.txt # output results to path/to/out.txt
protolint lint -plugin ./my_custom_rule1 -plugin ./my_custom_rule2 .   # run custom lint rules.
protolint list                              # list all current lint rules being used
protolint version                           # print protolint version
```

protolint does not require configuration by default, for the majority of projects it should work out of the box.

## Version Control Integration

protolint is available as a [pre-commit](https://pre-commit.com) hook.  Add this to your `.pre-commit-config.yaml` in your repository to run protolint with Go:
```yaml
repos:
  - repo: https://github.com/yoheimuta/protolint
    rev: <version> # Select a release here like v0.44.0
    hooks:
      - id: protolint
```
or alternatively use this to run protolint with Docker:
```yaml
repos:
  - repo: https://github.com/yoheimuta/protolint
    rev: <version> # Select a release here like v0.44.0
    hooks:
      - id: protolint-docker
```

## Editor Integration

Visual Studio Code

- [vscode-protolint](https://github.com/plexsystems/vscode-protolint)

JetBrains IntelliJ IDEA, GoLand, WebStorm, PHPStorm, PyCharm...

- [intellij-protolint](https://github.com/yoheimuta/intellij-protolint)

Vim([ALE engine](https://github.com/dense-analysis/ale))

- [ale](https://github.com/dense-analysis/ale)'s [built-in support](https://github.com/dense-analysis/ale/blob/master/supported-tools.md)

Vim([Syntastic](https://github.com/vim-syntastic/syntastic))

- [vim-protolint](https://github.com/yoheimuta/vim-protolint)

## GitHub Action

A [GitHub Action](https://github.com/features/actions) to run protolint in your workflows

- [github/super-linter](https://github.com/github/super-linter)
- [plexsystems/protolint-action](https://github.com/plexsystems/protolint-action)
- [yoheimuta/action-protolint](https://github.com/yoheimuta/action-protolint) - Integrated with [reviewdog](https://github.com/reviewdog/reviewdog)

## CI Integration

Jenkins Plugins

- [warnings-ng](https://github.com/jenkinsci/warnings-ng-plugin) and any that use [violatons-lib](https://github.com/tomasbjerre/violations-lib)

### Environment specific output

It is possible to format your linting according to the formatting of the CI/CD environment. The environment must be set using the output format. Currently, the following output is realized:

| Environment | Command Line Value | Description | Example |
|-------------|--------------------|-------------|---------|
| Github Actions | ci-gh | [Github Help](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-a-notice-message) | `::warning file=example.proto,line=10,col=20,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES` |
| Azure DevOps | ci-az | [Azure DevOps Help](https://learn.microsoft.com/en-us/azure/devops/pipelines/scripts/logging-commands?view=azure-devops&tabs=bash#task-commands) | `##vso[task.logissue type=warning;sourcepath=example.proto;linenumber=10;columnnumber=20;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES` |
| Gitlab CI/CD | ci-glab | Reverse Engineered from Examples | `WARNING: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(10,20) : EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES` |

You can also use the generic `ci` formatter, which will create a generic problem matcher.

With the `ci-env` value, you can specify the template from the following environment variables:

| Environment Variable | Priority | Meaning |
|----------------------|----------|---------|
| PROTOLINT_CIREPORTER_TEMPLATE_STRING | 1 | String containing a Go-template |
| PROTOLINT_CIREPORTER_TEMPLATE_FILE | 2 | Path to a file containing a Go-template |

The resulting line-feed must not be added, as it will be added automatically.

The following fields are available:

`Severity`
: The severity as string (either note, warning or error)

`File`
: Path to the file containing the error

`Line`
: Line within the `file` containing the error (starting position)

`Column`
: Column within the `file` containing the error (starting position)

`Rule`
: The name of the rule that is faulting

`Message`
: The error message that descibes the error

### Producing an output file and an CI/CD Error stream

You can create a specific output matching your CI/CD environment and also create an output file, e.g. for your static code analysis tools like github CodeQL or SonarQube.

This can be done by adding the `--add-reporter` flag.
Please note, that the value must be formatted `<reporter-name>:<output-file-path>` (omitting `<` and `>`).

```shell
$ protolint --reporter ci-gh --add-reporter sarif:/path/to/my/output.sarif.json proto/*.proto
```

## Use as a protoc plugin

protolint also maintains a binary [protoc-gen-protolint](cmd/protoc-gen-protolint) that performs the lint functionality as a protoc plugin.
See [cmd/protoc-gen-protolint/README.md](https://github.com/yoheimuta/protolint/blob/master/cmd/protoc-gen-protolint/README.md) in detail.

This is useful in situations where you already have a protoc plugin workflow.

## Call from Go code

You can also use protolint from Go code.
See [Go Documentation](https://pkg.go.dev/github.com/yoheimuta/protolint/lib) and [lib/lint_test.go](https://github.com/yoheimuta/protolint/blob/master/lib/lint_test.go) in detail.

```go
args := []string{"-config_path", "path/to/your_protolint.yaml", "."}
var stdout bytes.Buffer
var stderr bytes.Buffer

err := lib.Lint(test.inputArgs, &stdout, &stderr)
```

## Rules

See `internal/addon/rules` in detail.

The rule set follows:

- [Official Style Guide](https://protobuf.dev/programming-guides/style/). This is enabled by default. Basically, these rules can fix the violations by appending `-fix` option.
- Unofficial Style Guide. This is disabled by default. You can enable each rule with `.protolint.yaml`.

The `-fix` option on the command line can automatically fix all the problems reported by fixable rules.
See Fixable columns below.

The `-auto_disable` option on the command line can automatically disable all the problems reported by auto-disable rules.
This feature is helpful when fixing the existing violations breaks the compatibility.
See AutoDisable columns below.

- *1: These rules are not supposed to support AutoDisable because the fixes don't break their compatibilities. You should run the protolint with `-fix`.

| Official | Fixable | AutoDisable | ID                                | Purpose                                                                  |
|----------|---------|---------|-----------------------------------|--------------------------------------------------------------------------|
| Yes | ✅ | ✅ | ENUM_FIELD_NAMES_PREFIX | Verifies that enum field names are prefixed with its ENUM_NAME_UPPER_SNAKE_CASE.        |
| Yes | ✅ | ✅ | ENUM_FIELD_NAMES_UPPER_SNAKE_CASE | Verifies that all enum field names are CAPITALS_WITH_UNDERSCORES.        |
| Yes | ✅ | ✅ | ENUM_FIELD_NAMES_ZERO_VALUE_END_WITH | Verifies that the zero value enum should have the suffix (e.g. "UNSPECIFIED", "INVALID"). The default is "UNSPECIFIED". You can configure the specific suffix with `.protolint.yaml`. |
| Yes | ✅ | ✅ | ENUM_NAMES_UPPER_CAMEL_CASE       | Verifies that all enum names are CamelCase (with an initial capital).    |
| Yes | ✅ | *1 | FILE_NAMES_LOWER_SNAKE_CASE       | Verifies that all file names are lower_snake_case.proto. You can configure the excluded files with `.protolint.yaml`. |
| Yes | ✅ | ✅ | FIELD_NAMES_LOWER_SNAKE_CASE      | Verifies that all field names are underscore_separated_names.            |
| Yes | ✅ | *1 | IMPORTS_SORTED                    | Verifies that all imports are sorted. |
| Yes | ✅ | ✅ | MESSAGE_NAMES_UPPER_CAMEL_CASE    | Verifies that all message names are CamelCase (with an initial capital). |
| Yes | ✅ | *1 | ORDER                             | Verifies that all files should be ordered in the specific manner. |
| Yes | ✅ | *1 | PACKAGE_NAME_LOWER_CASE           | Verifies that the package name should only contain lowercase letters. |
| Yes | ✅ | ✅ | RPC_NAMES_UPPER_CAMEL_CASE        | Verifies that all rpc names are CamelCase (with an initial capital).     |
| Yes | ✅ | ✅ | SERVICE_NAMES_UPPER_CAMEL_CASE    | Verifies that all service names are CamelCase (with an initial capital). |
| Yes | ✅ | ✅ | REPEATED_FIELD_NAMES_PLURALIZED   | Verifies that repeated field names are pluralized names.            |
| Yes | ✅ | *1 | QUOTE_CONSISTENT   | Verifies that the use of quote for strings is consistent. The default is double quoted. You can configure the specific quote with `.protolint.yaml`.          |
| Yes | ✅ | *1 | INDENT    | Enforces a consistent indentation style. The default style is 2 spaces. Inserting appropriate new lines is also forced by default. You can configure the detail with `.protolint.yaml`. |
| Yes | ✅ | *1 | PROTO3_FIELDS_AVOID_REQUIRED      | Verifies that all fields should avoid required for proto3.            |
| Yes | _  | ✅ | PROTO3_GROUPS_AVOID      | Verifies that all groups should be avoided for proto3.            |
| Yes | _  | *1 | MAX_LINE_LENGTH    | Enforces a maximum line length. The length of a line is defined as the number of Unicode characters in the line. The default is 80 characters. You can configure the detail with `.protolint.yaml`. |
| No | _  | - | SERVICE_NAMES_END_WITH    | Enforces a consistent suffix for service names. You can configure the specific suffix with `.protolint.yaml`. |
| No | _  | - | FIELD_NAMES_EXCLUDE_PREPOSITIONS | Verifies that all field names don't include prepositions (e.g. "for", "during", "at"). You can configure the specific prepositions and excluded keywords with `.protolint.yaml`. |
| No | _  | - | MESSAGE_NAMES_EXCLUDE_PREPOSITIONS | Verifies that all message names don't include prepositions (e.g. "With", "For"). You can configure the specific prepositions and excluded keywords with `.protolint.yaml`. |
| No | _  | - | RPC_NAMES_CASE        | Verifies that all rpc names conform to the specified convention. You need to configure the specific convention with `.protolint.yaml`.     |
| No | _  | - | MESSAGES_HAVE_COMMENT | Verifies that all messages have a comment. You can configure to enforce Golang Style comments with `.protolint.yaml`. |
| No | _  | - | SERVICES_HAVE_COMMENT | Verifies that all services have a comment. You can configure to enforce Golang Style comments with `.protolint.yaml`. |
| No | _  | - | RPCS_HAVE_COMMENT | Verifies that all rps have a comment. You can configure to enforce Golang Style comments with `.protolint.yaml`. |
| No | _  | - | FIELDS_HAVE_COMMENT | Verifies that all fields have a comment. You can configure to enforce Golang Style comments with `.protolint.yaml`. |
| No | _  | - | ENUMS_HAVE_COMMENT | Verifies that all enums have a comment. You can configure to enforce Golang Style comments with `.protolint.yaml`. |
| No | _  | - | ENUM_FIELDS_HAVE_COMMENT | Verifies that all enum fields have a comment. You can configure to enforce Golang Style comments with `.protolint.yaml`. |
| No | _  | - | FILE_HAS_COMMENT | Verifies that a file starts with a doc comment. |
| No | _  | - | SYNTAX_CONSISTENT | Verifies that syntax is a specified version. The default is proto3. You can configure the version with `.protolint.yaml`. |

I recommend that you add `all_default: true` in `.protolint.yaml`, because all linters above are automatically enabled so that you can always enjoy maximum benefits whenever protolint is updated.

Here are some examples that show good style enabled by default.
`-` is a bad style, `+` is a good style:

__ENUM_FIELD_NAMES_PREFIX__

```diff
enum FooBar {
-  UNSPECIFIED = 0;
+  FOO_BAR_UNSPECIFIED = 0;
}
```

__ENUM_FIELD_NAMES_UPPER_SNAKE_CASE__

```diff
enum Foo {
-  firstValue = 0;
+  FIRST_VALUE = 0;
-  second_value = 1;
+  SECOND_VALUE = 1;
}
```

__ENUM_FIELD_NAMES_ZERO_VALUE_END_WITH__

```diff
enum Foo {
-  FOO_FIRST = 0;
+  FOO_UNSPECIFIED = 0;
}
```

__ENUM_NAMES_UPPER_CAMEL_CASE__

```diff
- enum foobar {
+ enum FooBar {
  FIRST_VALUE = 0;
  SECOND_VALUE = 1;
}
```

__FIELD_NAMES_LOWER_SNAKE_CASE__

```diff
message SongServerRequest {
-  required string SongName = 1;
+  required string song_name = 1;
}
```

__IMPORTS_SORTED__

```diff
- import public "new.proto";
+ import "myproject/other_protos.proto";
- import "myproject/other_protos.proto";
+ import public "new.proto";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
```

__MESSAGE_NAMES_UPPER_CAMEL_CASE__

```diff
- message song_server_request {
+ message SongServerRequest {
  required string SongName = 1;
  required string song_name = 1;
}
```

__ORDER__

```diff
- option java_package = "com.example.foo";
- syntax = "proto3";
- package examplePb;
- message song_server_request { }
- import "other.proto";
+ syntax = "proto3";
+ package examplePb;
+ import "other.proto";
+ option java_package = "com.example.foo";
+ message song_server_request { }
```

__PACKAGE_NAME_LOWER_CASE__

```diff
- package myPackage
+ package my.package
```

__RPC_NAMES_UPPER_CAMEL_CASE__

```diff
service FooService {
-  rpc get_something(FooRequest) returns (FooResponse);
+  rpc GetSomething(FooRequest) returns (FooResponse);
}
```

__RPC_NAMES_UPPER_CAMEL_CASE__

```diff
- service foo_service {
+ service FooService {
  rpc get_something(FooRequest) returns (FooResponse);
  rpc GetSomething(FooRequest) returns (FooResponse);
}
```

__REPEATED_FIELD_NAMES_PLURALIZED__

```diff
-  repeated string song_name = 1;
+  repeated string song_names = 1;
```

__INDENT__

```diff
 enum enumAllowingAlias {
   UNKNOWN = 0;
-        option allow_alias = true;
+  option allow_alias = true;
   STARTED = 1;
-     RUNNING = 2 [(custom_option) = "hello world"];
+  RUNNING = 2 [(custom_option) = "hello world"];
- }
+}
```

```diff
-   message TestMessage { string test_field = 1; }
+ message TestMessage {
+  string test_field = 1;
+}
```

__QUOTE_CONSISTENT__

```diff
 option java_package = "com.example.foo";
- option go_package = 'example';
+ option go_package = "example";
```

## Creating your custom rules

protolint is the pluggable linter so that you can freely create custom lint rules.

A complete sample project (aka plugin) is included in this repo under the [_example/plugin](_example/plugin) directory.

## Reporters

protolint comes with several built-in reporters(aka. formatters) to control the appearance of the linting results.

You can specify a reporter using the -reporter flag on the command line. For example, `-reporter junit` uses the junit reporter.

The built-in reporter options are:

- plain (default)
- junit
- json
- sarif
- sonar (SonarQube generic issue format)
- unix
- tsc (compatible to TypeScript compiler)

## Configuring

__Disable rules in a Protocol Buffer file__

Rules can be disabled with a comment inside a Protocol Buffer file with the following format.
The rules will be disabled until the end of the file or until the linter sees a matching enable comment:

```
// protolint:disable <ruleID1> [<ruleID2> <ruleID3>...]
...
// protolint:enable <ruleID1> [<ruleID2> <ruleID3>...]
```

It's also possible to modify a disable command by appending :next or :this for only applying the command to this(current) or the next line respectively.

For example:

```proto
enum Foo {
  // protolint:disable:next ENUM_FIELD_NAMES_UPPER_SNAKE_CASE
  firstValue = 0;    // no error
  second_value = 1;  // protolint:disable:this ENUM_FIELD_NAMES_UPPER_SNAKE_CASE
  THIRD_VALUE = 2;   // spits out an error
}
```

Setting the command-line option `-auto_disable` to `next` or `this` inserts disable commands whenever spotting problems. 

You can specify `-fix` option together. The rules supporting auto_disable suppress the violations instead of fixing them that cause a schema incompatibility.

__Config file__

protolint can operate using a config file named `.protolint.yaml`.

Refer to [_example/config/.protolint.yaml](_example/config/.protolint.yaml) for the config file specification.

protolint will automatically search a current working directory for the config file by default
and successive parent directories all the way up to the root directory of the filesystem.
And it can search the specified directory with `-config_dir_path` flag.
It can also search the specified file with `--config_path` flag.

## Exit codes

When linting files, protolint will exit with one of the following exit codes:

- `0`: Linting was successful and there are no linting errors.
- `1`: Linting was successful and there is at least one linting error.
- `2`: Linting was unsuccessful due to all other errors, such as parsing, internal, and runtime errors.

## Motivation

There exists the similar protobuf linters as of 2018/12/20.

One is a plug-in for Google's Protocol Buffers compiler.

- When you just want to lint the files, it may be tedious to create the compilation environment.
- And it generally takes a lot of time to compile the files than to parse the files.

Other is a command line tool which also lints Protocol Buffer files.

- While it has a lot of features other than lint, it seems cumbersome for users who just want the linter.
- The lint rule slants towards to be opinionated.
- Further more, the rule set and the official style guide don't correspond to each other exactly. It requires to understand both rules and the guide in detail, and then to combine the rules accurately.

### Other tools

I wrote an article comparing various Protocol Buffer Linters, including protolint, on 2019/12/17.

- https://qiita.com/yoheimuta/items/da7678fcd046b93a2637
  - NOTE: This one is written in Japanese.

## Dependencies

- [go-protoparser](https://github.com/yoheimuta/go-protoparser)

## License

The MIT License (MIT)

## Acknowledgement

Thank you to the prototool package: https://github.com/uber/prototool

I referred to the package for the good proven design, interface and some source code.
