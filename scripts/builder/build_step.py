#! /usr/bin/env python3

"""
    Used to build a step
"""

import build
import sys


def main():
    if not build.build():
        return 1

    return 0


# Main entry point
if __name__ == "__main__":
    sys.exit(main())
