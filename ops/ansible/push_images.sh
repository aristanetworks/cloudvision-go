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
TARGET_REPO=registry.$1.corp.arista.io

if [ -z "$1" ]; then
	echo "Cluster name not provided"
	exit 1
fi

if [ -z "$2" ]; then
	IMAGES="aeris/k8s-kafka
k8s/grafana
aeris/kafka-manager
haproxy
aeris/k8s-zookeeper
k8s/zk-prom"
else
	IMAGES=$2
fi

for image in $IMAGES; do
	docker pull "$SRC_REPO/$image"
	docker tag "$SRC_REPO/$image" "$TARGET_REPO/$image"
	docker push "$TARGET_REPO/$image"
done
