#!/bin/bash
# Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

# This script will push all the images needed byt an Aeris cluster to work
# from a source repo to the cluster repo.
#
# The first argument is mandatory: It's the name of the cluster (dev, staging, etc)
# This name is the folder name in the "inventories" folder
# The second argument is optional: It's the specific docker image you want to push

# TODO: Able to push specific (multiple) tags per image.
# Maybe use docker-library and all the tags in each file?

set -ex

SRC_REPO=registry.docker.sjc.aristanetworks.com:5000
TARGET_REPO=registry.$1-infra.prod.arista.io
CLUSTER=$1

# TODO: Uncomment when we have the VPN
# source ./util.sh

# check_cluster_name "$CLUSTER"

if [ -z "$2" ]; then
	IMAGES="aeris-apiserver
aeris-cli
aeris-dispatcher
aeris-ingest
aeris-recorder
aeris/k8s-hadoop
aeris/k8s-hbase
aeris/k8s-kafka
aeris/k8s-zookeeper
aeris/kafka-manager
k8s/grafana
k8s/proxy
k8s/proxy2
k8s/zk-prom"
else
	IMAGES=$2
fi

# TODO: Use VPN so we don't need all this (including TEMPORARY)
# Start kubectl port-forward to access the cluster internal docker registry
#kubectl port-forward "$(kubectl get po -l app=registry --template '{{(index .items 0).metadata.name}}')" 443:5000

# TEMPORARY:
# 1. Use to a r123sXX server (ssh core@r123sXX -A) (r123s17 for instance: ssh core@r123s17.sjc.aristanetworks.com -A)
# 2. Copy this script
# 3. Define "registry.ovh-bhs-infra.prod.arista.io 127.0.0.1" in local /etc/hosts
# 4. Start ssh tunnel to the prod docker registry:
#    ssh -L 4430:172.18.181.165:443 core@ns545236.ip-144-217-181.net -N
# 5. Start local port forwarding from 443 to 4430 as root:
#    sudo socat tcp-l:443,fork,reuseaddr tcp:127.0.0.1:4430
# 6. Run the script to push the images.

cleanup () {
	echo "Stopping port forward..."
	kill $PID
	wait
	echo "Port forward stopped."
}
trap cleanup EXIT

for image in $IMAGES; do
	docker pull "$SRC_REPO/$image"
	docker tag "$SRC_REPO/$image" "$TARGET_REPO/$image"
	docker push "$TARGET_REPO/$image"
done
