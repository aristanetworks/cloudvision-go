#!/usr/bin/env bash

# stop and disable the service if its being uninstalled
# not when its being upgraded
if [[ "$1" == "0" ]];then
	systemctl stop prometheus-dev
	systemctl disable prometheus-dev
fi
