#!/bin/bash

fullPath=$(realpath "$0")
utilDir=$(dirname "$fullPath")
repoDir=$(dirname "$utilDir")


if [[ ! -f "$utilDir/dcl-compiler/dcl.jar" ]]; then
    "$utilDir/downloadDCLCompiler.sh"
fi

cd "$utilDir/dcl-compiler" || exit 1

for i in "$repoDir/tests/"*; do
    if [[ -d "$i/dcl" ]]; then
        java -jar "$utilDir/dcl-compiler/dcl.jar" -compileTestToDcn -dcn=pretty -out="$i/dcn" "$i/dcl"
        find "$i" -name "*.rego" -type f -delete
    fi
done
