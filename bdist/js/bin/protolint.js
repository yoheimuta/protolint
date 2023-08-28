#!/usr/bin/env node

'use strict';

var path = require('path');
var execFile = require('child_process').execFile;

var exe_ext = process.platform === 'win32' ? '.exe' : '';

var protoc = path.resolve(__dirname, 'protolint' + exe_ext);

var args = process.argv.slice(2);

var child_process = execFile(protoc, args, null);

child_process.stdout.pipe(process.stdout);
child_process.stderr.pipe(process.stderr);

child_process.on("exit", (exit_code, _) => {
    process.exit(exit_code);
});

