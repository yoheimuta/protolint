# protoc-gen-protolint

## Installation

Download `protoc-gen-protolint` and make sure it's available in your PATH. Once it's
in your PATH, `protoc` will be able to make use of the plug-in.

### Via Homebrew

protoc-gen-protolint can be installed for Mac or Linux using Homebrew via the [yoheimuta/protolint](https://github.com/yoheimuta/homebrew-protolint) tap.

```
brew tap yoheimuta/protolint
brew install protolint
```

### Via GitHub Releases

You can also download a pre-built binary from this release page:

- https://github.com/yoheimuta/protolint/releases

In the downloads section of each release, you can find pre-built binaries in .tar.gz packages.

### Via Maven Central

This plugin is also available on Maven Central. For details about how to use it, check out the [gradle
example](../../_example/gradle).

### From Source

The binary can be installed from source if Go is available.
However, I recommend using one of the pre-built binaries instead because it doesn't include the version info.

```
go get -u -v github.com/yoheimuta/protolint/cmd/protoc-gen-protolint
```

## Usage

```
protoc --protolint_out=. *.proto
```

A version subcommand is supported.

```
protoc-gen-protolint version
```

All flags, which is supported by protolint are passed as an option to the plugin as a comma separated text. It should look like below.

```
protoc \
    --protolint_out=. \
    --protolint_opt=v,fix,config_dir_path=_example/config,reporter=junit,plugin=./plugin_example \
    *.proto
```

### With [Grpc.Tools package (.NET Build)](https://chromium.googlesource.com/external/github.com/grpc/grpc/+/HEAD/src/csharp/BUILD-INTEGRATION.md)

When you specify `ProtoRoot`, make sure to add `--proto_root` option like the below.

```
<ItemGroup>
  <Protobuf Include="protos\**\*.proto" AdditionalProtocArguments="--protolint_out=.;--protolint_opt=proto_root=protos" ProtoRoot="protos" />
</ItemGroup>
```

## Option

### proto_root

If you add [protoc's --proto_path](https://developers.google.com/protocol-buffers/docs/proto3#generating) to read your proto files in the specified directory,
protolint could fail to locate the proto files. You should tell protolint the root directory like the below.

```
❯ ls protos
helloworld.proto

❯ protoc \
    --proto_path=protos
    --protolint_out=. \
    --protolint_opt=proto_root=protos \
    helloworld.proto
```
