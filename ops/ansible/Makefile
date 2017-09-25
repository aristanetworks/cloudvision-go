# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

test:
	for i in inventories/*/hosts; do \
		ansible-playbook --syntax-check playbook.yml -i $$i; \
	done



# TODO: Make this work for the jenkins build as well
local-test:
	# dev: Test a basic host
	ansible-playbook --check -i inventories/dev/hosts -l bs200.sjc.aristanetworks.com playbook.yml
	# dev: Test a k8s master and etcd hosts
	ansible-playbook --check -i inventories/dev/hosts -l r123s24.sjc.aristanetworks.com playbook.yml
	# staging: Test a basic host
	ansible-playbook --check -i inventories/staging/hosts -l r123s12.sjc.aristanetworks.com playbook.yml

jenkins: test


.PHONY: jenkins test
