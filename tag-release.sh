#!/bin/zsh

v="$1"
if [[ -z $v ]]
then
    print "usage: $0 <version>"
    exit 1
fi

if [[ $v == v* ]]
then
    print "do not add the v prefix, tag should look like \"1.5\""
    exit 1
fi

grep -q "version = \"$v\"" version.go
if [[ $? != 0 ]]
then
    print "fix version.go first"
    exit 1
fi
echo git tag -s v$v -m \"tag gjfy v$v\"
