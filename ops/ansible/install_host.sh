#!/bin/bash

set -ex

# Get the current hostname
INSTALLHOST=$(hostname -s)
# Find the inventory for this host
INVENTORY=$(grep -l "^${INSTALLHOST}\\b" inventories/*/hosts)
if [ "$(echo "${INVENTORY}" | wc -l)" -ne "1" ]
then
	echo "ERROR: More than one inventory found for this hostname $INSTALLHOST"
	exit 1
fi

# Run ansible for this host
ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i "${INVENTORY}" -v playbook.yml -l "${INSTALLHOST}"
