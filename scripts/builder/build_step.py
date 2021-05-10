#! /usr/bin/env python3

"""
    Used to build a step
"""

import helpers
import sys
import os


# Step family is the top level directory containing steps of the same vendor
# essentially its steps/<step_family>
def step_family_path(split_step_path):
    family_dir = os.path.join(os.getcwd(), "/".join([".."] * (len(split_step_path) - 1)))
    return os.path.normpath(family_dir)


def main():
    step_rel_path = helpers.get_step_rel_path(os.getcwd())
    split_step_path = [x for x in step_rel_path.split("/") if x != ""]
    helpers.init_logger("_".join(split_step_path))

    image_repo = helpers.get_step_docker_repository(step_rel_path)
    dev_tag = helpers.get_current_branch()
    dev_image_tag = helpers.docker_image_tag(image_repo, dev_tag)

    if not helpers.docker_build(
            dev_image_tag,
        "Dockerfile",
        step_family_path(split_step_path),
        ["--build-arg", "BASE_BRANCH=latest",
         "--build-arg", "CURRENT_BRANCH=" + helpers.get_current_branch(),
         "--build-arg", "STEP_BASEPATH=" + "/".join(split_step_path[1:])]):
        return 1

    for tag in helpers.get_step_image_tags("./"):
        image_tag = helpers.docker_image_tag(image_repo, tag)

        if not helpers.docker_tag(image_repo, dev_tag, tag):
            return 1

        if not helpers.docker_push(image_tag):
            return 1


if __name__ == "__main__":
    sys.exit(main())
