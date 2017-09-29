#!/bin/sh

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

set -e

# Add our git config for "git push review"
if ! grep -q '\[remote "review"\]' $REPO_ROOT/.git/config
then
   cat $REPO_ROOT/gitconfig-review >> $REPO_ROOT/.git/config
fi

# Add the gerrit hook that adds change IDs
gitdir=$(git rev-parse --git-dir); scp -p -P 29418 gerrit:hooks/commit-msg ${gitdir}/hooks/
