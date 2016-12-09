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

# ASB Status
ASB_STATUS_VERRDY = "verify_ready"

# Dir with all the placeholders
ALL_DIR = "/home/asb/"

def main():
   db = MySQLdb.connect( user=DB_INFO[ "user" ],
                         db=DB_INFO[ "db" ],
                         host=DB_INFO[ "host" ] )
   cs = db.cursor()
   
   # Select all servers from the table that has status "verify_ready"
   cols = "name, domain"
   cond = 'status="%s"' % ( ASB_STATUS_VERRDY )
   stmt = "select %s from %s where %s" % ( cols, SERV_DB, cond )

   cs.execute( stmt )
   rows = list( cs.fetchall() )

   if rows:
      # Full FQDN of servers that need ticket on datacenter.servers
      needTicket = [ ( "%s.%s" % ( name, domain ) ) for name, domain in rows ] 
      
      # Get list of FQDNs in ASB account
      allowed = [ os.path.basename( p ) for p in glob.glob( 
                  os.path.join( ALL_DIR, "*" ) ) ]

      authServers = set( needTicket ) & set( allowed )
      for sv in authServers:
         ret = subprocess.call( [ "sh", "/root/restore_server", "%s" % sv ] )
         if ret:
            logging.error( "Could not issue ticket for %s" % sv )
         else:
            logging.info( "Issued ticket for %s; removing from whitelist" % sv )

         os.remove( os.path.join( ALL_DIR, sv ) )


if __name__ == "__main__":
   main()
