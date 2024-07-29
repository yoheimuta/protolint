# protolint

Protolint is a go based linter for `.proto` files for google protobuf und gRPC. It follows the best practices at the [protobuf.dev](https://protobuf.dev/programming-guides/style/) (and ff.). Please note, that this should be a dev-dependency.

The npm package provides a wrapper around the executables `protolint` and `protoc-gen-protolint`. During installation process, it will download the binaries matching the version and your operating system and CPU architecture from github.

If your behind a proxy, you can add the `PROTOLINT_PROXY` environment variable including the HTTP basic authentication information like username and password. **NOTE** that this will take precedence of the system `HTTP_PROXY`/`HTTPS_PROXY` environment variables. If these variables should be used, do not use `PROTOLINT_PROXY`.

If your running an airgapped environment, you can add the following environment variables:

`PROTOLINT_MIRROR_HOST`: The basic url you are using to serve the binaries. Defaults to `https://github.com`
`PROTOLINT_MIRROR_REMOTE_PATH`: The relative path on the mirror host. Defaults to `yoheimuta/protolint/releases/download/`

Within the remote path, make sure, that a folder `v<version>` exists containing the files downloaded from the github releases.

If you are required to authenticate against your mirror, use the following environment variables:

`PROTOLINT_MIRROR_USERNAME`: The user name. Defaults to an empty string.
`PROTOLINT_MIRROR_PASSWORD`: The password or identifaction token. Defaults to an empty string.

For node based projects, you can add the protobuf configuration to your `package.json` using a node called `protolint`.

For more information about protolint, its parameters and command-line arguments refer to the original ReadMe in the [github repository](https://github.com/yoheimuta/protolint).
