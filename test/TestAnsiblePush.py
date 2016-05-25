#!/usr/bin/env python
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import os
import sys
import subprocess
import logging
import unittest
import tempfile

# Description
# ===========
# This test provides a framework *ONLY* to test playbooks under 'playbooks' folder
# in 'ardc-config', via Ansible "push" mechanism. This framework should be used to
# test the affects of playbooks on a simulated datacenter environment.
#
# This framework is intended to be an aid as part of the automated server
# provisioning process (AID 3269).
#
# Assumptions/Preconditions
# =========================
# This test has following assumptions to set up the testing environment:
#     1. Ansible is already installed on main Ansible server and managed servers
#        (you can see which programs/packages are expected to be on these servers by
#         checking the Dockerfile under test/dockerfiles/). This step should've been
#        already completed by the previous steps in the server provisioninig
#        process.
#     2. The docker image created from test/dockerfiles/Dockerfile is already built
#        available on Arista's docker registry.
#     3. All server-to-server communication is done over SSH, and remote root login 
#        is performed by SSH host key authentication (RSA protocol used).
#     4. All playbooks that needs testing is under /playbooks and have been
#        appropriately added to the master_test.yml file.
#     5. An appropriate ansible_hosts file is filled out and ready in the production
#        environment. In this framework, an ad-hoc version is created with the docker
#        container information.
#
# Simulated Environment
# =====================
# This framework uses Docker containers to simulate the datacenter environment.
# The main agents are:
#
#        [ test ] ----> [ ansible server ] ----> [ server 1 ]
#                                          ----> [ server 2 ]
#                                          ----> [ server 3 ]
#                                                  ...
#                                          ----> [ server n ]
#
# The arrows in this diagram represent the direction of connection. Ansible server
# and the independent servers are represented by docker containers. The testing
# environment simulates a user logging into the ansible server to issue a command. 


# ================= Config info for this testing environment. =================
class Config( object ):
   '''
   Placeholder class for configuration settings for this testing environment 
   because python 2.7 doesn't natively support enum classes.

   '''
   # Number of client servers being simulated.
   num = 3

   debug = False
   if debug: 
      # Debug settings, referring to running tests on local MAC env.
      dockerImg = 'ar_fedora'
      ansible_serv = 'as_w'
      client_serv = 'sv%s_w'
   else:
      # Docker image built and available on Arista Docker registry.
      dockerImg = 'registry.docker.sjc.aristanetworks.com:5000/ardc-config:36d011e'

      # Ansible server name.
      ansible_serv = 'as'

      # Template for client server names.
      client_serv = 'sv%s'

   # Ansible server IP
   as_ip = ""

   # Dictionary of IP addresses to nodes
   servs = {}

   # Test-only 'ansible_hosts' template.
   ash_template = '[all]\n\n%s'

   # Test-only entry in 'ansible_hosts'.
   ash_host = '%s ansible_ssh_host=%s '

   # Test-only entry in 'known_hosts"
   kh_host = '%s,%s %s\n'

   # Path to RSA public host key.
   if debug:
      path_to_hostpub = 'dockerfiles/ar_fedora/ssh/id_rsa.pub'
   else:
      path_to_hostpub = 'test/dockerfiles/ar_fedora/ssh/id_rsa.pub'

   # Relative path to playbooks directory.
   if debug:
      playbooks = '../playbooks'
   else:
      playbooks = './playbooks'

   # Relative path to local.yml
   if debug:
      local_plbk = '../local.yml'
   else:
      local_plbk = './local.yml'
   local = 'local.yml'

   # Relative path to ping.yml
   if debug:
      ping_plbk = './ping.yml'
   else:
      ping_plbk = './test/ping.yml'
   ping = 'ping.yml'


