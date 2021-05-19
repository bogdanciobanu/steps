#! /usr/bin/env python3

"""
    Used to build a step
"""

import helpers
import os
import constants
import logging
import glob


class LintException(Exception):
    """ exception that happens due to linting failure """
    pass


class TestException(Exception):
    """ exception caused by testing failure """
    pass


class BuildException(Exception):
    """ exception caused by failed build """


# Current step image repository
def image_repo() -> str:
    return helpers.get_step_docker_repository(helpers.get_step_rel_path(os.getcwd()))


# Current step's dev tag (current branch)
def dev_tag() -> str:
    return helpers.get_current_branch()


# Current steps docker image with tag
def dev_image_tag() -> str:
    return helpers.docker_image_tag(image_repo(), dev_tag())


# Get current step path in filename format `redis/get` -> `redis_get`
def step_path_as_filename() -> str:
    step_rel_path = helpers.get_step_rel_path(os.getcwd())
    return step_rel_path.replace("/", "_")


# Get the base directory for tests results, defined by `TEST_OUTPTU_DIR` env var
def tests_output_path() -> str:
    test_output_dir = os.getenv("TEST_OUTPUT_DIR", "./")
    return os.path.join(test_output_dir, step_path_as_filename())


# Step family is the top level directory containing steps of the same vendor
# essentially its steps/<step_family>
def step_family_path() -> str:
    split_step_path = helpers.get_step_rel_path(os.getcwd()).split("/")
    family_dir = os.path.join(os.getcwd(), "/".join([".."] * (len(split_step_path) - 1)))
    return os.path.normpath(family_dir)


# Validate current directory is a valid baur application
def assert_step_cwd():
    assert os.path.isfile(constants.APP_TOML), f"must in directory with {constants.APP_TOML} file"


# Run integration tests; Integration tests run on docker image and have a build tag `integration`
def run_integration_tests():
    image = dev_image_tag()
    if len(glob.glob("./**integration_test.go")) == 0:
        logging.info("No integration tests found for step")
        return

    test_results_path = tests_output_path() + ".integration.junit.xml"
    logging.info(f"Running tests for {image}, output={test_results_path}")
    if not helpers.run_command(
            ["gotestsum", "-f", "standard-verbose", "--junitfile", test_results_path, "--", "--tags=integration"],
            env={"STEP_IMAGE": image}):
        raise TestException("Integration tests failed")


# Run unit tests
def run_unit_tests():
    unit_tests = [x for x in glob.glob("./**_test.go") if not "_integration" in x]
    if len(unit_tests) == 0:
        logging.info("No unit tests found for step")
        return

    test_results_path = tests_output_path() + ".unit.junit.xml"
    logging.info(f"Running unit tests, output={test_results_path}")
    if not helpers.run_command(["gotestsum", "-f", "standard-verbose", "--junitfile", test_results_path]):
        raise TestException("Unit tests failed")


# Check go code linting
def check_go_linting():
    go_files = glob.glob("./**.go")
    if len(go_files) == 0:
        logging.info("No Go source files found")
        return

    logging.info("Checking go-fmt")
    imports_result, result = helpers.run_command_with_output(["goimports", "-l", *go_files])
    if imports_result is None:
        return

    if len(imports_result) > 0 or not result:
        raise LintException(f"The following go files are not formatted: {imports_result}")

    logging.info("Running golangci-lint")
    if not helpers.run_command(["golangci-lint", "run", f"--path-prefix={os.getcwd()}"]):
        raise LintException("golangci-lint failed")


# Code executed before the step's docker image is built
def pre_image_build():
    # check code formatting
    check_go_linting()

    run_unit_tests()


# Build step image
def build_step_image():
    split_step_path = helpers.get_step_rel_path(os.getcwd()).split("/")
    if not helpers.docker_build(
            dev_image_tag(),
            "Dockerfile",
            step_family_path(),
            ["--build-arg", "BASE_BRANCH=latest",
             "--build-arg", "CURRENT_BRANCH=" + helpers.get_current_branch(),
             "--build-arg", "STEP_BASEPATH=" + "/".join(split_step_path[1:])]):
        raise BuildException("Failed building step image")


# Stuff that runs after the dev image has been built
def post_image_build():
    # run integration tests
    run_integration_tests()

    for tag in helpers.get_step_image_tags("./"):
        target_image_tag = helpers.docker_image_tag(image_repo(), tag)

        if not helpers.docker_tag(dev_image_tag(), target_image_tag):
            raise BuildException(f"docker tag failed {target_image_tag}")

        if not helpers.docker_push(target_image_tag):
            raise BuildException(f"docker push failed {target_image_tag}")


def build_base_image():
    if not helpers.docker_build(dev_image_tag()):
        raise BuildException("failed building base image")


# Assumes the script is executed by baur from step root directory (e.g: `steps/redis/get`)
# Build process:
# 1. pre_image_build()
# 2. Build image and tag with `dev_image_tag()`
# 3. post_image_build()
#
def build(base_image_build: bool = False) -> bool:
    helpers.init_logger(step_path_as_filename())

    assert_step_cwd()

    # Run pre image build steps
    pre_image_build()

    # check if there is a docker file in this directory
    if not os.path.isfile("Dockerfile"):
        logging.info("No Dockerfile in this step, skipping build")
        return True

    if not base_image_build:
        build_step_image()
    else:
        build_base_image()

    # run post image build steps
    post_image_build()

    return True
