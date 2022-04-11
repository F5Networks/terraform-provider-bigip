#!/usr/bin/env bash

function tfAccTestsLint {
  echo "==> Checking acceptance test terraform blocks are formatted..."
  terrafmt version

files=$(find ./bigip -type f -name "*_test.go")
error=false

for f in $files; do
  terrafmt diff -c -q -f "$f" || error=true
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make tffmtfix"
  echo ""
  echo "format only acceptance test config blocks:"
  echo "$ find ./bigip | egrep \"_test.go\" | sort | while read f; do terrafmt fmt -f \$f; done"
  echo ""
  echo "format a single test file:"
  echo "$ terrafmt fmt -f ./bigip/resource_test.go"
  echo ""
  exit 1
fi

}

function tfDocsLint {
  echo "==> Checking docs terraform blocks are formatted..."

files=$(find ./docs -type f -name "*.md")
error=false

for f in $files; do
  terrafmt diff -c -q -f "$f" || error=true
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make tffmtfix"
  echo ""
  echo "format only docs config blocks:"
  echo "$ find docs | egrep \".md\" | sort | while read f; do terrafmt fmt -f \$f; done"
  echo ""
  echo "format a single test file:"
  echo "$ terrafmt fmt -f ./docs/resources/resource.md"
  echo ""
  exit 1
fi

}

function main {
  tfAccTestsLint
  tfDocsLint

}

main
