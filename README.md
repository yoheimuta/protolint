# protolint [![CircleCI](https://circleci.com/gh/yoheimuta/protolint/tree/master.svg?style=svg)](https://circleci.com/gh/yoheimuta/protolint/tree/master)[![Go Report Card](https://goreportcard.com/badge/github.com/yoheimuta/protolint)](https://goreportcard.com/report/github.com/yoheimuta/protolint)[![License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/yoheimuta/protolint/blob/master/LICENSE)

protolint is a command line tool which lints Protocol Buffer files (proto3):

- Runs fast because this works without compiler.
- Easy to follow the official style guide. The rules and the style guide correspond to each other exactly.
- Allow to disable rules with a comment in a Protocol Buffer file.
  - It is useful for projects which must keep API compatibility while enforce the style guide as much as possible.
- Undergone testing for all rules.

## Installation

```
go get -u -v github.com/yoheimuta/protolint/cmd/protolint
```

For non-Go users, the simplest way to install the protolint is to download a pre-built binary from this release page:

- https://github.com/yoheimuta/protolint/releases

In the downloads section of each release, you can find pre-built binaries in .tar.gz packages.

## Usage

```
protolint lint example.proto example2.proto # file mode, specify multiple specific files
protolint lint .                            # directory mode, search for all .proto files recursively
protolint .                                 # same as "protolint lint ."
protolint lint -config_dir_path=path/to .   # search path/to for .protolint.yaml
protolint lint -fix .                       # automatically fix some of the problems reported by some rules
protolint lint -v .                         # with verbose output to investigate the parsing error
protolint list                              # list all current lint rules being used
```

## Rules

See `internal/addon/rules` in detail.

The rule set follows:

- [Official Style Guide](https://developers.google.com/protocol-buffers/docs/style). This is enabled by default.
- Formatting Style Guide. This is enabled by default.
  - Enforce a maximum line length. The length of a line is defined as the number of Unicode characters in the line. You can configure the detail with `.protolint.yaml`.
  - Enforce a consistent indentation style. The --fix option on the command line can automatically fix some of the problems reported by this rule. The default style is 4 spaces. You can configure the detail with `.protolint.yaml`.

| ID                                | Purpose                                                                  |
|-----------------------------------|--------------------------------------------------------------------------|
| ENUM_FIELD_NAMES_UPPER_SNAKE_CASE | Verifies that all enum field names are CAPITALS_WITH_UNDERSCORES.        |
| ENUM_NAMES_UPPER_CAMEL_CASE       | Verifies that all enum names are CamelCase (with an initial capital).    |
| FIELD_NAMES_LOWER_SNAKE_CASE      | Verifies that all field names are underscore_separated_names.            |
| MESSAGE_NAMES_UPPER_CAMEL_CASE    | Verifies that all message names are CamelCase (with an initial capital). |
| RPC_NAMES_UPPER_CAMEL_CASE        | Verifies that all rpc names are CamelCase (with an initial capital).     |
| SERVICE_NAMES_UPPER_CAMEL_CASE    | Verifies that all service names are CamelCase (with an initial capital). |
| MAX_LINE_LENGTH    | Enforces a maximum line length. |
| INDENT    | Enforces a consistent indentation style. |
| SERVICE_NAMES_END_WITH    | Enforces a consistent suffix for service names. |

`-` is a bad style, `+` is a good style:

__ENUM_FIELD_NAMES_UPPER_SNAKE_CASE__

```diff
enum Foo {
-  firstValue = 0;
+  FIRST_VALUE = 0;
-  second_value = 1;
+  SECOND_VALUE = 1;
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

__MESSAGE_NAMES_UPPER_CAMEL_CASE__

```diff
- message song_server_request {
+ message SongServerRequest {
  required string SongName = 1;
  required string song_name = 1;
}
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

__Config file__

protolint can operate using a config file named `.protolint.yaml`.

Refer to [_example/config/.protolint.yaml](_example/config/.protolint.yaml) for the config file specification.

protolint will search a current working directory for the config file by default.
And it can search the specified directory with `-config_dir_path` flag.

## Motivation

There exists the similar protobuf linters as of 2018-12-20.

One is a plug-in for Google's Protocol Buffers compiler.

- When you just want to lint the files, it may be tedious to create the compilation environment.
- And it generally takes a lot of time to compile the files than to parse the files.

Other is a command line tool which also lints Protocol Buffer files.

- While it has a lot of features other than lint, it seems cumbersome for users who just want the linter.
- The lint rule slants towards to be opinionated.
- Further more, the rule set and the official style guide don't correspond to each other exactly. It requires to understand both rules and the guide in detail, and then to combine the rules accurately.

## Dependencies

- [go-protoparser](https://github.com/yoheimuta/go-protoparser)

## Acknowledgement

Thank you to the prototool package: https://github.com/uber/prototool

I referred to the package for the good proven design, interface and some source code.
