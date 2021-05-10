#! /usr/bin/env python3
"""
    Runs
        - go get ./...
        - go mod download
        - go mod tidy
        in every go step directory
"""

import os
import steps
import sys

def main():
    for step in steps.get_go_steps():
        cur_dir = os.getcwd()
        print(f"> Running in: {step}")
        os.chdir(step)
        os.system("go get ./...")
        os.system("go mod download")
        os.system("go mod tidy")
        print("Done")
        os.chdir(cur_dir)

    return 0

if __name__ == "__main__":
    sys.exit(main())