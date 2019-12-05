#!/bin/bash
# Copyright (c) 2019 Arista Networks, Inc.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the COPYING file.

notice='Copyright \(c\) 20[0-9][0-9] Arista Networks, Inc.'
files=`git diff-tree --no-commit-id --name-only --diff-filter=ACMR -r HEAD | \
	egrep '\.(go|proto|py|sh|y|l)$' | grep -v '\.pb\.go$' | \
    grep -v '\.gen\.go$' | grep -v '.*\.proto$'| grep -vE '_pb2(_grpc)?\.py'`
status=0

for file in $files; do
	if ! [[ -e "$file" ]]; then
		continue
	fi
	if ! egrep -q "$notice" $file; then
		echo "$file: missing or incorrect copyright notice"
		status=1
	fi
done

exit $status
