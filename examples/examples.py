#!/usr/bin/env python3

import re

# refreshReadmeSnippet: snippet2
matches = re.findall("[a-z]+", "abc cde fgf")
# refreshReadmeSnippet: snippet2

assert matches == ["abc", "cde", "fgf"]

