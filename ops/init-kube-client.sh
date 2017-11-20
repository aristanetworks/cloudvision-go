#!/bin/bash
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

set -e

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

NARGS=$#

if [ "$NARGS" -eq 0 ]; then
	clusters="dev staging"
else
	clusters=$@
fi

initcluster() {
	local CLUSTER=$1

	# Create folder to save certs for this cluster
	mkdir -p "$HOME/.kube/certs/$CLUSTER"

	# Decrypt certificates for the cluster
	for f in ca admin-key admin;
	do
		ansible-vault \
			--vault-password-file "$ROOTDIR/ansible/.pass.$CLUSTER" \
			--output="$HOME/.kube/certs/$CLUSTER/${f}.pem" \
			decrypt \
			"$ROOTDIR/ansible/inventories/$CLUSTER/files/k8s/${f}.pem"
	done

	# Setup kubectl config for the cluster
	local SERVER
	case $CLUSTER in
		"dev"|"staging")
			SERVER=https://master.${CLUSTER}-infra.corp.arista.io
			;;
		*)
			SERVER=https://127.0.0.1:8240
			;;
	esac
	kubectl config set-cluster "$CLUSTER" \
		--server=$SERVER \
		--certificate-authority="$HOME/.kube/certs/$CLUSTER/ca.pem"
	kubectl config set-credentials "$CLUSTER" \
		--certificate-authority="$HOME/.kube/certs/$CLUSTER/ca.pem" \
		--client-key="$HOME/.kube/certs/$CLUSTER/admin-key.pem" \
		--client-certificate="$HOME/.kube/certs/$CLUSTER/admin.pem"
	kubectl config set-context "$CLUSTER" \
		--cluster="$CLUSTER" \
		--user="$CLUSTER"
}

for c in $clusters; do
	initcluster "$c"
done

if [ "$NARGS" -eq 0 ]; then
	# Let's have dev the default context for default case
	kubectl config use-context dev
fi
