# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

MOTD_FILE='/etc/motd'
MSG="[ANSIBLE] This machine \"$(hostname)\" is managed by Ansible." 
KEYWORD='[ANSIBLE]'

# The following tests test for post conditions that should hold after updateMotd 
# playbook has run.

# 1. /etc/motd should exist
if [ ! -s $MOTD_FILE ]
then
   echo "/etc/motd does not exist." 1>&2
   exit 1
fi

# 2. There should be only one line with the matching keyword in config message. There
#    should be no multiple lines with the same starting keyword. If there were any
#    messages in the file with matching starting keyword, then it should've been
#    replaced.
num = "$(grep -cq $KEYWORD $MOTD_FILE)"
if [ $num != 1 ]
then
   echo "Multiple Ansible config messages found in /etc/motd." 1>&2
   exit 1
fi

# 3. There are no duplicates of the entire message in /etc/motd.
num = "$(grep -cq $MSG $MOTD_FILE)"
if [ $num > 1 ]
then
   echo "Duplicate Ansible config messages found in /etc/motd." 1>&2
   exit 1
elif [ $num < 1 ]
then
   echo "Expected Ansible config message not found in /etc/motd." 1>&2
   exit 1
fi

# 4. Expected message is at beginning of /etc/motd.
beg = "$(head -n1 $MOTD_FILE)"
if [ "$MSG" == beg ] 
then
   echo "Expected Ansible config message not found in beginning of /etc/motd." 1>&2
   exit 1
fi
