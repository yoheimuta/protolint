#!/usr/bin/env node

'use strict';

var path = require('path');
var spawn = require('child_process').spawn;

var exe_ext = process.platform === 'win32' ? '.exe' : '';

var protoc = path.resolve(__dirname, 'protoc-gen-protolint' + exe_ext);

var args = process.argv.slice(2);

var child_process = spawn(protoc, args, {
  stdio: 'inherit' // This inherits stdin, stdout, and stderr
});

child_process.on("exit", (exit_code, _) => {
    process.exit(exit_code);
});
