#! /usr/bin/env python3

import glob
import yaml
import traceback
import sys
import os.path
import shutil
import builder.constants as constants

STEPS_PREFIX = "./steps/"


def get_image_name(path, image_prefix):
    path = path.replace(STEPS_PREFIX, "")
    return os.path.join(image_prefix, os.path.dirname(path))


def get_manifest_filename(path):
    return os.path.dirname(path).replace(STEPS_PREFIX, "").replace("/", "-") + constants.MANIFEST_SUFFIX


def patch_image_name(manifest_path, image_name):
    f = open(manifest_path)
    try:
        y = yaml.safe_load(f)
    except:
        traceback.print_exc()
        raise Exception(f"[!] Error parsing yaml {manifest_path}")

    if not 'metadata' in y:
        raise Exception(f"[!] No metadata object found in yaml {manifest_path}")

    if 'imageName' in y['metadata']:
        print(f"> imageName already set: {y['metadata']['imageName']}")
        return

    print(f"> Setting image name to: {image_name}")
    y['metadata']['imageName'] = image_name
    yaml.dump(y, open(manifest_path, "w"), default_flow_style=False)
    print("[+] Saved")


def main():
    try:
        prog_name, manifest_dir, output_dir = sys.argv
    except:
        print("Usage: prepare-manifests.py <manifest_dir> <output_dir>")
        return 1

    try:
        os.makedirs(output_dir)
    except:
        pass

    manifests = glob.glob(STEPS_PREFIX + "**/" + constants.MANIFEST_FILENAME, recursive=True)

    print(f"Manifests: {len(manifests)}")

    for m in manifests:
        print(f"Processing {m}")

        target_filename = os.path.join(output_dir, get_manifest_filename(m))
        image_name = get_image_name(m, constants.CONTAINER_REGISTRY)

        # Copy manifest
        print(f"> Copying to {target_filename}")

        shutil.copy(m, target_filename)
        patch_image_name(target_filename, image_name)

    return 0

if __name__ == "__main__":
    sys.exit(main())
