#!/usr/bin/env python
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import unittest
import os
import shlex
import subprocess
import sys

# Paths on container emulating ticket server
# XXX: Maintenance issue, should import from ASBAuthenticator
WHITELIST_DIR = "/home/asb"
LOGFILE = "/var/log/ASBAuthenticator.log"

# Relative paths for files on localhost
TMP_LOCAL = "test/tmp"
LOGFILE_LOCAL = os.path.join( TMP_LOCAL, "ASBAuthenticator.log" )
AUTH_ROOT_DIR = "setup/authenticator"

# Strings for testing whitelist file management
WL_TODELETE = "Found %r outdated whitelist files to delete"
WL_NODELETE = "No outdated whitelist files found"
WL_DELETED = "Deleted expired whitelist file: %s"
WL_FAILED = "Failed to delete whitelist file: %s"


# NOTE: When this testing is expanded to other parts of ASB the Config and
# Cmd classes should most likely be refactored to be some sort of DockerClient
# class which takes the image and container name as constructor params.
class Config( object ):
   #TODO: @asetty update test image
   test_image = "asetty.asbauth:671e34b58856"
   container_name = "asbauth_testcontainer"


class Cmd( object ):
   '''
   Class used to issue docker commands used in testing
   '''
   copy = "docker cp %s %s"
   run = "docker run -d -t --name %s %s"
   ex = "docker exec -t %s /bin/bash -c \"%s\""
   stop = "docker stop %s"
   rm = "docker rm %s"

   @classmethod
   def call( cls, command, *args ):
      # Format command string with correct number of args
      # and print to console in order to follow execution during testing
      command = command % args
      print command

      cmd = shlex.split( command )
      ret = subprocess.call( cmd )
      if ret:
         sys.stderr.write( "ERROR: Docker CLI call terminated with error: %s." % \
                           ret )
      return ret

