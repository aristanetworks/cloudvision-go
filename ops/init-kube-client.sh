#!/bin/bash
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p "$HOME/.kube/certs/dev"

# Decrypt certificates for dev cluster
for f in ca admin-key admin;
do
	ansible-vault \
		--vault-password-file "$ROOTDIR/ansible/.pass.dev" \
		--output="$HOME/.kube/certs/dev/${f}.pem" \
		decrypt \
		"$ROOTDIR/ansible/inventories/dev/files/k8s/${f}.pem"
done

# Setup dev cluster
kubectl config set-cluster dev \
	--server=https://master.dev-infra.corp.arista.io \
	--certificate-authority="$HOME/.kube/certs/dev/ca.pem"
kubectl config set-credentials dev \
	--certificate-authority="$HOME/.kube/certs/dev/ca.pem" \
	--client-key="$HOME/.kube/certs/dev/admin-key.pem" \
	--client-certificate="$HOME/.kube/certs/dev/admin.pem"
kubectl config set-context dev \
	--cluster=dev \
	--user=dev
# Let's have dev the default context for now
kubectl config use-context dev

mkdir -p "$HOME/.kube/certs/staging"

# Decrypt certificates for staging cluster
for f in ca admin-key admin;
do
	ansible-vault \
		--vault-password-file "$ROOTDIR/ansible/.pass.staging" \
		--output="$HOME/.kube/certs/staging/${f}.pem" \
		decrypt \
		"$ROOTDIR/ansible/inventories/staging/files/k8s/${f}.pem"
done

# Setup staging cluster
kubectl config set-cluster staging \
	--server=https://master.staging-infra.corp.arista.io \
	--certificate-authority="$HOME/.kube/certs/staging/ca.pem"
kubectl config set-credentials staging \
	--certificate-authority="$HOME/.kube/certs/staging/ca.pem" \
	--client-key="$HOME/.kube/certs/staging/admin-key.pem" \
	--client-certificate="$HOME/.kube/certs/staging/admin.pem"
kubectl config set-context staging \
	--cluster=staging \
	--user=staging
# Let's have dev the default context for now
# kubectl config use-context staging
