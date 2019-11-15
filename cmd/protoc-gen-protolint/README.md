# protoc-gen-protolint

## Installation

Download `protoc-gen-protolint` and make sure it's available in your PATH. Once it's
in your PATH, `protoc` will be able to make use of the plug-in.

### From Source

The binary can be installed from source if Go is available.

```
go get -u -v github.com/yoheimuta/protolint/cmd/protoc-gen-protolint
```

## Usage

```
protoc --protolint_out=. *.proto
```

All flags, which is supported by protolint are passed as an option to the plugin as a comma separated text. It should look like below.

```
protoc \
    --protolint_out=v,fix,config_dir_path=_example/config,reporter=junit,plugin=./plugin_example \
    *.proto
```

