# Plugin Example

### Build

```bash
go build main.go
```

NOTE: protolint plugin is backed by [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin), not the [plugin](https://golang.org/pkg/plugin/) standard library.

Therefore, you can build the plugin just as a normal Go main package.

### Run

```bash
protolint -plugin ./your_plugin_binary /path/to/files
```
