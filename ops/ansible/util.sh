#!/bin/bash
# Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

check_cluster_name() {
	local NAME=$1
	if [ -z "$NAME" ]; then
		echo "Cluster name not provided"
		exit 1
	fi
	if [ ! -f "inventories/$1/hosts" ]; then
		echo "Cluster $1 is unknown"
		exit 1
	fi
}