#!/usr/bin/env python
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import os
import sys
import glob
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

   # Debug settings, referring to running tests on local env.
   DEBUG = False

   # Docker image built and available on Arista Docker registry.
   dockerImg = 'registry.docker.sjc.aristanetworks.com:5000/ardc-config:9b1226b'

   # Ansible server name.
   ansible_sv = 'ardc_config_as'

   # Template for client server names.
   client_serv = 'ardc_config_sv%s'

   # Package directory
   pkgdir = "AroraConfig"

   # Misc files directory 
   misc_files = "AroraConfig_testfiles"
   
   # Path to RSA public host key.
   path_to_hostpub = 'test/dockerfiles/ar_fedora/ssh/id_rsa.pub'

   # Relative path to playbooks directory.
   playbooks = '%s/playbooks/' % pkgdir

   # Main master playbook to be run.
   master = 'AroraConfig.yml'

   # Relative path to master playbook
   master_plbk = '%s/%s' % ( pkgdir, master )

   # Relative path to ping.yml
   ping_plbk = 'test/%s/ping.yml' % misc_files

   # Test-only ping playbook
   ping = 'ping.yml'

   # Expected path to sentinel files
   sentinels_path = '/var/lib/AroraConfig/'

   # Destination path to put Ansible hosts file.
   dest_ansible_hosts = '/etc/ansible/hosts'

   # Destination path to put Ansible config file.
   as_cfg = '/etc/ansible/ansible.cfg'

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

   # Copy of Arora.conf from Arora18release 
   arora_conf = 'test/%s/Arora2010.conf_COPY' % misc_files

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
   # We can't just kill and remove every container because there are other containers
   # running on this same node as where the test runs.
   kill = ( 'echo %s | xargs -I %% sh -c "docker kill %%; docker rm %%"' 
            '> /dev/null 2>&1' )
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
         logger.error( 'Docker exec call terminated with error %s.', ret )
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


   def clean( self, preclean=False ): #pylint:disable-msg=R0201
      '''
      Clean up all spawned docker instances.
      '''
      logger.info( 'Killing all spawned docker containers for this test.' )

      if preclean:
         servs = [ conf.client_serv % n for n in range( conf.num ) ]
         containers = ' '.join( servs )
      else:
         containers = ' '.join( conf.servs.keys() ) 
      containers += ' %s' % conf.ansible_sv

      # Kill and remove all containers.
      ret = self.callcmd( cmd.kill % containers )
      if ret:
         if preclean:
            logging.warning( 'No runaway containers to remove.' )
         else:
            logging.error( 'Couldn\'t remove containers.' )


   def modifyAnsibleServerConfigs( self ): #pylint:disable-msg=R0201
      '''
      Some manual configuration settings specific for server. This should be in line
      with config.yml. This is the workaround of not keeping our a centralized copy
      of ansible.cfg, but instead using a playbook to apply the config changes.

      '''
      # XXX: If something is really different from deployment env then the settings
      # here are probably very different from the deployment env.
      # XXX: If something is really broken, check that the settings changed here are
      # the same ones as those in config.yml.

      # Disable host key checking
      self.callcmd( cmd.ex % ( conf.ansible_sv,
                    'sed -i \'/host_key_checking/c\host_key_checking = False\' %s' %
                    conf.as_cfg ) )

      # SSH args - stop auto disconnects after timeout
      self.callcmd( cmd.ex % ( conf.ansible_sv,
         'sed -i \'/ssh_args/c\ssh_args = -o ControlMaster=no\' %s' % conf.as_cfg ) )

      # Enable pipelining
      self.callcmd( cmd.ex % ( conf.ansible_sv,
         'sed -i \'/pipelining/c\pipelining = True\' %s' % conf.as_cfg ) )


   def checkPlaybooks( self ): #pylint:disable-msg=R0201
      '''
      Runs coarse sanity check against the playbooks to see if they meet the minimum
      requirements for testing. This is more to coerce tests to be written regardless
      of the test's actual content.

      '''

      # We want people to at least *try* to use the testing framework before their 
      # playbook gets added to the repo.
      # The structure we are looking are is:
      #           - name: "TEST"
      #             script: test/test_PLAYBOOK.py
      #             when: TEST is defined
      #
      # XXX: This testing would be much easier if we actually wrote a custom Ansible
      # module only for testing that people can easily include.
      # XXX: This does not guarantee that these lines follow one immediately after
      # another, but guarantees the order of the lines and the fact that these lines
      # do exist in the file.
      playbooks = glob.glob( os.path.join( conf.playbooks, '*.yml' ) )
      if not playbooks:
         logger.error( 'No playbooks found to sanity check.' )
         sys.exit( 1 )

      for play in playbooks:
         with open( play, 'r' ) as pl: 

            expectedLns = []
            expectedLns.append( '"TEST"' )
            expectedLns.append( 'script: test/' )
            expectedLns.append( 'when: TEST is defined' ) 

            for line in pl:
               if not expectedLns:
                  break

               line = line.lstrip()
               if not line:
                  continue
               elif line[0] == '#':
                  continue
               elif expectedLns[ 0 ] in line:
                  expectedLns.pop( 0 )
            
            if expectedLns:
               logger.error( 'No testing found for playbook "%s".', play )
               sys.exit( 1 )
            else:
               logger.info( 'Looks like there is testing in playbook "%s".', play )



   def thingsThatShouldBeInDockerImage( self ):
      # XXX: All the stuff here should really be in the Docker image, not here!!
      # But because changing Dockerfile is such a pain let's collect them here...
      # XXX: All the stuff here are TEMPORARY. 
      conts = list( conf.servs.keys() )
      conts.append( conf.ansible_sv )

      # Copy Arora.conf over to the containers. Need to discard Arora.conf anyway.
      logger.info( "Copying Arora.conf from Arora18release to containers." )
      for c in conts:
         self.callcmd( cmd.copy % ( conf.arora_conf,
                       '%s:%s' % ( c, '/etc/sysctl.d/Arora.conf' ) ) )


   def setUp( self ):
      '''
      Set up main ansible server and client servers as docker instances for 
      this test that simulates actual datacenter environment.

      '''   
      logger.info( 'Starting tests...' )

      # ========== SANITY CHECK PLAYBOOKS ==========
      # Sanity check that playbooks have testing framework included.
      logger.info( 'Sanity checking playbooks for tests.' )
      self.checkPlaybooks()


      # ========== TRY TO CLEAN UP ENV ==========
      # This isn't the cleanest way to do it since if these containers don't
      # already exist it would spew a lot of unwanted error messages to log. Also
      # doesn't delete all previous containers if number of containers for this run
      # is different. For now, this is good enough.
      logging.warning( 'Attempting to remove runaway containers from previous run.' )
      self.clean( preclean=True )


      # ========== CREATE SERVERS ==========
      # Create main ansible server. If ansible server couldn't be created, don't
      # even bother with rest of testing.
      self.callcmd( cmd.create % ( conf.ansible_sv, conf.dockerImg ) )
      conf.as_ip = self.getcmd( cmd.inspectIP % conf.ansible_sv ) 
      assert conf.as_ip
      logger.info( 'Initialized test Ansible host.' )

      # Create client servers.
      for i in range( conf.num ):
         curr = conf.client_serv % i
         ret = self.callcmd( cmd.create % ( curr, conf.dockerImg ) )
         if ret:
            logger.error( 'Couldn\'t initialize host %s. Skipping.', curr )
            continue

         conf.servs[ curr ] = self.getcmd( cmd.inspectIP % curr ) 
         assert conf.servs[ curr ]


      # Let's make sure that we don't have zero client servers ready to test on. It
      # would be pointless to run the test if we don't have any client servers. 
      assert conf.servs
      logger.info( 'Initialized %d/%d test hosts.', len( conf.servs ), conf.num )


      # ========== SET UP ANSIBLE SERVER ==========
      # Create mock 'ansible_hosts' file for the pseudo-servers we just created.
      logger.info( 'Setting Ansible configuration settings on ansible server.' )
      self.modifyAnsibleServerConfigs()

      logger.info( 'Creating and copying over ansible_hosts file to as.' )
      hosts = conf.ash_template % ( '\n'.join( [ conf.ash_host % ( sv, ip ) for sv, 
                                                 ip in conf.servs.items() ] ) )
      fd, tmp_host_f = tempfile.mkstemp()
      os.write( fd, hosts )
      os.close( fd )
      # SCP mock ansible hosts file over to as, where it is expected to be read:
      # '/etc/ansible/hosts'
      self.callcmd( cmd.copy % ( tmp_host_f,
                    '%s:%s' % ( conf.ansible_sv, conf.dest_ansible_hosts ) ) )
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
         #self.callcmd( cmd.copy % ( tmp_kh_f, '%s:/root/.ssh/known_hosts' 
         #                                      % conf.ansible_sv ) )
         os.unlink( tmp_kh_f )


      # ========== SET UP HOST SERVERS ==========
      # Get a list of playbooks available. These are the names of the sentinels,
      # without the '.yml' extension.
      yamlFiles  = os.listdir( conf.playbooks )
      sentinels = [ os.path.splitext( f )[ 0 ] for f in yamlFiles ]
      
      # Create sentinel files on the host servers.
      logger.info( 'Creating sentinel files on host servers.' )
      for sv,_ in conf.servs.iteritems():
         for st in sentinels:
            self.callcmd( cmd.ex % ( sv, 'mkdir -p %s' % conf.sentinels_path ) )
            self.callcmd( cmd.ex % ( sv, 'touch %s%s' % 
                                         ( conf.sentinels_path, st ) ) )

      self.thingsThatShouldBeInDockerImage()
   

   def tearDown( self ):
      '''
      Tearing down testing environment.
      '''
      self.clean()
      logging.shutdown()


   def test( self ):
      '''
      Runs AroraConfig.yml, which is the master test playbook to run in order to run
      all other AroraConfig related ansible playbooks to test for both test and 
      production playbooks.
      '''   

      # Copy playbooks directory and its playbooks over to ansible server.
      # This could be done by the Dockerfile, but wasn't in case the 
      # directory/name/location of the playbooks changed.
      logger.info( 'Copying playbooks to Ansible server.' )
      self.callcmd( cmd.copy % ( conf.playbooks, '%s:/' % conf.ansible_sv ) )

      # Copy over master playbook to ansible server container.
      self.callcmd( cmd.copy % ( conf.master_plbk, '%s:/' % conf.ansible_sv ) )
   
      # Copy over ping.yml ansible server container.
      self.callcmd( cmd.copy % ( conf.ping_plbk, '%s:/' % conf.ansible_sv ) )

      # On Ansible server, run ping.yml as a test.
      logger.info( 'Running sanity ping test on the containers.' )
      ret = self.callcmd( cmd.ex % ( conf.ansible_sv, cmd.ans_pl % conf.ping ) )

      if ret:
         logger.error( 'Failed ping test. Aborting.' )
         sys.exit( 1 )

      # Run all the playbooks once. This simulates provisioning the servers from
      # a blank slate state.
      logger.info( 'Test iteration 1: simulating running playbooks on blank slate '
                   'servers, newly provisioned.' )
      ret1 = self.callcmd( cmd.ex % ( conf.ansible_sv, cmd.ans_pl % conf.master ) )

      # Run all the playbooks again. This simulates maintaining the servers that are
      # already provisioned.
      logger.info( 'Test iteration 2: simulating maintaining servers already '
                   'provisioned.' )
      ret2 = self.callcmd( cmd.ex % ( conf.ansible_sv, cmd.ans_pl % conf.master ) )

      if ret1 or ret2:
         if ret1:
            # Running the playbooks the first time breaks things.
            logger.error( 'Error with Test iteration 1.' )
         if ret2:
            # Trying to maintain with the playbooks breaks things.
            logger.error( 'Error with Test iteration 2.' )
         if ret1 and not ret2:
            # Playing the playbooks the first time broke things, but running the
            # playboks again fixed the problem itself.
            logger.error( 'Error with Test iteration 1, but iteration 2 fixed it.' )
         if not ret1 and ret2:
            logger.error( 'Test iteration 2 broke test iteration 1.' )

         sys.exit( 1 )
      else:
         logger.info( 'Both test iterations passed without error.' )


if __name__ == '__main__':
   unittest.main( verbosity=2 )
   print '\n'
