#!/usr/bin/env python
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import logging
import unittest
import subprocess
import sys
import shlex

# Basic testing parameters
testhost = "bs306"
TOPLEVEL_PLBK = "AroraConfig.yml"
DIR_PLBK = "playbooks"

# Logging setup
logging.basicConfig( level=logging.INFO )
logger = logging.getLogger( 'TestPlaybooks' )

# Prints error message and exits with failure code 1
def errExit( msg ):
   logger.error( msg )
   sys.exit( 1 )

# Runs command on Ansible Master (ansible@) and checks output or return code
def cmdAnsibleMaster( cmd, success_msg, err_msg, expected="" ):
   try:
      ssh_cmd = ( "sshpass -p arastra ssh -o StrictHostKeyChecking=no "
                  "root@ansible \"su ansible && %s\"" )
      if expected:
         if expected not in subprocess.check_output( shlex.split( ssh_cmd % cmd ) ):
            errExit( err_msg )
         else:
            logger.info( success_msg )
      else:
         subprocess.check_call( shlex.split( ssh_cmd % cmd ) )
         logger.info( success_msg )
   except subprocess.CalledProcessError as e:
      errExit( "%s : %s" % ( e, err_msg ) )

def scpAnsibleMaster( src, dest, success_msg, err_msg, directory=False ):
   try:
      scp_cmd = ( "sshpass -p arastra scp -o StrictHostKeyChecking=no %s %s "
                  "root@ansible:/home/ansible/%s" )
      opt = "-r" if directory else ""
      subprocess.check_call( shlex.split( scp_cmd % ( opt, src, dest ) ) )
      logger.info( success_msg )
   except subprocess.CalledProcessError as e:
      errExit( "%s : %s" % ( e, err_msg ) )

class TestPlaybooks( unittest.TestCase ):
   def setUp( self ):
      # Install packages that should be installed on the jenkins node.
      subprocess.call( shlex.split( "sudo yum install -y sshpass" ) )

      # Sanity host ping test
      cmdAnsibleMaster( "ansible all -m ping -i %s," % testhost,
         "Test host %s successfully pinged and checked to be alive." % testhost,
         "Failed to check if test host %s is alive." % testhost,
         "SUCCESS" )

      # Check test host can be used for testing at all
      cmdAnsibleMaster( "ssh %s svp info bs306" % testhost,
         "Test host %s in installed state. Proceeding with test." % testhost,
         "Test host %s not in installed state. Aborting test." % testhost,
         "installed" )

      # Check test host is set in build role
      cmdAnsibleMaster( "ssh %s svp info bs306" % testhost,
         "Test host %s is setup with build role. Proceeding with test." % testhost,
         "Test host %s not setup with build role. Aborting test." % testhost,
         "build" )

      # Check ansible version on test host
      expected_ansible_ver = "2.2.0.0"
      cmdAnsibleMaster( "ssh %s ansible --version" % testhost,
         "Test host has expected ansible version %s" % expected_ansible_ver,
         "Test host %s does not have expected ansible version (required %s)." %
            ( testhost, expected_ansible_ver ),
         expected_ansible_ver )

      # Copy over playbooks to ansible master to be run
      scpAnsibleMaster( TOPLEVEL_PLBK,
         TOPLEVEL_PLBK,
         "Copied top level Ansible playbook %s to Ansible Master." % TOPLEVEL_PLBK,
         "Failed to copy top level Ansible playbook %s to Ansible Master."
            % TOPLEVEL_PLBK )

      # Copy over rest of the required playbooks to Ansible Master
      scpAnsibleMaster( DIR_PLBK,
         "",
         "Copied over rest of the playbooks in %s to Ansible Master." % DIR_PLBK,
         "Failed to copy rest of the playbooks in %s to Ansible Master." % DIR_PLBK,
         directory=True )

   def tearDown( self ):
      # Clean up.
      cmdAnsibleMaster( "rm /home/ansible/%s " % TOPLEVEL_PLBK,
                        "Cleaned up %s from Ansible Master" % TOPLEVEL_PLBK,
                        "Failed to clean up top level playbook %s" % TOPLEVEL_PLBK )
      cmdAnsibleMaster( "rm -r /home/ansible/%s" % DIR_PLBK,
                        "Cleaned up %s from Ansible Master" % DIR_PLBK,
                        "Failed to clean up rest of the playbooks in %s" % DIR_PLBK )

      # Uninstall the packages temporarily installed just for this test
      subprocess.call( shlex.split( "sudo yum erase sshpass" ) )

   def test( self ):
      # Run the playbooks on test host
      cmdAnsibleMaster( "ansible-playbook -i %s, /home/ansible/%s" %
                           ( testhost, TOPLEVEL_PLBK ),
                        "Testing playbooks finished.",
                        "Failed to push the playbooks that needed to be tested." )

if __name__ == '__main__':
   unittest.main( verbosity=2 )
