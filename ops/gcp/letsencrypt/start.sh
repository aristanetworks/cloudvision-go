#!/bin/bash
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

# We don't set "-e" option because we don't want this script to fail.
# This script is the CMD for the docker container.
# If this script fails, the container will exit, and kubernetes will restart it
# For an usual container, it's fine, but this container will run
# certbot at startup.
# It means that if k8s keeps restarting the container because this script is
# failing, let's encrypt will rate limit us and we will not be able to
# fix the script and retry.
# (or connect to the container in order to fix the script).
# So, we let this script and its infinite for loop run even in case of errors
# We might need to find a way to monitor this (to monitor possible errors)

set -x

# Create nginx config for all the declared domains
for d in $DOMAINS;
	do
	cat <<EOF > "/etc/nginx/sites-enabled/$d.conf"
server {
	listen 80;
	listen [::]:80;

	server_name $d;

	location /.ping {
		add_header Content-Type text/plain;
		return 200 "pong";
	}

	location / {
		return 407 "407 VPN required";
	}
}
EOF
	done;

# Start nginx
nginx &

# Start sshd
/usr/sbin/sshd &

# Run for the first time
# shellcheck disable=SC2034
for d in $DOMAINS;
do
	certbot --nginx --agree-tos -n -m ops-dev@arista.com -d "$d"
done;

# Run every day
# Poor man/simple cron job
while true
do
	sleep 86400
	certbot renew
done
