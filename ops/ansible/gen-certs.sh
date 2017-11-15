#!/bin/bash
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

# This script will generate the certificates for a k8s cluster

set -ex

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CLUSTER=$1

source $ROOTDIR/util.sh

check_cluster_name "$CLUSTER"

shift

ansible-playbook \
	gen-certs.yml \
	-i "$ROOTDIR/inventories/$CLUSTER/hosts" \
	--vault-password-file "$ROOTDIR/.pass.$CLUSTER" \
	"$@"
