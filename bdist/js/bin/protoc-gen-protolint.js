#!/usr/bin/env node

import { spawn } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const extension = process.platform === 'win32' ? '.exe' : '';

const command = path.resolve(__dirname, 'protoc-gen-protolint' + extension);

/** @type {Promise<void>} */
const promise = new Promise((resolve) => {
  const protoc = spawn(command, process.argv.slice(2), {
    stdio: 'inherit', // This inherits stdin, stdout, and stderr
    windowsHide: true,
  });

  protoc.on('close', (code) => {
    process.exitCode = code ?? 0;
    resolve();
  });

  protoc.on('error', (error) => {
    console.error(`Failed to start protoc-gen-protolint: ${error.message}`);
    process.exitCode = 1;
    resolve();
  });
});

await promise;
