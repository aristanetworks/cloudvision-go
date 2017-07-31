# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

test:
	ansible-playbook --syntax-check playbook.yml

# TODO: Make this work for the jenkins build as well
local-test:
	# dev: Test a basic host
	ansible-playbook --check -i inventories/dev/hosts -l bs200 playbook.yml
	# dev: Test a k8s master and etcd hosts
	ansible-playbook --check -i inventories/dev/hosts -l r12s24 playbook.yml
	# staging: Test a basic host
	ansible-playbook --check -i inventories/staging/hosts -l r12s12 playbook.yml

jenkins: test


.PHONY: jenkins test
