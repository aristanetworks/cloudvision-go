#!/bin/bash

# Helper script to get the yaml file of a kubernetes tasks on stdout
# without running any task

# $1 is the cluster, either "dev" or "staging"
# $2 is the kubernetes task name for which you want to output the yaml

source ./util.sh

CLUSTER=$1

check_cluster_name "$CLUSTER"

if [ -z "$2" ]; then
	echo "Task name not provided"
	exit 1
fi

DISPLAY_K8S_YAML="$2" ansible-playbook \
	-i "inventories/$CLUSTER/hosts" \
	--vault-password-file ".pass.$CLUSTER" \
	-l localhost \
	--start-at-task \
	"$2" \
	--check \
	cluster.yml
