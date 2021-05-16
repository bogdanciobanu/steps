#! /usr/bin/env python3

import steps
import sys
import os
import builder.helpers
import traceback
import argparse
import logging

DOCKERFILE = "Dockerfile"

TOML_TEMPLATE = """
name = "{step_name}"
includes = ["{{{{ .Root }}}}/baur-includes/step.toml#build-step"]
"""
APP_FILE = ".app.toml"


# Check app name is set correctly in .app.toml
def check_app_toml(step):
    toml_path = get_app_toml_file(step)
    if not os.path.isfile(toml_path):
        logging.error(f"{toml_path}: Step toml file not found")
        return False

    f = open(toml_path, "r").read()
    expected_string = f"name = \"{get_step_toml_name(step)}\""
    if expected_string not in f:
        logging.error(f"{toml_path} Step name mismatch: {expected_string} not found")
        return False
    
    return True


def create_app_toml(step):
    toml_path = get_app_toml_file(step)

    try:
        open(toml_path, "w").write(TOML_TEMPLATE.format(**{"step_name": get_step_toml_name(step)}))
    except:
        traceback.print_exc()
        logging.error(f"Error: failed creating {toml_path}")
        return 1

    logging.info(f"Created {toml_path}")


def get_app_toml_file(step):
    return os.path.join(step, APP_FILE)


def get_step_toml_name(step):
    return builder.helpers.get_step_rel_path(step)


# Main entry point
def main():
    logging.basicConfig(format='%(message)s', level=logging.INFO)

    parser = argparse.ArgumentParser(description='Baur app help tool for StackPulse steps.')
    parser.add_argument('--init', help=f"create {APP_FILE} in steps which don't already have it", action='store_true')
    parser.add_argument('--check', help=f"check that all steps contain a valid {APP_FILE}", action='store_true')
    args = parser.parse_args()

    if not args.init:
        args.check = True

    missing_apps = []
    total_steps = 0
    for step in steps.get_steps():
        if not os.path.isfile(os.path.join(step, DOCKERFILE)):
            continue

        total_steps += 1

        logging.debug(f"Checking step: {step}")
        if check_app_toml(step):
            continue

        missing_apps.append(step)

    if len(missing_apps) > 0:
        logging.info(f"{total_steps} steps found / {len(missing_apps)} steps missing or invalid {APP_FILE}")
        for a in missing_apps:
            if args.init:
                create_app_toml(a)
        if args.check:
            return 1


if __name__ == "__main__":
    sys.exit(main())
