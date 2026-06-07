#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");

const root = path.resolve(__dirname, "..", "..");
const exeName = process.platform === "win32" ? "cc-notify.exe" : "cc-notify";
const binary = path.join(root, "dist", exeName);

if (!fs.existsSync(binary)) {
  console.error(`cc-notify binary not found: ${binary}`);
  console.error("Try reinstalling @fredzhang/cc-notify.");
  process.exit(1);
}

const result = spawnSync(binary, process.argv.slice(2), {
  stdio: "inherit",
  windowsHide: true
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status === null ? 1 : result.status);
