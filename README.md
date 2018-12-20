# protolinter [![GoDoc](https://godoc.org/github.com/yoheimuta/protolinter?status.svg)](https://godoc.org/github.com/yoheimuta/protolinter)

protolinter is a command line tool which lints Protocol Buffer files (proto3):

- Runs fast because this works without compiler.
- Easy to follow the official style guide. The rules and the style guide correspond to each other exactly.
- Undergone testing for all rules.

### Installation

```
go get -u -v github.com/yoheimuta/protolinter/cmd/pl
```

### Usage

```
pl lint example.proto example2.proto # file mode, specify multiple specific files
pl lint .                            # directory mode, search for all .proto files recursively
pl .                                 # same as "pl lint ."
pl list                              # list all current lint rules being used
```

### Rules

See `internal/addon/rules` in detail.

The rule set follows:

- [the official Style Guide](https://developers.google.com/protocol-buffers/docs/style).

| ID                                | Purpose                                                                  |
|-----------------------------------|--------------------------------------------------------------------------|
| ENUM_FIELD_NAMES_UPPER_SNAKE_CASE | Verifies that all enum field names are CAPITALS_WITH_UNDERSCORES.        |
| ENUM_NAMES_UPPER_CAMEL_CASE       | Verifies that all enum names are CamelCase (with an initial capital).    |
| FIELD_NAMES_LOWER_SNAKE_CASE      | Verifies that all field names are underscore_separated_names.            |
| MESSAGE_NAMES_UPPER_CAMEL_CASE    | Verifies that all message names are CamelCase (with an initial capital). |
| RPC_NAMES_UPPER_CAMEL_CASE        | Verifies that all rpc names are CamelCase (with an initial capital).     |
| SERVICE_NAMES_UPPER_CAMEL_CASE    | Verifies that all service names are CamelCase (with an initial capital). |

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

### Motivation

There exists the similar protobuf linters as of 2018-12-20.

[protoc-gen-lint](https://github.com/ckaznocha/protoc-gen-lint) is a plug-in for Google's Protocol Buffers compiler.

- When you just want to lint the files, it may be tedious to create the compilation environment.
- And it generally takes a lot of time to compile the files than to parse the files.

[prototool](https://github.com/uber/prototool) is a Swiss Army Knife for Protocol Buffers.

- The lint rule slants towards to be opinionated, because the rule set basically follows the Uber's Style Guide.
- While it has a lot of features other than lint, it seems cumbersome for users who just want the linter.
- Further more, the rule set and the official style guide don't correspond to each other exactly. It requires to understand both rules and the guide in detail, and then to combine the rules accurately.
- There are no tests about linter rules.

### TODO

- [ ] Enable to turn off the rule with a comment line
- [ ] Auto-Register binaries to GitHub Releases
- [ ] Support a configuration file to turn on/off the rules
- [ ] More rules
  - [ ] oneofNamesUpperCamelCaseRule
  - [ ] mapFieldNamesLowerSnakeCaseRule
  - [ ] oneofFieldNamesLowerSnakeCaseRule
  - [ ] requestResponseNamesMatchRPCRule
  - [ ] fieldNumbersStartFromOneRule
  - [ ] enumFieldNumbersStartFromZeroRule
  - [ ] enumsHaveCommentsRule
  - [ ] enumFieldsHaveCommentsRule
  - [ ] messagesHaveCommentsRule
  - [ ] fieldsHaveCommentsRule
  - [ ] mapFieldsHaveCommentsRule
  - [ ] rpcsHaveCommentsRule
  - [ ] servicesHaveCommentsRule
  - [ ] oneofsHaveCommentsRule
  - [ ] oneofFieldsHaveCommentsRule
  - [ ] commentsBeginWithTheNameRule
  - [ ] commentsNoCStyleRule
  - [ ] maxLineLengthRule
  - [ ] spaceIndentationRule
  - [ ] tabIndentationRule
