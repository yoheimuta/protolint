package protocgenprotolint

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds/lint"

	"github.com/golang/protobuf/proto"
	protogen "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/yoheimuta/protolint/internal/osutil"
)

var (
	version  = "master"
	revision = "latest"
)

const (
	subCmdVersion = "version"
)

// Do runs the command logic.
func Do(
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	if 0 < len(args) && args[0] == subCmdVersion {
		return doVersion(stdout)
	}

	err := signalSupportProto3Optional()
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return osutil.ExitInternalFailure
	}

	subCmd, err := newSubCmd(stdin, stdout, stderr)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return osutil.ExitInternalFailure
	}
	return subCmd.Run()
}

func signalSupportProto3Optional() error {
	// supports proto3 field presence
	// See https://github.com/protocolbuffers/protobuf/blob/cdc11c2d2d314ce0382fe0eaa715e5e0e1270438/docs/implementing_proto3_presence.md#signaling-that-your-code-generator-supports-proto3-optional
	var supportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	data, err := proto.Marshal(&protogen.CodeGeneratorResponse{
		SupportedFeatures: &supportedFeatures,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

func newSubCmd(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (*lint.CmdLint, error) {
	data, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, err
	}

	var req protogen.CodeGeneratorRequest
	err = proto.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	flags, err := newFlags(&req)
	if err != nil {
		return nil, err
	}

	subCmd, err := lint.NewCmdLint(
		*flags,
		stdout,
		stderr,
	)
	if err != nil {
		return nil, err
	}
	return subCmd, nil
}

func newFlags(
	req *protogen.CodeGeneratorRequest,
) (*lint.Flags, error) {
	flags, err := lint.NewFlags(req.FileToGenerate)
	if err != nil {
		return nil, err
	}

	var pf subcmds.PluginFlag
	for _, p := range strings.Split(req.GetParameter(), ",") {
		params := strings.SplitN(strings.TrimSpace(p), "=", 2)
		switch params[0] {
		case "":
			continue
		case "config_path":
			if len(params) != 2 {
				return nil, fmt.Errorf("config_path should be specified")
			}
			flags.ConfigPath = params[1]
		case "config_dir_path":
			if len(params) != 2 {
				return nil, fmt.Errorf("config_dir_path should be specified")
			}
			flags.ConfigDirPath = params[1]
		case "fix":
			flags.FixMode = true
		case "reporter":
			if len(params) != 2 {
				return nil, fmt.Errorf("reporter should be specified")
			}
			value := params[1]
			r, err := lint.GetReporter(value)
			if err != nil {
				return nil, err
			}
			flags.Reporter = r
		case "output_file":
			if len(params) != 2 {
				return nil, fmt.Errorf("output_file should be specified")
			}
			flags.OutputFilePath = params[1]
		case "plugin":
			if len(params) != 2 {
				return nil, fmt.Errorf("plugin should be specified")
			}
			err = pf.Set(params[1])
			if err != nil {
				return nil, err
			}
		case "v":
			flags.Verbose = true
		case "proto_root":
			if len(params) != 2 {
				return nil, fmt.Errorf("proto_root should be specified")
			}
			for i, f := range flags.FilePaths {
				flags.FilePaths[i] = filepath.Join(params[1], f)
			}
		default:
			return nil, fmt.Errorf("unmatched parameter: %s", p)
		}
	}

	plugins, err := pf.BuildPlugins(flags.Verbose)
	if err != nil {
		return nil, err
	}
	flags.Plugins = plugins

	return &flags, nil
}

func doVersion(
	stdout io.Writer,
) osutil.ExitCode {
	_, _ = fmt.Fprintln(stdout, "protoc-gen-protolint version "+version+"("+revision+")")
	return osutil.ExitSuccess
}
