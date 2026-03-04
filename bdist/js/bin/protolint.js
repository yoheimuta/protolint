#!/usr/bin/env node

import { spawn } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const extension = process.platform === 'win32' ? '.exe' : '';

const command = path.resolve(__dirname, 'protolint' + extension);

/** @type {Promise<void>} */
const promise = new Promise((resolve) => {
  const protolint = spawn(command, process.argv.slice(2), {
    stdio: 'inherit',
    windowsHide: true,
  });

  protolint.on('close', (code) => {
    process.exitCode = code ?? 0;
    resolve();
  });

  protolint.on('error', (error) => {
    console.error(`Failed to start protolint: ${error.message}`);
    process.exitCode = 1;
    resolve();
  });
});

await promise;
