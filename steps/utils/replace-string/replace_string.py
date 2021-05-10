#!/usr/bin/env python3

import os
import json
import base64
import re


def main():
    # required params
    original = os.environ["ORIGINAL_BASE64_STRING"]
    subStrSrc = os.environ["SUB_STR_SRC"]
    subStrDest = os.environ["SUB_STR_DEST"]
    isRegex = bool(os.environ.get("REGEX") == "True")

    decoded_content = base64.b64decode(original).decode("utf-8")
    if isRegex:
        replaced_string = re.sub(subStrSrc, subStrDest, decoded_content)
    else:
        replaced_string = decoded_content.replace(subStrSrc, subStrDest)

    result = base64.b64encode(replaced_string.encode("ascii")).decode("ascii")

    output_object = {}
    output_object["output"] = result
    json_result = json.dumps(output_object)
    print("<-- END -->")
    print(json_result)


if __name__ == "__main__":
    main()
