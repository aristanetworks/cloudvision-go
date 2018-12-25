#!/bin/sh
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

if [ "$#" -eq 0 ]; then
    echo "No arguments provided"
    echo "Usage: ./add-godeps.sh GO_PACKAGE_URL [...]"
    exit 1
fi

set -ex

go get -u github.com/golang/dep/cmd/dep

# If you are reading this script after receiving an error similar to this:
# ✗ Gopkg.toml already contains rules for github.com/golang/mock/mockgen/model, cannot specify a version constraint or alternate source
# see https://github.com/golang/dep/issues/705: "Constraints cannot be applied to arbitrary packages - only the root of a project"
# and the discussion at https://chat.google.com/room/AAAAgmRUTSU/ysGU8ZJSxKk
# The gist of it is dep will not work with the @master constraint if an ancestor of the subpackage being vendored exists.

dep ensure -add "$@"@master
./clean-deps.sh

set +x

echo Added dependencies: "$@"
echo If changes look good, run
echo "    git add Gopkg.toml Gopkg.lock vendor/ gopenconfig/testmodules"
