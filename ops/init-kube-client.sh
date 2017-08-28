#!/bin/bash
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Setup dev cluster
kubectl config set-cluster dev \
	--server=https://172.24.32.24 \
	--certificate-authority="$ROOTDIR/ansible/inventories/dev/files/ca.pem"
kubectl config set-credentials dev \
	--certificate-authority="$ROOTDIR/ansible/inventories/dev/files/ca.pem" \
	--client-key="$ROOTDIR/ansible/inventories/dev/files/admin-key.pem" \
	--client-certificate="$ROOTDIR/ansible/inventories/dev/files/admin.pem"
kubectl config set-context dev \
	--cluster=dev \
	--user=dev
# Let's have dev the default context for now
kubectl config use-context dev

# Setup staging cluster
kubectl config set-cluster staging \
	--server=https://172.24.32.7 \
	--certificate-authority="$ROOTDIR/ansible/inventories/staging/files/ca.pem"
kubectl config set-credentials staging \
	--certificate-authority="$ROOTDIR/ansible/inventories/staging/files/ca.pem" \
	--client-key="$ROOTDIR/ansible/inventories/staging/files/admin-key.pem" \
	--client-certificate="$ROOTDIR/ansible/inventories/staging/files/admin.pem"
kubectl config set-context staging \
	--cluster=staging \
	--user=staging
# Let's have dev the default context for now
# kubectl config use-context staging
