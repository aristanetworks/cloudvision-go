# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

RCLOCAL='/etc/rc.d/rc.local'

# Make sure rc.local exists.
if [ ! -s $RCLOCAL ]
then
   echo "/etc/rc.d/rc.local doesn't exist"
   exit 1
fi

# This mask would be representative of "chmod +x".
# Check permissions on rc.local.
if [ ! -x $RCLOCAL ]
then
   echo "Wrong permissions set for rc.local"
   exit 1
fi

# Check there is no code that invokes /etc/sysctl-Arora.conf
if grep -q '\/etc\/sysctl-Arora.conf' $RCLOCAL
then
   echo "sysctl-Arora.conf invocation code still in rc.local"
   exit 1
fi

# Check only one invocatio to /etc/sysctl.d/Arora.conf exists
num=$(grep -cq '\/etc\/sysctl.d\/Arora.conf' $RCLOCAL)
if [ $num > 1 ]
then
   echo "More than one code that involes etc/sysctl.d/Arora.conf found"
   exit 1
elif [ $num < 1 ]
then
   echo "No code to invoke sysctl.d/Arora.conf found"
   exit 1
fi

exit 0
