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
   num = 5

   # Debug settings, referring to running tests on local MAC env.
   DEBUG = False

   # Docker image built and available on Arista Docker registry.
   dockerImg = ( 'ar_fedora' if DEBUG else
                 'registry.docker.sjc.aristanetworks.com:5000/ardc-config:36d011e' )

   # Ansible server name.
   ansible_serv = 'as_w' if DEBUG else 'as'

   # Template for client server names.
   client_serv = 'sv%s_w' if DEBUG else 'sv%s'

   # Path to RSA public host key.
   path_to_hostpub = ( 'dockerfiles/ar_fedora/ssh/id_rsa.pub' if DEBUG else 
                       'test/dockerfiles/ar_fedora/ssh/id_rsa.pub' )

   # Relative path to playbooks directory.
   playbooks = '../playbooks/' if DEBUG else './playbooks/'

   # Relative path to local.yml
   local_plbk = '../local.yml' if DEBUG else './local.yml'

   # Relative path to ping.yml
   ping_plbk = './ping.yml' if DEBUG else './test/ping.yml'

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

   # Main master playbook to be run.
   local = 'local.yml'

   # Test-only ping playbook
   ping = 'ping.yml'

   # Expected path to sentinel files
   sentinels_path = '/var/lib/AroraConfig/'

class Cmd( object ):
   '''
   Placeholder class for often used commands for this testing environment because
   python 2.7 doesn't natively support enum classes.

   '''
   create = 'docker run -d -P -t --name %s %s > /dev/null'
   inspectIP = 'docker inspect --format \'{{.NetworkSettings.IPAddress}}\' %s'
   copy = 'docker cp %s %s'
   ex = 'docker exec -t %s /bin/bash -c "%s"'
   ex_no_output = 'docker exec -t %s /bin/bash -c "%s" > /dev/null'
   delete = 'rm %s'
   kill = 'echo %s | xargs -I %% sh -c "docker stop %%; docker rm %%" > /dev/null'
   ans_pl = 'ansible-playbook %s -e "TEST=true"'


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


   # XXX: MOVE THIS TO DOCKERFILE ASAP
   def stuffThatReallyShouldBeInTheDockerFile( self, node ):
      '''
      Collecting changes that should go to Dockerfile. Really shouldn't be here.
      ALL THE STUFF HERE SHOULD GO IN THE DOCKERFILE LATER!!!!!!
      '''
      # Install Ansible 2.0
      logger.info( '%s: Remove Ansible 1.9', node )
      self.callcmd( cmd.ex_no_output % ( node, 'yum -y remove ansible' ) )
      logger.info( '%s: Install Ansible 2.0', node )
      self.callcmd( cmd.ex_no_output % ( node, 
         ( 'yum install -y https://archive.fedoraproject.org/pub/fedora/linux/'
           'updates/23/x86_64/a/ansible-2.0.2.0-1.fc23.noarch.rpm' ) ) )


   def setAnsibleConfig( self ): #pylint:disable-msg=R0201
      '''
      Centralized function to set custom Ansible Server configuration settings.
      '''
      # Ansible -- stop disconnecting SSH sessions on your own!
      self.callcmd( cmd.ex % ( conf.ansible_serv, 
         'sed -i \'s/#ssh_args = -o ControlMaster=auto -o ControlPersist=60s/'
         'ssh_args = -o ControlMaster=no/\' /etc/ansible/ansible.cfg' ) )


   def setUp( self ):
      '''
      Set up main ansible server and client servers as docker instances for 
      this test that simulates actual datacenter environment.
      '''   
      # ========== TRY TO CLEAN UP ENV ==========
      # This isn't the cleanest way to do it since if these containers don't
      # already exist it would spew a lot of unwanted error messages to log. Also
      # doesn't delete all previous containers if number of containers for this run
      # is different. For now, this is good enough.
      logging.warning( '\nAttempting to clean up runaway containers.' )
      self.clean()
      

      # ========== CREATE SERVERS ==========
      # Create main ansible server. If ansible server couldn't be created, don't
      # even bother with rest of testing.
      self.callcmd( cmd.create % ( conf.ansible_serv, conf.dockerImg ) )
      conf.as_ip = self.getcmd( cmd.inspectIP % conf.ansible_serv ) 
      assert conf.as_ip
      logger.info( 'Initialized test Ansible host.' )
      self.stuffThatReallyShouldBeInTheDockerFile( conf.ansible_serv ) # XXX

      # Create client servers.
      for i in range( conf.num ):
         curr = conf.client_serv % i
         ret = self.callcmd( cmd.create % ( curr, conf.dockerImg ) )
         if ret:
            logger.error( 'Couldn\'t initialize host %s. Skipping.', curr )
            continue

         conf.servs[ curr ] = self.getcmd( cmd.inspectIP % curr ) 
         assert conf.servs[ curr ]
         self.stuffThatReallyShouldBeInTheDockerFile( curr ) # XXX

      # Let's make sure that we don't have zero client servers ready to test on. It
      # would be pointless to run the test if we don't have any client servers. 
      assert len( conf.servs )
      logger.info( 'Initialized %d/%d test hosts.', len( conf.servs ), conf.num )


      # ========== SET UP ANSIBLE SERVER ==========
      # Set custom Ansible Server (as) configurations.
      self.setAnsibleConfig()

      # Create mock 'ansible_hosts' file for the pseudo-servers we just created.
      logger.info( 'Creating and copying over ansible_hosts file to as.' )
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
      logger.info( 'Creating and copying over known_hosts file to as.' )
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

      # ========== SET UP HOST SERVERS ==========
      # Get a list of playbooks available. These are the names of the sentinels,
      # without the '.yml' extension.
      yamlFiles  = os.listdir( conf.playbooks )
      sentinels = [ os.path.splitext( f )[ 0 ] for f in yamlFiles ]
      
      # Create sentinel files on the host servers.
      logger.info( 'Creating sentinel files on host servers.' )
      for sv in conf.servs.keys():
         for st in sentinels:
            self.callcmd( cmd.ex % ( sv, 'mkdir -p %s' % conf.sentinels_path ) )
            self.callcmd( cmd.ex % ( sv, 'touch %s%s' % 
                                         ( conf.sentinels_path, st ) ) )


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
      # This could be done by the Dockerfile, but wasn't in case the 
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
