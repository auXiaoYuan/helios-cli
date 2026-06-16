#!/usr/bin/env node
// Thin shim: spawn the bundled helios-cli binary with the user's argv.

'use strict';

const path = require('path');
const fs = require('fs');
const { spawn } = require('child_process');

const binaryName = process.platform === 'win32' ? 'helios-cli.exe' : 'helios-cli';
const binaryPath = path.join(__dirname, 'bin', binaryName);

if (!fs.existsSync(binaryPath)) {
  console.error(
    `[helios-cli] binary not found at ${binaryPath}. ` +
      'Try reinstalling, or run `node npm/install.js` to redownload.',
  );
  process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit' });

child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code === null ? 1 : code);
});

child.on('error', (err) => {
  console.error(`[helios-cli] failed to start binary: ${err.message}`);
  process.exit(1);
});
