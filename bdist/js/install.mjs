import { got } from 'got';
import { createFetch } from 'got-fetch';
import npmlog from 'npmlog';
import { HttpProxyAgent } from 'http-proxy-agent';
import * as fs from 'fs';
import { temporaryFile } from 'tempy';
import * as tar from 'tar';
import * as pipeline from 'stream';

const _arch_mapping = { "x64": "amd64" };
const _platform_mapping = { "win32": "windows" };

const _platform = process.platform;
const _arch = process.arch;

// TODO FIND correct paths -> "./bin" goes to node_modules

const script_name = "protolint-install";

const module_name = "protolint";
const protolint_host = process.env.PROTOLINT_MIRROR_HOST ?? "https://github.com";
const protolint_path = process.env.PROTOLINT_MIRROR_REMOTE_PATH ?? `yoheimuta/${module_name}/releases/download/`;
const protolint_version = process.env.npm_package_version;
const platform = _platform_mapping[_platform] ?? _platform;
const arch = _arch_mapping[_arch] ?? _arch;

const url = `${protolint_host}/${protolint_path}/v${protolint_version}/${module_name}_${protolint_version}_${platform}_${arch}.tar.gz`;

let agent;

let proxy_address = process.env.PROTOLINT_PROXY;
if (!proxy_address)
{
    proxy_address = protolint_host.startsWith("https") ? process.env.HTTPS_PROXY : process.env.HTTP_PROXY;
}

if (proxy_address) {
    agent = new HttpProxyAgent(proxy_address);
}

const agent_config = {
    http: agent
};

const got_config = {
    followRedirect: true,
    maxRedirects: 3,
    username: process.env.PROTOLINT_MIRROR_USERNAME ?? '',
    password: process.env.PROTOLINT_MIRROR_PASSWORD ?? '',
    agent: agent_config,
};

const instance = got.extend(got_config);

function get_filename_with_extension(fileName) {
    const ext = process.platform == "win32" ? ".exe" : "";
    return `${fileName}${ext}`;
}

npmlog.info(script_name, "Fetching protolint executable from %s", url);

const fetch = createFetch(instance);

fetch(url).then(
    async response => {
        if (response.ok)
        {
            const targetFile = temporaryFile({ name: "_protolint.tar.gz"});
            const out = fs.createWriteStream(targetFile, {
                flags: "w+"
            });
            var success = undefined;
            const streaming = pipeline.pipeline(
                response.body,
                out,
                (err) => {
                    if (err)
                    {
                        npmlog.error(script_name, "Failed to save downloaded file: %s", err);
                        success = false;
                    }
                    else
                    {
                        npmlog.info(script_name, "Protolint saved to %s", targetFile);
                        success = true;
                    }
                }
            );

            while (success === undefined)
            {
                await new Promise(resolve => setTimeout(resolve, 1000));
            }

            if (success)
            {
                return targetFile;
            }

            return null;
        }
        else
        {
            npmlog.error(script_name, "Failed to download %s. Got status: %i", response.url, response.status);
            return null;
        }
    }
).then(
    previous => {
        if (!fs.existsSync("./bin"))
        {
            fs.mkdirSync("./bin");
        }

        return previous;
    }
).then(
    async file => {
        if (file)
        {
            const result = await tar.x(
                {
                    "keep-existing": false,
                    cwd: "./bin/",
                    sync: false,
                    file: file,
                    strict: true,
                },
                [
                    get_filename_with_extension("protolint"),
                    get_filename_with_extension("protoc-gen-protolint"),
                ],
                (err) => {
                    if (err) {
                        npmlog.error(script_name, "Failed to extract protlint executables: %s");
                    }
                },
            )
            .then(
                () => {
                    return {
                       protolint: `./bin/${get_filename_with_extension("protolint")}`, 
                       protoc_gen_protolint: `./bin/${get_filename_with_extension("protoc-gen-protolint")}`,
                    };
                }
            ).catch(
                (err) => {
                    npmlog.error(script_name, "Failed to extract files from downloaded tar file: %s", err);
                    return {
                       protolint: undefined, 
                       protoc_gen_protolint: undefined,
                    };
                }
            );

            return result;
        }
        else
        {
            npmlog.warn(script_name, "Could not find downloaded protolint archive.");
            return {
                protolint: undefined, 
                protoc_gen_protolint: undefined,
             };
        }
    }
).then(
    (protolint_obj) => {
        return (protolint_obj != null && protolint_obj.protolint != null && protolint_obj.protoc_gen_protolint != null);
    }
).then(
    (result) => {
        if (result){
            npmlog.info(script_name, "Protolint installed successfully.");
        }
        else {
            npmlog.warn(script_name, "Failed to download protolint. See previous messages for details");
        }
    }
).catch(
    reason => {
        npmlog.error(script_name, "Failed to install protolint: %s", reason);
        process.exit(1);
    }
);