import constants
import os
import subprocess
import yaml
import logging
import sys


# Will initialize logger for given step.
# If env var LOG_OUPUT_DIR is set it will also store log files to that directory
def init_logger(step_name: str = None) -> None:
    stdout_handler = logging.StreamHandler(sys.stdout)
    handlers = [stdout_handler]

    log_output_dir = os.getenv("LOG_OUTPUT_DIR")
    if log_output_dir is not None and step_name is not None:
        handlers.append(logging.FileHandler(filename=os.path.join(log_output_dir, f"{step_name}.log")))

    logging.basicConfig(
        level=logging.DEBUG,
        format="[%(asctime)s] %(levelname)s - %(message)s",
        handlers=handlers
    )


# Run command and return true/false depending on exit code
def run_command(args: list[str], env: dict[str, str] = {}) -> bool:
    output, result = run_command_with_output(args, env)
    return result


# Run command and return stdout
def run_command_with_output(args: list[str], env: dict[str, str] = {}) -> bool:
    cmdline = " ".join(args)
    logging.info(f"> Running: {cmdline}")
    output = None
    try:
        output = subprocess.check_output(args, env=dict(os.environ, **env))
    except subprocess.CalledProcessError as e:
        if e.output:
            logging.error(str(e.output, 'utf-8'))
        logging.exception("Failed running command")
        return None, False

    if len(output) != 0:
        output = str(output, 'utf-8')
        logging.info(output)
    return output, True


# Get current branch using git cli tool
def get_current_branch() -> str:
    branch = subprocess.check_output(["git", "rev-parse", "--abbrev-ref", "HEAD"]).strip().decode("utf-8")
    logging.info(f"Current branch: {branch}")
    return branch


# Get step relative path (to /steps/ directory)
def get_step_rel_path(path: str) -> str:
    lookup = constants.STEPS_ROOT + "/"
    return path[path.rfind(lookup) + len(lookup):]


# Get step docker image path
def get_step_docker_repository(step_path: str) -> str:
    return os.path.join(constants.CONTAINER_REGISTRY, step_path)


# Try to read versions.step / versions file
def get_version_from_file(step_path: str) -> str:
    for fn in constants.VERSION_FILENAMES:
        fn = os.path.join(step_path, fn)
        if os.path.isfile(fn):
            return open(fn).read().strip()

    return None


# Get manifest version, first it will try either VERSION or VERSION.txt files then manifest.yaml
def get_manifest_version(step_path: str) -> str:
    # Check for version file
    version = get_version_from_file(step_path)
    if version is not None:
        logging.info(f"Step version by version file: {version}")
        return version

    manifest_path = os.path.join(step_path, constants.MANIFEST_FILENAME)

    if not os.path.exists(manifest_path):
        raise Exception("no manifest found")

    try:
        f = open(manifest_path)
        y = yaml.safe_load(f)
        f.close()
        version = y['metadata']['version']
        logging.info(f"Step version from manifest: {version}")
        return version
    except:
        logging.error(f"[X] Failed parsing version from {manifest_path}")
        raise Exception("failed parsing manifest")

    if version == "":
        raise Exception(f"no valid version set for step: {manifest_path}")


# Build full docker image path given a repo and tag
def docker_image_tag(repo: str, tag: str) -> str:
    return f"{repo}:{tag}"


# return the relevant tags for the given branch
# on side branch we tag with current branch on master we tag with latest and imageVersion
def get_step_image_tags(step_path: str) -> list[str]:
    current_branch = get_current_branch()
    if current_branch == "master":
        return ["latest", get_manifest_version(step_path)]
    else:
        return [current_branch]


# Build step docker image
def docker_build(tag: str, dockerfile: str = "Dockerfile", root_path: str = ".", args: list[str] = []) -> bool:
    logging.info(f"Building docker image {tag}")
    cmd = ["docker", "build"]
    cmd += args
    cmd += ["--iidfile", constants.CONTAINER_ID_FILE,
            "-t", tag, "-f", dockerfile, root_path]

    return run_command(cmd)


# Tags the given docker image with the tags
def docker_tag(old_tag: str, new_tag: str) -> bool:
    logging.info(f"> Tagging {old_tag} -> {new_tag}")

    return run_command(["docker", "tag", old_tag, new_tag])


# Will ush docker image to remote repository if running with env `CI=true`
def docker_push(image: str) -> bool:
    if os.getenv("CI") != "true":
        logging.info("Skipping push for local build")
        return True

    logging.info(f"Pushing docker image {image}")
    return run_command(["docker", "push", image])
