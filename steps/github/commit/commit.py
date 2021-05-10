#!/usr/bin/env python3

import sys
import github
import json
import os
import random
from subprocess import check_output
import shlex


END_STRING = "<-- END -->"


def main():
    # required params
    token = os.environ["TOKEN"]
    repo = os.environ["REPO"]
    name = os.environ["NAME"]
    filepath = os.environ["FILEPATH"]
    cmds = os.environ["COMMANDS"]

    # optional params
    base = os.environ.get("BASE") or "master"
    title = os.environ.get("TITLE") or f"StackPulse {name}"
    body = os.environ.get("BODY") or f"StackPulse {name} remediation pull request"
    commitmsg = os.environ.get("COMMIT_MSG") or f"StackPulse {name} remediation commit"
    debug = bool(os.environ.get("DEBUG") or False)

    gh = github.Github(token)
    repo = gh.get_repo(repo)

    branch = f"SP-Playbook-{name}"
    ref = repo.create_git_ref(f"refs/heads/{branch}", repo.get_branch(base).commit.sha)

    update_file = True
    try:
        orig_file = repo.get_contents(filepath)
        content = orig_file.decoded_content
    except Exception as e:
        content = ""
        update_file = False

    new_content = check_output(cmds, input=content, shell=True)
    if not update_file:
        repo.create_file(filepath, commitmsg, new_content, branch)
    else:
        repo.update_file(filepath, commitmsg, new_content, orig_file.sha, branch)

    pull_request = repo.create_pull(title=title, body=body, head=branch, base=base)
    output_dict = {
        "url": pull_request.url,
        "html_url": pull_request.html_url,
        "created_at": str(pull_request.created_at),
        "title": pull_request.title,
        "number": pull_request.number,
    }

    print(f"{END_STRING}{json.dumps(output_dict)}")


if __name__ == "__main__":
    main()
