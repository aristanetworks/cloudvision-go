The files in this directory will need to be copied to the appropriate
locations on an NFS/HTTP server - see redis.anaconda-ks.cfg for the location
of the current (NFS) server and modify as appropriate.

For example, if $NFSROOT is the directory on the NFS server that the
kickstart file mounts from, the following copy command
   cp bond* rc.local redis.conf $NFSROOT
should be sufficient to populate the directory to allow kickstart to work.

The kickstart config file (redis.anaconda-ks.cfg) should be copied to
wherever the inst.ks points at in the PXE boot directive. For example,
using $NFSROOT as above, you can copy it like this:
   cp redis.anaconda-ks.cfg $NFSROOT
and add the PXE boot directive as:
kernel .../vmlinuz ... inst.ks=nfs:$NFSSERVER:$NFSROOT/redis.anaconda-ks.cfg ...
See AID/3212 for details

Obviously you don't have to use NFS - http (or even tftp, if you are very UDP
inclined) should work fine (most of the time). The nice thing about using NFS
is that a temporary glitch on the network will not leave you with a half-baked
semi-installed system.