class Cmd( object ):
   '''
   Placeholder class for often used commands for this testing environment because
   python 2.7 doesn't natively support enum classes.

   '''
   create = 'docker run -d -P -t --name %s %s > /dev/null'
   inspectIP = 'docker inspect --format \'{{.NetworkSettings.IPAddress}}\' %s'
   copy = 'docker cp %s %s'
   ex = 'docker exec -t %s /bin/bash -c "%s"'
   delete = 'rm %s'
   kill = 'echo %s | xargs -I %% sh -c "docker stop %%; docker rm %%" > /dev/null'
   ans_pl = 'ansible-playbook %s'


# ============================= Global Variables =================================
logging.basicConfig( level=logging.INFO )
logger = logging.getLogger( 'TestPlaybook' ) 
conf = Config()
cmd = Cmd()


# ========================== Ansible Push Test Class  ===========================
class TestAnsiblePushTests( unittest.TestCase ):

   def callcmd( self, command ): #pylint:disable-msg=R0201
      '''
      Wrapper for subprocess.call() to process exit code from running the command.

      Parameters
      ----------
      command : string
         The command (one of the members in Cmd) to run.

      '''
      ret = subprocess.call( command, shell=True )
      if ret:
         logger.error( 'Returned code %d: %s', ret, command )
      return ret


   def getcmd( self, command ): #pylint:disable-msg=R0201
      '''
      Wrapper for subprocess.check_output() to strip newline at the end of the
      returned output and return this output.

      Parameters
      ----------
      command : string
         The command (one of the members in Cmd) to run.
      '''
      return subprocess.check_output( command, shell=True ).strip( '\n' )


   def clean( self ): #pylint:disable-msg=R0201
      '''
      Clean up all spawned docker instances.
      '''
      logger.info( 'Killing all spawned docker containers for this test.' )
      containers = ' '.join( conf.servs.keys() ) 
      containers += ' %s' % conf.ansible_serv

      # Kill and remove all containers.
      self.callcmd( cmd.kill % containers )


   def setUp( self ):
      '''
      Set up main ansible server and client servers as docker instances for 
      this test that simulates actual datacenter environment.
      '''   
      # Try to clean environment first.
      # XXX: This isn't the cleanest way to do it since if these containers don't
      # already exist it would spew a lot of unwanted error messages to log. Also
      # doesn't delete all previous containers if number of containers for this run
      # is different. For now, this is good enough.
      logging.warning( '\nAttempting to clean up runaway containers.' )
      self.clean()
      
      # Create main ansible server. If ansible server couldn't be created, don't
      # even bother with rest of testing.
      self.callcmd( cmd.create % ( conf.ansible_serv, conf.dockerImg ) )
      conf.as_ip = self.getcmd( cmd.inspectIP % conf.ansible_serv ) 
      assert conf.as_ip
      logger.info( 'Initialized test Ansible host.' )
      # =========================================================================
      # XXX: ALL THE STUFF HERE SHOULD GO IN THE DOCKERFILE LATER.
      # Install Ansible 2.0
      logger.info( '%s: Remove Ansible 1.9 and install Ansible 2.0',
                   conf.ansible_serv )
      self.callcmd( cmd.ex % ( conf.ansible_serv, 'yum -y remove ansible' ) )
      self.callcmd( cmd.ex % ( conf.ansible_serv, 
         ( 'yum install -y https://archive.fedoraproject.org/pub/fedora/linux/'
           'updates/23/x86_64/a/ansible-2.0.2.0-1.fc23.noarch.rpm' ) ) )
      # =========================================================================

      # Create client servers.
      for i in range( conf.num ):
         curr = conf.client_serv % i
         ret = self.callcmd( cmd.create % ( curr, conf.dockerImg ) )
         if ret:
            logger.error( 'Couldn\'t initialize host %s. Skipping.', curr )
            continue

         conf.servs[ curr ] = self.getcmd( cmd.inspectIP % curr ) 
         assert conf.servs[ curr ]
         # =========================================================================
         # XXX: ALL THE STUFF HERE SHOULD GO IN THE DOCKERFILE LATER.
         # Install Ansible 2.0
         logger.info( '%s: Remove Ansible 1.9 and install Ansible 2.0', curr )
         self.callcmd( cmd.ex % ( curr, 'yum -y remove ansible' ) )
         self.callcmd( cmd.ex % ( curr, 
            ( 'yum install -y https://archive.fedoraproject.org/pub/fedora/linux/'
              'updates/23/x86_64/a/ansible-2.0.2.0-1.fc23.noarch.rpm' ) ) )
         # =========================================================================

      # Let's make sure that we don't have zero client servers ready to test on. It
      # would be pointless to run the test if we don't have any client servers. 
      assert len( conf.servs )
      logger.info( 'Initialized %d/%d test hosts.', len( conf.servs ), conf.num )

      # Ansible -- stop disconnecting SSH sessions on your own!
      self.callcmd( cmd.ex % ( conf.ansible_serv, 
         'sed -i \'s/#ssh_args = -o ControlMaster=auto -o ControlPersist=60s/'
         'ssh_args = -o ControlMaster=no/\' /etc/ansible/ansible.cfg' ) )

      # Create mock 'ansible_hosts' file for the pseudo-servers we just created.
      hosts = conf.ash_template % ( '\n'.join( [ conf.ash_host % ( sv, ip ) for sv, 
                                                 ip in conf.servs.items() ] ) )
      fd, tmp_host_f = tempfile.mkstemp()
      os.write( fd, hosts )
      os.close( fd )
      # SCP mock ansible hosts file over to as
      self.callcmd( cmd.copy % ( tmp_host_f, '%s:/test_ansible_hosts' % 
                                 conf.ansible_serv ) )
      os.unlink( tmp_host_f )

      # Create mock known_hosts file.
      with open( conf.path_to_hostpub, 'r' ) as f:
         key = f.read().strip('\n')
         known_hosts = ''.join( [ conf.kh_host % ( sv, ip, key ) for sv, ip 
                                  in conf.servs.items() ] )
         fd, tmp_kh_f = tempfile.mkstemp()
         os.write( fd, known_hosts )
         os.close( fd )
         # Copy over mock known_hosts file to ansible server container. 
         self.callcmd( cmd.copy % ( tmp_kh_f, '%s:/root/.ssh/known_hosts' 
                                               % conf.ansible_serv ) )
         os.unlink( tmp_kh_f )


   def tearDown( self ):
      '''
      Tearing down testing environment.
      '''
      self.clean()
      logging.shutdown()


   def test( self ):
      '''
      Runs master_test.yml, which is the master test playbook to run in order to run
      all other ansible playbooks to test for both test and production playbooks.
      '''   
      # Copy playbooks directory and its playbooks over to ansible server container.
      # XXX: This could be done by the Dockerfile, but wasn't in case the
      # directory/name/location of the playbooks changed.
      self.callcmd( cmd.copy % ( conf.playbooks, '%s:/' % conf.ansible_serv ) )

      # Copy over local.yml ansible server container.
      self.callcmd( cmd.copy % ( conf.local_plbk, '%s:/' % conf.ansible_serv ) )
   
      # Copy over ping.yml ansible server container.
      self.callcmd( cmd.copy % ( conf.ping_plbk, '%s:/' % conf.ansible_serv ) )

      # On Ansible server, run ping.yml as a test.
      ret = self.callcmd( cmd.ex % ( conf.ansible_serv, cmd.ans_pl % conf.ping ) )

      if ret:
         logger.error( 'Failed ping test. Aborting.' )
         sys.exit( 1 )

      # On Ansible server, run local.yml which plays all the playbooks.
      ret = self.callcmd( cmd.ex % ( conf.ansible_serv, cmd.ans_pl % conf.local ) )

      if ret:
         logger.error( 'Something went wrong playing the playbooks. Aborting.' )
         sys.exit( 1 )


if __name__ == '__main__':
   unittest.main( verbosity=2 )
