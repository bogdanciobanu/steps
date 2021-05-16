#!/usr/bin/env python3

import glob
import sys
import yaml
import hashlib
import datetime
import collections
import os
import re
from datetime import datetime, timezone


_mapping_tag = yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG


def dict_representer(dumper, data):
    return dumper.represent_dict(data.items())


def dict_constructor(loader, node):
    return collections.OrderedDict(loader.construct_pairs(node))


yaml.add_representer(collections.OrderedDict, dict_representer)
yaml.add_constructor(_mapping_tag, dict_constructor)


def filetime(path):
    local_time = datetime.fromtimestamp(os.path.getmtime(path), timezone.utc)
    return local_time.isoformat()


def gettime():
    local_time = datetime.now(timezone.utc).astimezone()
    return local_time.isoformat()


def main(inputdir="./out/manifests/", outputfile="./out/manifests/index.yml"):
    files = glob.glob("./out/manifests/*.yml")
    if outputfile in files:
        files.remove(outputfile)

    result = collections.OrderedDict(
        {"apiVersion": "stackpulse.io/v1", "kind": "Index", "entries": []}
    )
    all_digest = hashlib.sha256()
    for filename in files:
        data = open(filename, "rb").read()
        digest = hashlib.sha256(data).hexdigest()
        name = os.path.basename(filename)
        time = filetime(filename)
        result["entries"].append(
            collections.OrderedDict({"name": name, "digest": digest, "created": time})
        )
        all_digest.update(data)

    all_digest = all_digest.hexdigest()
    time = gettime()
    result["digest"] = all_digest
    result["created"] = time

    with open(outputfile, "wb") as f:
        f.write(yaml.dump(result).encode())
    print('[+] Generated "%s" hash: %s' % (outputfile, all_digest))


if __name__ == "__main__":
    main(*sys.argv[1:])
