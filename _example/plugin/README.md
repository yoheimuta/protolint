# Plugin Example

### Build

```bash
go build -o plugin_example main.go
```

NOTE: protolint plugin is backed by [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin), not the [plugin](https://golang.org/pkg/plugin/) standard library.

Therefore, you can build the plugin just as a normal Go main package.

### Run

```bash
protolint -plugin ./plugin_example /path/to/files

# Or you can pass some flags to your plugin:
protolint -plugin "./plugin_example -go_style=false" /path/to/files

# You can see that your plugin is loaded correctly.
protolint list -plugin ./plugin_example
```

NOTE: `sh` must be in your PATH.

NOTE2: Even when you specify the plugin, the configuration defined in `.protolint.yaml` can stop your plugin from working. Check the yaml when your plugin doesn't appear to be working. See https://github.com/yoheimuta/protolint/issues/260 in detail.
