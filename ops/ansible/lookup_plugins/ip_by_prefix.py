# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

from ansible.errors import AnsibleError
from ansible.plugins.lookup import LookupBase

class LookupModule(LookupBase):

    def run(self, terms, variables, **kwargs):

    	if len(terms) != 1:
    		raise AnsibleError("ip_by_prefix requires one argument only")

    	for a in variables["ansible_all_ipv4_addresses"]:
    		if a.startswith(terms[0]):
    			return [a]

    	for a in variables["ansible_all_ipv6_addresses"]:
    		if a.startswith(terms[0]):
    			return [a]

        raise AnsibleError("Unable to find ip address with prefix %s in list %s" % (terms[0], variables["ansible_all_ipv4_addresses"]))
