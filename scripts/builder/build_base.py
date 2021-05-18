#! /usr/bin/env python3
"""
    This script is used to build ./base image
"""

import sys
import build


def main():
    if not build.build(base_image_build=True):
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
