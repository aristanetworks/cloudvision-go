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
  if [ -f "/etc/nginx/sites-enabled/$d.conf" ]; then
    echo "SKIPPING: $d nginx config already present"
  else
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
  return 302 https://www.arista.com/;
}
}
EOF

  fi

done

function publish {
  local NAME
  for d in /etc/letsencrypt/live/*; do
    NAME=$(basename "$d")
    echo "Publishing letsencrypt cert $NAME"

    echo "apiVersion: v1
  kind: Secret
  metadata:
  name: lecert-${NAME}
  type: Opaque
  data:
  cert.pem: $(base64 -w 0 < "$d/cert.pem")
  chain.pem: $(base64 -w 0 < "$d/chain.pem")
  fullchain.pem: $(base64 -w 0 < "$d/fullchain.pem")
  privkey.pem: $(base64 -w 0 < "$d/privkey.pem")
  " | kubectl apply -f -

    echo "Published letsencrypt cert $NAME"
  done
}

# Start nginx
nginx &

# Run for the first time
# shellcheck disable=SC2034
for d in $DOMAINS;
do
  if [ -d "/etc/letsencrypt/live/$d" ]; then
    echo "SKIPPING: $d certificate already present"
  else
    certbot --nginx --agree-tos -n -m ops-dev@arista.com -d "$d"
  fi
  publish
done;

# Run every day
# Poor man/simple cron job
while true
do
  sleep 86400
  certbot renew
  publish
done