class AuthenticatorSanityTests( unittest.TestCase ):
   @classmethod
   def setUpClass( cls ):
      '''
      Set up test environment
      '''
      # Run test container, from basic image with Fedora20 with arastra
      # and asb user accounts. This container should be thought of
      # as ticketserver in our test environment
      ret = Cmd.call( Cmd.run, Config.container_name, Config.test_image )
      if ret:
         print "Unable to start docker container for image %s" % Config.test_image
         sys.exit( 1 )
      else:
         print "ASBAuthenticator test container running with name %s" % \
            Config.container_name

      # Copy over ASBAuthenticator and restore_server executables in the same place
      # as on ticketserver
      Cmd.call( Cmd.copy, os.path.join( AUTH_ROOT_DIR, "root/ASBAuthenticator" ),
                Config.container_name + ":/root" )
      # NOTE: restore_server is not yet used in testing
      Cmd.call( Cmd.copy, os.path.join( AUTH_ROOT_DIR, "root/restore_server" ),
                Config.container_name + ":/root" )


   @classmethod
   def tearDownClass( cls ):
      '''
      Clean up testing environment
      '''
      # Stop and remove the test container
      Cmd.call( Cmd.stop, Config.container_name )
      Cmd.call( Cmd.rm, Config.container_name )

   def runAuthenticator( self ):
      '''
      Run ASBAuthenticator executable on test container
      '''
      Cmd.call( Cmd.ex, Config.container_name, "/root/ASBAuthenticator --test" )

   def getLog( self ):
      '''
      Get ASBAuthenticator log from test container.
      By default, the log file copied over from the container is
      deleted, but can be kept with the delete parameter.
      Similarly, the log file and whitelist files on the test container
      are deleted by default, but can be kept via the clear parameter
      '''
      Cmd.call( Cmd.copy, Config.container_name + ":" + LOGFILE, TMP_LOCAL )
      try:
         f = open( LOGFILE_LOCAL, 'r' )
      except IOError as e:
         err = "Error opening log file copied from container: %s" % e.strerror
         sys.stderr.write( err )
         return None
      log = f.read()
      f.close()
      return log

   def clearContainerState( self ):
      '''
      Clear the log file for ASBAuthenticator and all files and directories
      in the whitelist directory on the test container
      '''
      Cmd.call( Cmd.ex, Config.container_name, "rm -f %s" % LOGFILE )
      Cmd.call( Cmd.ex, Config.container_name, "rm -rf %s" % \
                os.path.join( WHITELIST_DIR, "*" ) )

   def createFilesOnContainer( self, files, directory=False ):
      '''
      Create files on test container using the touch command

      Parameters:
         files - List of pairs (filename, date string) which specifies
            the name and age of the file to be created. For more
            information on the date string look at the touch man page, in
            particular the -d option.
         directory - Specify if files created should be directories.
            By default this is False and regular files are created
      '''
      for f, time in files:
         if directory:
            Cmd.call( Cmd.ex, Config.container_name, "mkdir %s" % \
                      os.path.join( WHITELIST_DIR, f) )
         Cmd.call( Cmd.ex, Config.container_name, "touch -d '%s' %s " % \
                   ( time, os.path.join( WHITELIST_DIR, f) ) )


   def testWhitelistExpiry( self ):
      def runAndCheck( expired, notExpired ):
         self.runAuthenticator()
         log = self.getLog()

         # Cleanup for next test
         self.clearContainerState()
         if not log:
            self.fail( "Log file expected, but none found" )
            return
         else:
            os.remove( LOGFILE_LOCAL )

         if len( expired ):
            self.assertIn( WL_TODELETE % ( len( expired ) ), log )
            self.assertNotIn( WL_NODELETE, log )
         else:
            self.assertIn( WL_NODELETE, log )
            self.assertNotIn( WL_TODELETE % ( len( expired ) ), log )

         # Check that each expired whitelist file is logged as being
         # deleted and that it was actually deleted
         # The opposite is checked for non-expired whitelist files,
         # and any other files which should not be deleted by
         # whitelist pruning (i.e. directories, non-FQDNs).
         for e, _ in expired:
            self.assertIn( WL_DELETED % e, log )
         for n, _ in notExpired:
            self.assertNotIn( WL_DELETED % n, log )

         # Whitelist file which should be deleted should never fail to be deleted
         self.assertNotIn( WL_FAILED, log )

      # Test when no files are expired or not expired
      # Equivalent of no whitelist files
      expired = []
      notExpired = []
      self.createFilesOnContainer( expired + notExpired )
      runAndCheck( expired, notExpired )

      # Test with two expired files
      expired = [ ( "test.aristanetworks.com", "24 hours ago" ),
                  ( "bs306.sjc.aristanetworks.com", "2 weeks ago" ) ]
      notExpired = []
      self.createFilesOnContainer( expired + notExpired )
      runAndCheck( expired, notExpired )

      # Test with only not expired files
      expired = []
      notExpired = [ ( "us121.sjc.aristanetworks.com", "1 hour ago" ),
                     ( "test1.lol.aristanetworks.com", "3 hours ago"),
                     ( "bs306.sjc.aristanetworks.com", "0 minutes ago" ) ]
      self.createFilesOnContainer( expired + notExpired )
      runAndCheck( expired, notExpired )

      # Test file which has just expired, along with one which is one minute away
      # Note: 1439 minutes is one minute until 24 hours
      expired = [ ( "test.aristanetworks.com", "24 hours ago" ),
                  ( "bs306.sjc.aristanetworks.com", "2 weeks ago" ) ]
      notExpired = [ ( "us121.sjc.aristanetworks.com", "1439 minutes ago" ) ]
      self.createFilesOnContainer( expired + notExpired )
      runAndCheck( expired, notExpired )

      # Test varying timed expired and non expired files
      expired = [ ( "test.aristanetworks.com", "2 years ago" ),
                  ( "test1.lol.aristanetworks.com", "1 week ago"),
                  ( "bs306.sjc.aristanetworks.com", "2 weeks ago" ) ]
      notExpired = [ ( "us121.sjc.aristanetworks.com", "1 hour ago" ),
                     ( "test2.lol.aristanetworks.com", "3 hours ago"),
                     ( "bs305.sjc.aristanetworks.com", "0 minutes ago" ) ]
      self.createFilesOnContainer( expired + notExpired )
      runAndCheck( expired, notExpired )
      
      # Test that files which are "expired", but not FQDNs are ignored
      nonFQDN = [ ( "bs306", "24 hours ago" ),
                  ( "bs306.sjc.aristanetworks", "1 week ago" ) ]
      self.createFilesOnContainer( nonFQDN )
      runAndCheck( [], nonFQDN )

      # Test that FQDN named directories which are "expired" are ignored
      directories = [ ( "test.aristanetworks.com", "2 years ago" ),
                  ( "test1.lol.aristanetworks.com", "1 week ago"),
                  ( "bs306.sjc.aristanetworks.com", "2 weeks ago" ) ]
      self.createFilesOnContainer( directories, directory=True )
      runAndCheck( [], directories )
