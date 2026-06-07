#!/usr/bin/env node

const crypto = require("crypto");
const fs = require("fs");
const https = require("https");
const os = require("os");
const path = require("path");

const packageJson = require("../../package.json");

if (process.env.CC_NOTIFY_SKIP_DOWNLOAD === "1") {
  console.log("Skipping cc-notify binary download because CC_NOTIFY_SKIP_DOWNLOAD=1");
  process.exit(0);
}

const root = path.resolve(__dirname, "..", "..");
const distDir = path.join(root, "dist");
const platform = mapPlatform(process.platform);
const arch = mapArch(process.arch);
const exeName = process.platform === "win32" ? "cc-notify.exe" : "cc-notify";
const assetName = `${exeName}_${packageJson.version}_${platform}_${arch}`;
const releaseBase = `https://github.com/ipfred/cc-hooks-notify/releases/download/v${packageJson.version}`;
const binaryURL = `${releaseBase}/${assetName}`;
const checksumsURL = `${releaseBase}/checksums.txt`;
const target = path.join(distDir, exeName);

main().catch((error) => {
  console.error(`Failed to install cc-notify binary: ${error.message}`);
  process.exit(1);
});

async function main() {
  if (!platform || !arch) {
    throw new Error(`unsupported platform ${os.platform()} ${os.arch()}`);
  }
  fs.mkdirSync(distDir, { recursive: true });
  const [binary, checksums] = await Promise.all([
    download(binaryURL),
    download(checksumsURL)
  ]);
  verifyChecksum(binary, checksums.toString("utf8"), assetName);
  fs.writeFileSync(target, binary, { mode: 0o755 });
  if (process.platform !== "win32") {
    fs.chmodSync(target, 0o755);
  }
  console.log(`Installed cc-notify binary: ${target}`);
}

function mapPlatform(value) {
  return {
    darwin: "darwin",
    linux: "linux",
    win32: "windows"
  }[value];
}

function mapArch(value) {
  return {
    x64: "amd64",
    arm64: "arm64"
  }[value];
}

function verifyChecksum(binary, checksums, name) {
  const actual = crypto.createHash("sha256").update(binary).digest("hex");
  const line = checksums.split(/\r?\n/).find((item) => item.trim().endsWith(` ${name}`));
  if (!line) {
    throw new Error(`checksum entry not found for ${name}`);
  }
  const expected = line.trim().split(/\s+/)[0];
  if (actual !== expected) {
    throw new Error(`checksum mismatch for ${name}`);
  }
}

function download(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if ([301, 302, 303, 307, 308].includes(response.statusCode)) {
        response.resume();
        return resolve(download(response.headers.location));
      }
      if (response.statusCode !== 200) {
        response.resume();
        return reject(new Error(`${url} returned HTTP ${response.statusCode}`));
      }
      const chunks = [];
      response.on("data", (chunk) => chunks.push(chunk));
      response.on("end", () => resolve(Buffer.concat(chunks)));
    }).on("error", reject);
  });
}
