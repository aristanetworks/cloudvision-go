# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import os
import sys
import subprocess
import logging

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


# ==================== Config info for this testing environment. ====================
class Config( object ):
   '''
   Placeholder class for configuration settings for this testing environment because
   python 2.7 doesn't natively support enum classes.

   '''
   # Number of client servers being simulated.
   num = 10 

   #TODO connect to registry after code review.
   dockerImg = 'ar_fedora'

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
   ash_host = '%s ansible_ssh_host=%s ansible_ssh_pass=arista ansible_ssh_user=root'

   # Test-only 'ansible_hosts' file.
   ash_f = '/tmp/test_ansible_hosts'

   # Test-only entry in 'known_hosts"
   kh_host = '%s,%s %s\n'

   # Test-only 'known_hosts' file.
   kh_f = '/tmp/test_known_hosts'

   # Test-only main playbook to run.
   testplbk = 'master_test.yml'

   # Path to RSA public host key.
   path_to_hostpub = 'dockerfiles/ar_fedora/test_id_rsa.pub'

   # Relative path to playbooks directory.
   playbooks_dir = '../playbooks/'


class Cmd( object ):
   '''
   Placeholder class for often used commands for this testing environment because
   python 2.7 doesn't natively support enum classes.

   '''
   create = 'docker run -d -P --name %s %s'
   inspectIP = 'docker inspect --format \'{{.NetworkSettings.IPAddress}}\' %s'
   copy = 'docker cp %s %s'
   ex = 'docker exec -it as /bin/bash -c "%s"'
   delete = 'rm %s'
   kill = 'echo %s | xargs -I %% sh -c "docker stop %%; docker rm %%"'

# ================================ Helper Methods ===================================
def callcmd( command ):
   '''
   Wrapper for subprocess.call() to process exit code from running the command.

   Parameters
   ----------
   command : string
      The command (one of the members in Cmd) to run.

   '''
   ret = subprocess.call( command, shell=True )
   if ret:
      logger.error( 'FAILURE: returned code %d: %s' % ( ret, command ) )
   return ret


def getcmd( command ):
   '''
   Wrapper for subprocess.check_output() to strip newline at the end of the
   returned output and return this output.

   Parameters
   ----------
   command : string
      The command (one of the members in Cmd) to run.
   '''
   return subprocess.check_output( command, shell=True ).strip( '\n' )


def setup( conf, cmd ):
   '''
   Set up main ansible server and client servers as docker instances for 
   this test that simulates actual datacenter environment.

   Parameters
   ----------
   conf : Config object
      Contains configuration details for this test.

   cmd : Cmd object
      Set of hardcoded commands to use.

   '''   
   # Create main ansible server. If ansible server couldn't be created, don't
   # even bother with rest of testing.
   callcmd( cmd.create % ( conf.ansible_serv, conf.dockerImg ) )
   conf.as_ip = getcmd( cmd.inspectIP % conf.ansible_serv ) 
   assert conf.as_ip
   logger.info( 'Initialized test Ansible server.' )

   # Create client servers.
   for i in range( conf.num ):
      curr = conf.client_serv % i
      ret = callcmd( cmd.create % ( curr, conf.dockerImg ) )
      if ret:
         logger.error( 'FAILURE: Couldn\'t initialize node %s. Skipping.' % curr )
         continue

      conf.servs[ curr ] = getcmd( cmd.inspectIP % curr ) 
      assert conf.servs[ curr ]

   # Let's make sure that we don't have zero client servers ready to test on. It
   # would be pointless to run the test if we don't have any client servers. 
   assert len( conf.servs )
   logger.info( 'Initialized %d/%d test custom servers.' % ( len( conf.servs ), 
                                                             conf.num ) )

   # Create mock 'ansible_hosts' file for the pseudo-servers we just created.
   hosts = conf.ash_template % ( '\n'.join( [ conf.ash_host % ( sv, ip ) for sv, ip
                                              in conf.servs.items() ] ) )
   with open( conf.ash_f, 'w' ) as f:
      f.write( hosts )
      f.flush()
      
      # SCP mock ansible hosts file over to as
      callcmd( cmd.copy % ( conf.ash_f, 'as:/' ) )

      # Let's remove the temporary file we wrote.
      callcmd( cmd.delete % ( conf.ash_f ) )

   # Create mock known_hosts file.
   with open( conf.path_to_hostpub, 'r' ) as f:
      key = f.read().strip('\n')
      known_hosts = ''.join( [ conf.kh_host % ( sv, ip, key ) for sv, ip 
                               in conf.servs.items() ] )
      
      with open( conf.kh_f, 'w' ) as hf:
         hf.write( known_hosts )
         hf.flush()
         
         # Copy over mock known_hosts file to ansible server container. 
         callcmd( cmd.copy % ( conf.kh_f, 'as:/.ssh/known_hosts' ) )
         callcmd( cmd.copy % ( conf.kh_f, 'as:/etc/ssh/ssh_known_hosts' ) )

         # Let's remove the temporary file we wrote.
         callcmd( cmd.delete % ( conf.kh_f ) )


def runTest( conf, cmd ):
   '''
   Runs master_test.yml, which is the master test playbook to run in order to run
   all other ansible playbooks to test for both test and production playbooks.

   Parameters
   ----------
   conf : Config object
      Contains configuration details for this test.

   cmd : Cmd object
      Set of hardcoded commands to use.

   '''   
   # SCP playbooks to ansible server container. We have to copy at runtime instead
   # of creating a docker image with playbooks already copied over because the
   # number of playbooks tested may vary vastly at each test runtime.
   #
   # Change this if you want to run certain/different playbooks.
   playbooks = os.listdir( conf.playbooks_dir )
   for f in playbooks:
      callcmd( cmd.copy % ( os.path.join( os.path.abspath( conf.playbooks_dir ), f ), 
                            'as:/' ) )

   # From Ansible-server, run master test playbook that will run all other
   # playbooks.
   callcmd( cmd.ex % 'ansible-playbook master_test.yml' ) 


def clean( conf, cmd ):
   '''
   Clean up all spawned docker instances.

   Parameters
   ----------
   conf : Config object
      Contains configuration details for this test.

   cmd : Cmd object
      Set of hardcoded commands to use.

   '''
   logger.info( 'Killing all spawned docker containers for this test.' )
   containers = ' '.join( conf.servs.keys() ) 
   containers += ' %s' % conf.ansible_serv

   # Kill and remove all containers.
   callcmd( cmd.kill % containers )


# =============================== Main test runner ==================================
def main():
   # Set up logger for this test.
   logging.basicConfig( level=logging.INFO )
   global logger
   logger = logging.getLogger( 'TestPlaybook' ) 

   conf = Config()
   cmd = Cmd()

   setup( conf, cmd )
   runTest( conf, cmd )
   clean( conf, cmd )



if __name__ == '__main__':
   main()
