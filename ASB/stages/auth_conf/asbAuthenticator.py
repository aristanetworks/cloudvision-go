#!/usr/bin/env python
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

# ASB Authenticator: Wrapper script for "restore_server"
# ----------------------------------------------------------------------
# Collects FQDN placeholders in automation account which signals they are
# ready to be issued a ticket and issues ticket. Assumes this script is 
# initiated as root by cron job


import MySQLdb
import subprocess
import glob
import os
import sys
import logging


# Enable logging
logging.basicConfig( filename='/var/log/asbAuthenticator.log',
                     format='[ASB Authenticator] %(asctime)s: %(message)s',
                     datefmt='%m/%d/%Y %H:%M:%S',
                     level=logging.DEBUG )

# Database information
DB_INFO = { "user": "arastra",
            "db": "datacenter",
            "host": "mysql" }

# Datatable to read from
SERV_DB = "servers"

# ASB Status to look for
ASB_STATUS_AUTH = "auth"

# Dir with all the placeholders (whitelist files)
WHITELIST_DIR = "/home/asb/"

# Dir with all the user home directories
HOME_DIR = "/home"

# FQDN of all servers to be authenticated should contain this
CONTAINS_DOMAIN = "aristanetworks.com"

def main():
   db = MySQLdb.connect( user=DB_INFO[ "user" ],
                         db=DB_INFO[ "db" ],
                         host=DB_INFO[ "host" ] )
   cs = db.cursor()

   # Select all servers from the table that has status "verify_ready"
   cols = "name, domain"
   cond = 'status="%s"' % ( ASB_STATUS_AUTH )
   stmt = "select %s from %s where %s" % ( cols, SERV_DB, cond )

   cs.execute( stmt )
   rows = list( cs.fetchall() )

   if rows:
      # Full FQDN of servers that need ticket on datacenter.servers
      needTicket = [ ( "%s.%s" % ( name, domain ) ) for name, domain in rows ]

      # Get list of whitelisted FQDNs in ASB account
      whitelist = [ os.path.basename( p ) for p in glob.glob(
                    os.path.join( WHITELIST_DIR, "*" ) ) ]

      # Get list of servers already authenticated, ie. servers that has user account
      preAuthServers = os.listdir( HOME_DIR )
      preAuthServers = [ sv for sv in preAuthServers if CONTAINS_DOMAIN in sv ]

      # Only authenicate servers that are either:
      #    (1) whitelisted, ie. have a whitelisted file in ASB account, or
      #    (2) already authenticated, ie. have an user account on ticketserver.
      allowed = set( whitelist ).union( preAuthServers )
      authServers = set( needTicket ) & allowed

      for sv in authServers:
         ret = subprocess.call( [ "sh", "/root/restore_server", "%s" % sv ] )
         if ret:
            logging.error( "Could not issue ticket for %s" % sv )
         else:
            if sv in whitelist:
               logging.info( "Issued ticket for %s" % sv )
               # Remove the whitelist file after authentication
               logging.info( "Removed whitelist file for %s" % sv )
               os.remove( os.path.join( WHITELIST_DIR, sv ) )
            elif sv in preAuthServers:
               logging.info( "Re-issued ticket for %s" % sv )


if __name__ == "__main__":
   main()
