#!/usr/bin/env node
/* eslint-disable no-console */
// Postinstall: download the prebuilt helios-cli binary matching the current
// platform/arch from GitHub Releases and place it at npm/bin/helios-cli(.exe).
//
// Asset naming convention (produced by `goreleaser` or your release CI):
//   helios-cli_<os>_<arch>.tar.gz   (linux/darwin)
//   helios-cli_<os>_<arch>.zip      (windows)

'use strict';

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execFileSync } = require('child_process');

const pkg = require('../package.json');

const PLATFORM_MAP = {
  'darwin-arm64': { os: 'darwin', arch: 'arm64', ext: 'tar.gz' },
  'darwin-x64': { os: 'darwin', arch: 'amd64', ext: 'tar.gz' },
  'linux-x64': { os: 'linux', arch: 'amd64', ext: 'tar.gz' },
  'linux-arm64': { os: 'linux', arch: 'arm64', ext: 'tar.gz' },
  'win32-x64': { os: 'windows', arch: 'amd64', ext: 'zip' },
};

const BINARY_DIR = path.join(__dirname, 'bin');
const BINARY_NAME = process.platform === 'win32' ? 'helios-cli.exe' : 'helios-cli';
const BINARY_PATH = path.join(BINARY_DIR, BINARY_NAME);

function log(msg) {
  console.log(`[helios-cli] ${msg}`);
}

function bail(msg) {
  console.error(`[helios-cli] ${msg}`);
  process.exit(1);
}

function detectAsset() {
  const key = `${process.platform}-${process.arch}`;
  const meta = PLATFORM_MAP[key];
  if (!meta) {
    bail(
      `unsupported platform: ${key}. ` +
        'Build from source instead: see https://github.com/auXiaoYuan/helios-cli',
    );
  }
  const version = (pkg.helios && pkg.helios.binaryVersion) || pkg.version;
  const base =
    (pkg.helios && pkg.helios.releaseBase) ||
    'https://github.com/auXiaoYuan/helios-cli/releases/download';
  const filename = `helios-cli_${meta.os}_${meta.arch}.${meta.ext}`;
  return {
    url: `${base}/v${version}/${filename}`,
    filename,
    ext: meta.ext,
  };
}

function download(url, destPath, redirects) {
  if (redirects === undefined) redirects = 0;
  return new Promise((resolve, reject) => {
    const req = https.get(url, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        if (redirects > 5) {
          reject(new Error('too many redirects'));
          return;
        }
        res.resume();
        resolve(download(res.headers.location, destPath, redirects + 1));
        return;
      }
      if (res.statusCode !== 200) {
        reject(new Error(`download failed: HTTP ${res.statusCode} for ${url}`));
        res.resume();
        return;
      }
      const out = fs.createWriteStream(destPath);
      res.pipe(out);
      out.on('finish', () => out.close(resolve));
      out.on('error', reject);
    });
    req.on('error', reject);
  });
}

function extract(archivePath, ext, outDir) {
  if (ext === 'tar.gz') {
    // Use system tar — available on macOS/Linux/Windows 10+.
    execFileSync('tar', ['-xzf', archivePath, '-C', outDir], { stdio: 'inherit' });
    return;
  }
  if (ext === 'zip') {
    if (process.platform === 'win32') {
      execFileSync(
        'powershell',
        [
          '-NoProfile',
          '-Command',
          `Expand-Archive -Force -Path '${archivePath}' -DestinationPath '${outDir}'`,
        ],
        { stdio: 'inherit' },
      );
    } else {
      execFileSync('unzip', ['-o', archivePath, '-d', outDir], { stdio: 'inherit' });
    }
    return;
  }
  throw new Error(`unknown archive ext: ${ext}`);
}

async function main() {
  if (process.env.HELIOS_CLI_SKIP_DOWNLOAD === '1') {
    log('HELIOS_CLI_SKIP_DOWNLOAD=1 set, skipping binary download.');
    return;
  }

  fs.mkdirSync(BINARY_DIR, { recursive: true });

  if (fs.existsSync(BINARY_PATH)) {
    log(`binary already present at ${BINARY_PATH}, skipping.`);
    return;
  }

  const asset = detectAsset();
  const archivePath = path.join(BINARY_DIR, asset.filename);

  log(`downloading ${asset.url}`);
  try {
    await download(asset.url, archivePath);
  } catch (err) {
    bail(`failed to download binary: ${err.message}`);
  }

  log(`extracting ${asset.filename}`);
  try {
    extract(archivePath, asset.ext, BINARY_DIR);
  } catch (err) {
    bail(`failed to extract binary: ${err.message}`);
  }

  try {
    fs.unlinkSync(archivePath);
  } catch (_) {
    // best effort cleanup
  }

  if (!fs.existsSync(BINARY_PATH)) {
    bail(`binary ${BINARY_NAME} not found after extraction in ${BINARY_DIR}`);
  }

  if (process.platform !== 'win32') {
    fs.chmodSync(BINARY_PATH, 0o755);
  }

  log(`installed binary at ${BINARY_PATH}`);
}

main().catch((err) => bail(err.stack || String(err)));
