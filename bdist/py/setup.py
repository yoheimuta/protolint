#!/usr/bin/env python3

import contextlib
import hashlib
import logging
import os
import pathlib
import shutil
import subprocess
import sys
import zipfile


@contextlib.contextmanager
def as_cwd(path):
    """Changes working directory and returns to previous on exit."""
    prev_cwd = pathlib.Path.cwd()
    os.chdir(path)
    try:
        yield
    finally:
        os.chdir(prev_cwd)


def clear_dir(path):
    if path.is_dir():
        for e in path.glob("**/*"):
            if e.is_dir():
                logger.debug("Clearing entries from %s", e)
                clear_dir(e)
            else:
                logger.debug("Removing file %s", e)
                e.unlink()
        logger.debug("Removing directory %s", path)
        path.rmdir()


logger = logging.getLogger("BUILD")
logger.setLevel(logging.INFO)
logger.addHandler(logging.StreamHandler(sys.stdout))

file_dir: pathlib.Path = pathlib.Path(os.path.dirname(__file__))
bdist = file_dir / "_bdist"
wheel = file_dir / "_wheel"

if wheel.is_dir():
    clear_dir(wheel)

bdist.mkdir(exist_ok=True, parents=True)
wheel.mkdir(exist_ok=True, parents=True)

logger.info("Building files from %s", file_dir)
repo_root = file_dir / ".." / ".."
license_file = repo_root / "LICENSE"
readme_rst = file_dir / "ReadMe.rst"

logger.info("Using repository root %s", repo_root)

dist = repo_root / "dist"

logger.info("Using previously files from %s", dist)

cp: subprocess.CompletedProcess = subprocess.run(["git", "describe", "--tag"], capture_output=True)
version_id = cp.stdout.decode("utf8").lstrip("v").rstrip("\n")
del cp

logger.info("Assuming version is %s", version_id)

ap_map: dict[str, str] = {
    "darwin_amd64_v1": "macosx_10_0_x86_64",
    "darwin_arm64_v8.0": "macosx_10_0_arm64",
    "linux_amd64_v1": "manylinux2014_x86_64",
    "linux_arm64_v8.0": "manylinux2014_aarch64",
    "linux_arm_7": "manylinux2014_armv7l",
    "windows_amd64_v1": "win_amd64",
    "windows_arm64_v8.0": "win_arm64",
}

executables = {"protolint", "protoc-gen-protolint"}

PY_TAG = "py2.py3"
ABI_TAG = "none"

package_name = "protolint-bin"


for arch_platform in ap_map.keys():
    tag = f"{PY_TAG}-{ABI_TAG}-{ap_map[arch_platform]}"
    logger.info("Packing files for %s using tag %s", arch_platform, tag)
    suffix = ".exe" if "windows" in arch_platform else ""

    pdir = bdist / arch_platform
    clear_dir(pdir)
    pdir.mkdir(exist_ok=True, parents=True)

    p_executables = [dist / f"{exe}_{arch_platform}" / f"{exe}{suffix}" for exe in executables]

    file_name = package_name.replace('-', '_')

    logger.debug("Creating wheel data folder")
    dataFolder = pdir / f"{file_name}-{version_id}.data"
    logger.debug("Creating wheel data folder")
    distInfoFolder = pdir / f"{file_name}-{version_id}.dist-info"

    dataFolder.mkdir(parents=True, exist_ok=True)
    distInfoFolder.mkdir(parents=True, exist_ok=True)

    with as_cwd(pdir):
        logger.debug("Creating scripts folder")
        scripts = dataFolder / "scripts"
        scripts.mkdir(parents=True, exist_ok=True)
        for p in p_executables:
            logger.debug("Copying executable %s to scripts folder %s", p, scripts)
            shutil.copy(p, scripts)

        logger.debug("Copying LICENSE from %s", license_file)
        shutil.copy(license_file, distInfoFolder)

        with (distInfoFolder / "WHEEL").open("w+", encoding="utf-8") as wl:
            logger.debug("Writing WHEEL file")
            wl.writelines([
                "Wheel-Version: 1.0\n",
                "Generator: https://github.com/yoheimuta/protolint/\n",
                "Root-Is-PureLib: false\n",
                f"Tag: {tag}\n"]
            )

        with (distInfoFolder / "METADATA").open("w+", encoding="utf-8") as ml:
            logger.debug("Writing METADATA file")
            ml.writelines([
                "Metadata-Version: 2.1\n",
                f"Name: {package_name}\n",
                f"Version: {version_id} \n",
                "Summary: A pluggable linter and fixer to enforce Protocol Buffer style and conventions.\nThis package contains the pre-compiled binaries.\n",
                "Home-page: https://github.com/yoheimuta/protolint/\n",
                "Author: yohei yoshimuta\n",
                "Maintainer: yohei yoshimuta\n",
                "License: MIT\n",
                "Project-URL: Official Website, https://github.com/yoheimuta/protolint/\n",
                "Project-URL: Source Code, https://github.com/yoheimuta/protolint.git\n",
                "Project-URL: Issue Tracker, https://github.com/yoheimuta/protolint/issues\n",
                "Classifier: Development Status :: 5 - Production/Stable\n",
                "Classifier: Environment :: Console\n",
                "Classifier: Intended Audience :: Developers\n",
                "Classifier: License :: OSI Approved :: MIT License\n",
                "Classifier: Natural Language :: English\n",
                "Classifier: Operating System :: MacOS\n",
                "Classifier: Operating System :: Microsoft :: Windows\n",
                "Classifier: Operating System :: POSIX :: Linux\n",
                "Classifier: Programming Language :: Go\n",
                "Classifier: Topic :: Software Development :: Pre-processors\n",
                "Classifier: Topic :: Utilities\n", 
                "Requires-Python: >= 3.0\n",
                "Description-Content-Type: text/rst\n",
                "License-File: LICENSE\n",
            ])

            ml.writelines(["\n"])

            with readme_rst.open("r", encoding="utf-8") as readme:
                ml.writelines(readme.readlines())

        wheel_content = list(distInfoFolder.glob("**/*")) + list(dataFolder.glob("**/*"))
        elements_to_relative_paths = {entry: str(entry).lstrip(str(pdir)).lstrip("/").lstrip("\\") for entry in wheel_content if entry.is_file()}
        with (distInfoFolder / "RECORD").open("w+", encoding="utf-8") as rl:
            logger.debug("Writing RECORD file")
            for entry in elements_to_relative_paths.keys():
                relPath = elements_to_relative_paths[entry]
                sha256 = hashlib.sha256(entry.read_bytes())
                fs = entry.stat().st_size
                rl.write(f"{relPath},sha256={sha256.hexdigest()},{str(fs)}\n")

            rl.write(distInfoFolder.name + "/RECORD,,\n")
            wheel_content.append(distInfoFolder / "RECORD")

        whl_file = wheel / f"{file_name}-{version_id}-{tag}.whl"
        if whl_file.is_file():
            logger.debug("Removing existing wheel file")
            whl_file.unlink()

        with zipfile.ZipFile(whl_file, "w", compression=zipfile.ZIP_DEFLATED) as whl:
            logger.info("Creating python wheel %s", whl_file)
            for content in wheel_content:
                whl.write(
                    content,
                    content.relative_to(pdir),
                    zipfile.ZIP_DEFLATED,
                )

logger.info("Done")
