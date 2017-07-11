#!/usr/bin/env bash

# in upgrade scenario, restart service only if it was already running
if [[ "$1" == "1" ]];then
	systemctl try-restart prometheus-dev
fi

# in uninstall scenario, reload
if [[ "$1" == "0" ]];then
	systemctl --system daemon-reload
fi
