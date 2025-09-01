#!/bin/bash

fullPath=$(realpath "$0")
utilDir=$(dirname "$fullPath")
# repoDir=$(dirname "$utilDir")

# echo cd "$utilDir"
# echo cd "$repoDir"

if [[ -d "$utilDir/dcl-compiler/cas-dcl-ide" ]]; then
    rm -rf "$utilDir/dcl-compiler/cas-dcl-ide"
fi

mkdir -p "$utilDir/dcl-compiler/cas-dcl-ide"
cd "$utilDir/dcl-compiler/cas-dcl-ide" || exit 1
git init
git remote add -f origin "git@github.wdf.sap.corp:CPSecurity/cas-dcl-ide.git"
git config core.sparseCheckout true
echo "/etc/scripts/" >> ".git/info/sparse-checkout"
git pull --depth=1 origin master
cd "$utilDir/dcl-compiler/cas-dcl-ide/etc/scripts/" || exit 1
./downloadDCLCompiler.sh
rm "$utilDir/dcl-compiler/dcl.jar"
cp dcl.jar "$utilDir/dcl-compiler/dcl.jar"
rm -rf "$utilDir/dcl-compiler/cas-dcl-ide"
