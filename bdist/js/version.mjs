import * as child_process from 'child_process';
import * as fs from 'fs';
import * as path from 'path';
import npmlog from 'npmlog';
import semver from 'semver';

import { dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const command = child_process.spawn("git", ["describe", "--tag"], {cwd: __dirname});
var gdt = "";
var done = undefined;
command.stdout.on("data", (data) => { gdt += data; });
command.stderr.on("data", (data) => { npmlog.warn("protolint-version[git]", "git ran into the following error: %s", data) });
command.on("error", (err) => { npmlog.error("protolint-version[git]", "Failed to start git executable: %s",  err) });
command.on("close", (exit_code) => { done = exit_code; });

while (done === undefined) {
    await new Promise(resolve => setTimeout(resolve, 1000));
}

if (done !== 0)
{
    npmlog.error("protolint-version", "Failed to get git tag: %i", done);
    process.exit(done);
}

var version = semver.coerce(gdt);
if (!semver.valid(version)) {
    npmlog.error("protolint-version", "Cannot parse %s to a valid version: %s", gdt, version);
}

npmlog.info("protolint-version", "Preparing to publish %s", version);
const package_json_file = path.join(__dirname, "package.json");
var package_json = JSON.parse(fs.readFileSync(package_json_file));
package_json["version"] = version.version;
fs.writeFileSync(package_json_file, JSON.stringify(package_json, undefined, 2));

npmlog.info("protolint-version", "Successfully written version %s to %s", version, package_json_file);
