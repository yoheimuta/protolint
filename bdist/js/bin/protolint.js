#!/usr/bin/env node

import { spawn } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const extension = process.platform === 'win32' ? '.exe' : '';

const command = path.resolve(__dirname, 'protolint' + extension);

const protolint = spawn(command, process.argv.slice(2), {
  stdio: 'inherit', // This inherits stdin, stdout, and stderr
  windowsHide: true,
});

protolint.on('close', (code) => {
  process.exitCode = code ?? 0;
});

protolint.on('error', (error) => {
  console.error(`Failed to start protolint: ${error.message}`);
  process.exitCode = 1;
});
