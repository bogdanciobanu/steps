#! /usr/bin/env python

import steps
import sys
import os
import builder.helpers
import traceback
import argparse

DOCKERFILE = "Dockerfile"

TOML_TEMPLATE = """
name = "{step_name}"
includes = ["{{{{ .Root }}}}/baur-includes/step.toml#build-step"]
"""
APP_FILE = ".app.toml"

def main():

    parser = argparse.ArgumentParser(description='Initialize baur app in step directories.')
    parser.add_argument('quiet', help='dont prompt to overwrite app.toml if exists', action='store_true')
    args = parser.parse_args()

    for step in steps.get_steps():
        if not os.path.isfile(os.path.join(step, DOCKERFILE)):
            continue

        print("> Step: %s" % (step,))
        app_toml = os.path.join(step, APP_FILE)
        if os.path.isfile(app_toml):
            if args.quiet:
                print("[!] app.toml already exists")
                continue

            if "y" != input("[!] app.toml already exists, overwrite? [y/N]"):
                continue

        try:
            open(app_toml, "w").write(TOML_TEMPLATE.format(**{"step_name": builder.helpers.get_step_rel_path(step)}))
        except:
            traceback.print_exc()
            print("Error: failed creating .app.toml")
            return 1

        print("  Created %s" % (APP_FILE,))


if __name__ == "__main__":
    sys.exit(main())