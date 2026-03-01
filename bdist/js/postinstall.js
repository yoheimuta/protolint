import { mkdir } from 'node:fs/promises';
import path from 'node:path';
import { pipeline } from 'node:stream/promises';
import { fileURLToPath } from 'node:url';
import { createGunzip } from 'node:zlib';

import { unpackTar } from 'modern-tar/fs';
import { EnvHttpProxyAgent, fetch, Headers, ProxyAgent } from 'undici';

const SCRIPT_NAME = 'protolint-install';

const baseUrl = process.env.PROTOLINT_MIRROR_HOST ?? 'https://github.com';
const remotePath =
  process.env.PROTOLINT_MIRROR_REMOTE_PATH ??
  `yoheimuta/protolint/releases/download`;

const version = process.env.npm_package_version;

if (version === undefined) {
  throw new Error(
    `${SCRIPT_NAME}: Failed getting the package version from env`,
  );
}

const platform = process.platform === 'win32' ? 'windows' : process.platform;
const arch = process.arch === 'x64' ? 'amd64' : process.arch;

const urlPath = `/${remotePath}/v${version}/protolint_${version}_${platform}_${arch}.tar.gz`;

const downloadUrl = new URL(urlPath, baseUrl);

let dispatcher;

if (!process.env.PROTOLINT_NO_PROXY) {
  dispatcher = process.env.PROTOLINT_PROXY
    ? new ProxyAgent(process.env.PROTOLINT_PROXY)
    : new EnvHttpProxyAgent();
}

const headers = new Headers();
const username = process.env.PROTOLINT_MIRROR_USERNAME;
const password = process.env.PROTOLINT_MIRROR_PASSWORD;

if (username && password) {
  const auth = Buffer.from(`${username}:${password}`).toString('base64');
  headers.set('Authorization', `Basic ${auth}`);
}

console.info(
  `${SCRIPT_NAME}: Fetching protolint executable from ${downloadUrl.href}`,
);

const response = await fetch(downloadUrl, {
  dispatcher,
  headers,
  redirect: 'follow',
});

if (!response.ok || !response.body) {
  throw new Error(
    `${SCRIPT_NAME}: Failed downloading the archive: ${response.statusText}`,
  );
}

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const binDirectory = path.join(__dirname, 'bin');

await mkdir(binDirectory, { recursive: true });

const executables = new Set(
  ['protolint', 'protoc-gen-protolint'].map((basename) =>
    process.platform === 'win32' ? `${basename}.exe` : basename,
  ),
);

await pipeline(
  response.body,

  createGunzip(),

  unpackTar(binDirectory, {
    filter: ({ name }) => executables.has(name),
    map: (header) => {
      header.mode = 0o755;

      return header;
    },
  }),
);

console.info(`${SCRIPT_NAME}: Protolint installed successfully.`);
