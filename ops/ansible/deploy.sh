#!/bin/bash
# Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

# This script will deploy a cluster

# The first argument is mandatory: It's the name of the cluster (dev, staging, etc)
# This name is the folder name in the "inventories" folder

# The other args will be passed to "ansible-playbook"

set -ex

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CLUSTER=$1

source ./util.sh

check_cluster_name "$CLUSTER"

shift

ansible-playbook \
	cluster.yml \
	-i "$ROOTDIR/inventories/$CLUSTER/hosts" \
	--vault-password-file "$ROOTDIR/.pass.$CLUSTER" \
	"$@"
