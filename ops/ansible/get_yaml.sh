#!/bin/bash

# Helper script to get the yaml file of a kubernetes tasks on stdout
# without running any task

# $1 is the cluster, either "dev" or "staging"
# $2 is the kubernetes task name for which you want to output the yaml

DISPLAY_K8S_YAML="$2" ansible-playbook \
	-i "inventories/$1/hosts" \
	-l localhost \
	--start-at-task \
	"$2" \
	--check \
	cluster.yml
