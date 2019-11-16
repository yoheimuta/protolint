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
    --protolint_out=v,fix,config_dir_path=_example/config,reporter=junit,plugin=./plugin_example \
    *.proto
```

