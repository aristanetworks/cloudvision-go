# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

MOTD_FILE='/etc/motd'
MSG="[ANSIBLE] This machine \"$(hostname)\" is managed by Ansible."
KEYWORD='[ANSIBLE]'

if [ -s $MOTD_FILE ]
then
   if grep -q "$KEYWORD" "$MOTD_FILE"
   then
      echo "Replacing previous Ansible information in /etc/motd."

      #sed -i'' '/[ANSIBLE]/${MSG}/' $MOTD_FILE
      #sed -i ','"[ANSIBLE]"',c\'"$MSG"',' $MOTD_FILE
      sed -i '/[ANSIBLE]/c\'"$MSG" $MOTD_FILE
   else
      echo "Adding Ansible information to /etc/motd."
      echo $MSG >> $MOTD_FILE
   fi
else
   echo "Creating a new /etc/motd and writing Ansible information."
   echo $MSG > $MOTD_FILE
fi
