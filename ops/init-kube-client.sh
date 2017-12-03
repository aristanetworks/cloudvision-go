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
			SERVER=https://master.${CLUSTER}-infra.prod.arista.io
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

	if [ ! -f "$ROOTDIR/ansible/inventories/$CLUSTER/files/openvpn/ta.key" ]; then
		echo "No VPN to configure for this cluster"
		return
	fi

	echo "Configuring VPN for this cluster"
	# Generates the client.conf for the vpn
	mkdir -p "$HOME/.openvpn/certs/$CLUSTER"
	for f in ca.pem ta.key client-key.pem client.pem ;
	do
		ansible-vault \
			--vault-password-file "$ROOTDIR/ansible/.pass.$CLUSTER" \
			--output="$HOME/.openvpn/certs/$CLUSTER/${f}" \
			decrypt \
			"$ROOTDIR/ansible/inventories/$CLUSTER/files/openvpn/${f}"
	done
	cat > "$HOME/.openvpn/$CLUSTER-client.conf" <<EOF
client
dev tun
proto udp
remote ns545238.ip-144-217-181.net 1194
persist-key
persist-tun
ca $HOME/.openvpn/certs/$CLUSTER/ca.pem
cert $HOME/.openvpn/certs/$CLUSTER/client.pem
key $HOME/.openvpn/certs/$CLUSTER/client-key.pem
tls-auth $HOME/.openvpn/certs/$CLUSTER/ta.key 1
key-direction 1
remote-cert-tls server
cipher AES-256-CBC
tls-version-min 1.2
auth SHA256
comp-lzo
verb 3
EOF

	echo "VPN configured for this cluster. Conf file path is $HOME/.openvpn/$CLUSTER-client.conf"
}

for c in $clusters; do
	initcluster "$c"
done

if [ "$NARGS" -eq 0 ]; then
	# Let's have dev the default context for default case
	kubectl config use-context dev
fi
