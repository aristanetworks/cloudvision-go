# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

CONFIG_FILE='/etc/ansible/ansible.cfg'

if ! grep -q 'host_key_checking = False' $CONFIG_FILE
then
   echo 'Incorrect host key checking setting' 1>&2
   exit 1
fi

if ! grep -q "ssh_args = -o ControlMaster=no" "$CONFIG_FILE"
then
   echo 'Incorrect ssh arguments setting' 1>&2
   exit 1
fi

if ! grep -q "pipelining = True" "$CONFIG_FILE"
then
   echo 'Incorrect SSH pipelining setting' 1>&2
   exit 1
fi

exit 0

