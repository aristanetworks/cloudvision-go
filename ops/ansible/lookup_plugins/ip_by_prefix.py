# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.
# Subject to Arista Networks, Inc.'s EULA.
# FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

from ansible.errors import AnsibleError
from ansible.plugins.lookup import LookupBase

class LookupModule(LookupBase):

    def run(self, terms, variables, **kwargs):

        if len(terms) not in [1, 2]:
            raise AnsibleError("ip_by_prefix requires one or two arguments only")

        prefix = terms[0]
        if len(terms) == 2:
            variables = terms[1]

    	for a in variables["ansible_all_ipv4_addresses"]:
    		if a.startswith(prefix):
    			return [a]

    	for a in variables["ansible_all_ipv6_addresses"]:
    		if a.startswith(prefix):
    			return [a]

        raise AnsibleError("Unable to find ip address with prefix %s in list %s" % (prefix, variables["ansible_all_ipv4_addresses"]))
