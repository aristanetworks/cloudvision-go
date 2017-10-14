# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import sys
import BaseHTTPServer
import MySQLdb
import logging
import argparse

# MySQL constants
ARDC_HOST = "ardcdb"
ARDC_USER = "arastra"
ARDC_DB = "datacenter"
ARDC_TABLE_SERVERS = "servers"
ARDC_TABLE_SERVER_PROPERTY = 'server_property'

def _getMySQLConnection( host=ARDC_HOST, user=ARDC_USER, db=ARDC_DB ):
   try:
      db = MySQLdb.connect( user=user, db=db, host=host )
   except MySQLdb.error as e:
      logging.error( "Unable to connect to MySQL: %s", str( e ) )
      return None
   return db


def getServer( name, domain="sjc.aristanetworks.com" ):
   '''
   Return a dictionary containing info of the specified server on datacenter.db
   '''
   logging.info( "Fetching info for server with name %s and domain %s...",
                 name, domain )
   serverInfo = {}
   db = _getMySQLConnection()
   if not db:
      logging.error( "Cannot fetch server with name %s and domain %s.",
            name, domain )
      return None
   cursor = db.cursor()

   stmt = "SELECT role FROM %s WHERE name = %%s AND domain = %%s" %\
            ARDC_TABLE_SERVERS
   cursor.execute( stmt, ( name, domain ) )
   info = cursor.fetchone()
   if not info:
      logging.error( "Cannot fetch server with name %s and domain %s: server not "
            "found", name, domain )
      cursor.close()
      db.close()
      return None
   role = info[ 0 ]

   stmt = "SELECT propertyValue from %s WHERE serverName = %%s AND domain = %%s "\
         "AND propertyName = 'numDisks'" % ARDC_TABLE_SERVER_PROPERTY
   cursor.execute( stmt, ( name, domain ) )
   info = cursor.fetchone()
   if not info:
      logging.error( "Cannot fetch server with name %s and domain %s: server "
            "does not have property 'numDisks'", name, domain )
      cursor.close()
      db.close()
      return None
   numDisks = int( info[ 0 ] )

   serverInfo[ "name" ] = name
   serverInfo[ "domain" ] = domain
   serverInfo[ "role" ] = role
   serverInfo[ "numDisks" ] = numDisks
   cursor.close()
   db.close()
   logging.info( "Fetched server: %s", serverInfo )
   return serverInfo

def getDiskScheme( role, numDisks ):
   '''
   Return the disk+partition scheme for the given server role and number of disks.
   '''
   # TODO: Query the scheme from datacenter db and convert it to the syntax used
   #       by anaconda kickstart.
   logging.info( "Fetching disk and partition scheme..." )
   import os
   DIR_PATH = "/root"
   filePath = DIR_PATH + "/%s-%d-disks" % ( role, numDisks )
   if os.path.exists( filePath ):
      with open( filePath ) as f:
         data = f.read()
         logging.info( "Fetched scheme for %s server with %d disks",
                       role, numDisks )
         logging.info( "Scheme is:\n%s", data )
         return data
   else:
      logging.error( "Cannot fetch scheme: %s server with %d disks is not "
                     "supported.", role, numDisks )
      return ""

class DiskHTTPRequestHandler( BaseHTTPServer.BaseHTTPRequestHandler ):
   def do_GET( self ):
      '''
      Respond to GET request with the appropriate disk+partition scheme.
      '''
      fqdn = self.address_string()
      logging.info( "Serving GET request from %s %s...", fqdn,
                     str( self.client_address ) )
      name, domain = fqdn.split( ".", 1 )
      serverInfo = getServer( name, domain )
      if not serverInfo:
         logging.info( "Unable to serve GET request from %s %s, sending 404 "
                       "response", fqdn, str( self.client_address ) )
         self.send_error( 404 )
         return
      diskScheme = getDiskScheme( serverInfo[ "role" ], serverInfo[ "numDisks" ] )
      if diskScheme:
         logging.info( "Successfully served GET request from %s %s", fqdn,
                       str( self.client_address ) )
         res = ""
         res += diskScheme
         self.wfile.write( res )
      else:
         logging.info( "Unable to serve GET request from %s %s, sending 404 "
                       "response", fqdn, str( self.client_address ) )
         self.send_error( 404 )

   def log_message( self, format, *args ):
      msg = "%s - - [%s] %s" % ( self.client_address[ 0 ],
                                 self.log_date_time_string(),
                                 format % args )
      logging.info( "HTTP server message:  %s", msg )

def run( server_class=BaseHTTPServer.HTTPServer, handler=DiskHTTPRequestHandler,
         port=8080 ):
   try:
      server_address = ( '', port )
      httpd = server_class( server_address, handler )
      logging.info( "Starting disk server, listening on port %d...", port )
      httpd.serve_forever()
   except Exception as e:
      logging.error( "Error in running disk server: %s", str( e ) )
      logging.info( "Aborting." )
      sys.exit( 1 )
   finally:
      if httpd:
         # Clean up the server
         logging.info( "Shutting down disk server..." )
         httpd.server_close()
         logging.info( "Shut down disk server." )

def createParser():
   parser = argparse.ArgumentParser( description=( """Start a http server to \
service GET request for disk+partition scheme from servers undergoing anaconda \
kickstart.""" ), formatter_class=argparse.RawTextHelpFormatter )
   parser.add_argument( '--verbose', action='store_true',
                        help="Enable verbose output" )
   parser.add_argument( '--port', type=int, default=8080,
                        help="The port this server will listen to." )
   return parser

if __name__ == "__main__":
   # Parse command line arguments
   ps = createParser()
   args = ps.parse_args()

   # Setup logging
   if args.verbose:
      level = logging.DEBUG
   else:
      level = logging.INFO
   logging.basicConfig( filename="/var/log/DiskServer.log",
                        format='%(asctime)s: %(message)s',
                        datefmt='%m/%d/%Y %H:%M:%S',
                        level=level )
   # Start the server
   run( port=args.port )

