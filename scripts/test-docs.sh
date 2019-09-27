#!/usr/bin/env bash

set -e
set -x

echo "Building docs with Sphinx"
make clean
make html

echo "Checking grammar and style"
vale --glob='*.rst' .

echo "Checking links"
make linkcheck
