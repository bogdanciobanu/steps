#!/usr/bin/env python3

import glob
import json
from html.parser import HTMLParser


def validate_catalog():
    catalog = json.loads(open("./vendors/vendors.json", "rb").read())
    svgs = [x.split("/")[1] for x in glob.glob("vendors/*.svg")]
    for vendor in catalog:
        vendor_type, icon = vendor["vendor_type"], vendor["icon"]
        assert vendor_type != ""
        if vendor_type == "custom":
            continue
        assert icon in svgs


    icons = set(vendor["icon"] for vendor in catalog if vendor["icon"] != "")
    vendors = set(
        vendor["vendor_type"] for vendor in catalog if vendor["vendor_type"] != ""
    )
    assert len(set(vendors)) == len(vendors)

    icons_diffrence = set(svgs).symmetric_difference(icons)
    assert len(icons_diffrence) == 0, icons_diffrence



def validate_svg():
    class SVGParser(HTMLParser):
        def __init__(self):
            super().__init__()

            self.attr = {"width": False, "height": False}

        def has_width(self):
            return self.attr["width"]

        def has_height(self):
            return self.attr["height"]

        def handle_starttag(self, tag, attrs):
            if "svg" == tag:
                for attr in attrs:
                    self.attr[attr[0]] = True

    files = glob.glob("./vendors" + "/*.svg")

    for f in files:
        print(f"Validating: {f}", end=" ")

        svg_parser = SVGParser()
        with open(f, "r+") as svgfile:
            svg_parser.feed(svgfile.read())
        assert svg_parser.has_height() == True
        assert svg_parser.has_width() == True

        print("- OK")

validate_catalog()
validate_svg()