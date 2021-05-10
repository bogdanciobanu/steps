#! /usr/bin/env python3
"""
    This script is used to build ./base image
"""

import helpers
import sys

STEP_PATH = "base"


def main():
    helpers.init_logger(STEP_PATH)
    image_repo = helpers.get_step_docker_repository(STEP_PATH)
    dev_tag = helpers.get_current_branch()
    dev_image_tag = helpers.docker_image_tag(image_repo, dev_tag)

    if not helpers.docker_build(dev_image_tag):
        return 1

    for tag in helpers.get_step_image_tags("./"):
        image_tag = helpers.docker_image_tag(image_repo, tag)

        if not helpers.docker_tag(image_repo, dev_tag, tag):
            return 1

        if not helpers.docker_push(image_tag):
            return 1


if __name__ == "__main__":
    sys.exit(main())
