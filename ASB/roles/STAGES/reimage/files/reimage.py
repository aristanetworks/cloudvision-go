#!/usr/bin/env python
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import subprocess
import sys
import re

# ----------------------------------------------------------------------------------
# Variable Definitions
# ----------------------------------------------------------------------------------
config_stats_file = "/tmp/CURR_CONFIG"


# ----------------------------------------------------------------------------------
# Util Methods
# ----------------------------------------------------------------------------------
def _c( cmd ):
   '''
   Make a call, ignoring the output.
   '''
   try:
      # XXX: shell=True is a security hole but need this option for chained commands
      rc = subprocess.call( "sudo " + cmd, shell=True )
   except:
      if rc:
         print "ERROR/IGNORED: <'%s'> returned error code %d" % ( cmd, rc )

def _o( cmd ):
   '''
   Make a call and return output.
   '''
   # XXX: shell=True is a security hole but need this option for chained commands
   try:
      return subprocess.check_output( "sudo " + cmd, shell=True ).strip( '\n' )
   except:
      print "ERROR/IGNORED: <'%s'> return error, ignoring" % cmd
      return []

def disks():
   '''
   Return list of available disks on this machine.
   '''
   devs = _o( "parted -s -l | awk -F [/:] '/^Disk \/dev\/[^md]/ { print $3 }'" )
   devs = devs.split( '\n' )
   return devs if all( devs ) else None

def md_devices():
   '''
   Return list of md devices that are part of the RAID arrays.
   '''
   mddevs = _o( "cat /proc/mdstat | awk '/^[md]/ {print $1}'" )
   mddevs = mddevs.split( '\n' )
   return mddevs if all( mddevs ) else None

def partitions():
   '''
   Return list of partitions on the disks.
   '''
   prts = _o( "ls /dev/sd[abc]* | awk -F [/] '/\/dev\/sd[abc][1-9]/ {print $3}'" )
   prts = prts.split( '\n' )
   return prts if all( prts ) else None


# ----------------------------------------------------------------------------------
# Methods
# ----------------------------------------------------------------------------------
def sanity():
   '''
   Check that this machine is in a sane state before being reimaged.
   '''
   def _write( configs ):
      # Write out the config stats to a file
      with open( config_stats_file, 'w' ) as f:
         for conf in configs:
            f.write( '%s\n' % conf[ 0 ] )
            f.write( '%s\n\n\n' % conf[ 1 ] )
      f.close()

   configs = []

   # Store previous config states for this machine
   configs.append( ( 'ifconfig', _o( "ifconfig" ) ) )
   configs.append( ( 'ip link', _o( "ip link" ) ) )
   configs.append( ( 'cat /proc/mdstat', _o( "cat /proc/mdstat" ) ) )
   configs.append( ( 'df -h', _o( "df -h" ) ) )
   configs.append( ( 'parted -s -l',  _o( "parted -s -l" ) ) )

   # get list of disk names ( i.e.: sda, sdb, sdc... )
   dsks = disks()
   if not dsks:
      print >> sys.stderr, "No disk drives detected. Something is very wrong!"
      _write( configs )
      sys.exit( 1 )

   # Check every disk
   for d in dsks:
      _c( "smartctl -a /dev/%s | grep -q PASSED" % d )
      configs.append( ( 'smartctl -a /dev/%s' % d,  
                        _o( "smartctl -a /dev/%s" % d ) ) )

      template = "smartctl -a /dev/%s | awk '/%%s/ {print $10}'" % d
      rrer = int( _o( template % 'Raw_Read_Error_Rate' ) )

      # XXX: These output from smartctl are vendor specific! 
      # ST/Seagate does not have "Reallocated_Event_Count"
      rec = _o( template % 'Reallocated_Event_Count' )
      if rec:
         rec = int( rec )
      else:
         rec = -1
         configs.append( "Reallocated_Event_Count field check was skipped because it was unavailable" )

      ou = int( _o( template % 'Offline_Uncorrectable' ) )
   
      # If these three fields are not 0, then disk is bad, do not proceed.
      if ( rrer > 100 ) or ( rec > 100 ) or ( ou > 100 ):
         err = """ 
Sanity Check Failed - these values from smartctl failed sanity test:
Raw_Read_Error_Rate: %s
Reallocated_Event_Count: %s
Offline_Uncorrectable: %s
"""            % ( rrer, "unavailable" if rec == -1 else rec, ou )

         print >> sys.stderr, err
         _write( configs )
         sys.exit( 1 )

   _write( configs )


def wipe():
   '''
   Wipe the machine of old partitions and MD arrays.
   '''
   mddevs = md_devices() # i.e. md124, md125
   dsks = disks() # i.e. sda, sdb, sdc
   prts = partitions() # i.e. sda1, sda2, sda3

   if not dsks:
      print >> sys.stderr, "No disk drives detected. Something is very wrong!"
      sys.exit( 1 )
   if not mddevs:
      print >> sys.stderr, "No MD arrays detected."
   if not prts:
      print >> sys.stderr, "No partitions on the disk drives detected."

   # Disable page swaps.
   _c( "swapoff -a" )

   # Unmount and stop all MD arrays.
   if mddevs:
      for md in mddevs:
         _c( "umount /dev/%s" % md )
         _c( "mdadm --stop /dev/%s" % md )

   # Zero out the superblock for every partition on each disk
   if prts:
      for p in prts:
         _c( "mdadm --zero-superblock --force /dev/%s" % p )

   # Zero out any potential superblock sitting at the header of each disk
   for d in dsks:
      _c( "mdadm --zero-superblock --force /dev/%s" % d )

   # Remove each partition from the table
   if prts:
      for p in prts:
         d, pn = re.match( r'([a-z]+)([0-9]+)', p ).groups(0)
         _c( "parted -s /dev/%s rm %s" % (d, pn) )

   # Inform Kernel of the partition table changes
   _c( "partprobe" )


def verify():
   '''
   Verify that disk wiping process was successful before setting up disk again.
   '''
   # Check there are no RAID arrays
   raids = _o( 'cat /proc/mdstat' )
   expected = "Personalities : \nunused devices: <none>"
   if raids != expected:
      print >> sys.stderr, "RAID arrays were not properly removed"
      sys.exit( 1 )

   # Check either only /dev/md0 exists or no /dev/md* exists
   md0Exists = False
   output = _o( 'ls -l /dev/md*' ).split( '\n' )
   if len( output ) == 1 and 'md0' in output[ 0 ]:
      print "md0 still exists but this is okay. Ignoring."
      md0Exists = True
   elif len( output ) > 1:
      print >> sys.stderr, "Not all MD devices were removed!"
      sys.exit( 1 )

   # Check that disks have all of their partitions wiped
   dsks = disks()
   for d in dsks:
      output = _o( 'parted -s /dev/%s print | tail -n2' % d )
      if output != 'Number  Start  End  Size  File system  Name  Flags':
         if not md0Exists:
            print >> sys.stderr, "Drive /dev/%s was not properly wiped." % d
            sys.exit( 1 )


# ----------------------------------------------------------------------------------
# Main
# ----------------------------------------------------------------------------------
def main():
   arg = sys.argv[ 1 ]
   if arg == "sanity":
      sanity()
   elif arg == "wipe":
      wipe()
   elif arg == "verify":
      verify()

if __name__ == "__main__":
   main()

