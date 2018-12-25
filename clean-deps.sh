#!/bin/bash
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

# This script removes unneeded files from vendored dependencies.
# TODO: filter out arista/* packages to make it work with gopenconfig

case $# in
  (0)
    targetdir=vendor
    ;;
  (1)
    targetdir=$1
    ;;
  (*)
    echo >&2 'usage: ./clean-deps.sh [directory]'
    exit 1
    ;;
esac

case `uname -s` in
  (Darwin)
    DARWIN_FIND_FLAGS='-E'
    ;;
  (Linux)
    LINUX_FIND_FLAGS='-regextype posix-extended'
    ;;
  (*)
    echo >&2 'unsupported platform'
    exit 1
    ;;
esac

# Filter out files for platforms we don't intend to ever support.
find ${DARWIN_FIND_FLAGS} "$targetdir" ${LINUX_FIND_FLAGS} \
  -regex '.*_(arm[64lbex]*|mips[64lbex]*|ppc[64lbex]*|s390x|dragonfly|freebsd|netbsd|openbsd|solaris|plan9|aix|sparc64|riscv64|windows|test)(_[a-z0-9]+)?\.(s|go|pl)$' -print0 \
  | xargs -0 rm -v
# Filter out these platforms based on build tags.
egrep -lr '^// \+build (ignore|appengine|android|dragonfly|freebsd|nacl|netbsd|openbsd|plan9|aix|solaris|windows'\
'|amd64p32|arm|armbe|arm64|arm64be|ppc64|ppc64le|mips|mipsle|mips64|mips64le|mips64p32|mips64p32le|ppc|s390|s390x|sparc|sparc64|riscv64| |,)+$' "$targetdir" \
  | xargs rm -v

# Remove other useless files: README files and markdown files, generated
# Python code (usually protobuf code), build-related files, etc.  The second
# part of the matching below is to make sure we preserve files that are used
# by license.py (e.g. when the license file is named LICENSE.md, we don't want
# to delete it just because it matched '*.md').
find "$targetdir" \( \
  -name .travis.yml \
  -o -name 'README*' \
  -o -name '*.pdf' \
  -o -name '*.md' \
  -o -name '*.sh' \
  -o -name '*.p[ly]' \
  -o -name Dockerfile \
  -o -name Makefile \
  -o -name Makefile.am \
  -o -name configure.ac \
  -o -name CMakeLists.txt \
  -o -name BUILD.bazel \
  \) -a ! \( \
  -name 'AUTHORS*' \
  -o -name 'CONTRIBUTORS*' \
  -o -name 'COPY*' \
  -o -name 'LICEN*' \
  -o -name 'NOTICE*' \
  \) -print0 | xargs -0 rm -v
