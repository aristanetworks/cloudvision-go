####################################################################################
#                                  ASBTrigger API                                  #
#                          (for organizational purposes)                           #
####################################################################################

import MySQLdb
import subprocess
import time
import shlex


# ASB Statuses
# XXX: This is a duplicate of group_vars/all
ASB_STATUS_FRESH = "fresh"
ASB_STATUS_REIMG_HEALTH = "reimage_healthcheck"
ASB_STATUS_REIMG_WIPE = "reimage_wipe"
ASB_STATUS_REIMG_VERIFY = "reimage_verify"
ASB_STATUS_DISK = "disk"
ASB_STATUS_OS = "os"
ASB_STATUS_CONF = "config"
ASB_STATUS_VERIFY = "verify"
ASB_STATUS_AUTH = "auth"
ASB_STATUS_INSTALLED = "installed"
ASB_STATUS_HELP = "HELP_ME"


def _connect():
   db = MySQLdb.connect( user="arastra", 
                         db="datacenter", 
                         host="mysql" )
   cs = db.cursor()
   return db, cs

def _disconnect( cursor, dbcon ):
   cursor.close()
   dbcon.close()


def fetch_servers( condition ):
   """
   Get the list of servers that meet CONDITION and are not reserved.
   """
   db, cs = _connect()

   # Select from datacenter.servers given condition
   stmt = "select name, domain from servers where %s" % ( condition )
   cs.execute( stmt )
   rows = list( cs.fetchall() )

   # Select all servers from datacenter.servers that is marked as reserved
   ## Note that LIKE is case insensitive
   cond = 'propertyName="reserved" AND propertyValue LIKE "true"'
   stmt2 = "SELECT serverName, domain FROM server_property WHERE %s" % ( cond )
   cs.execute( stmt2 )
   reserved_row = list( cs.fetchall() )

   _disconnect( cs, db )

   # Combine name with domain and create tuple to create full FQDN
   rows = [ "%s.%s" % ( n, d ) for n, d in rows ]
   reserved_row =  [ "%s.%s" % ( n, d ) for n, d in reserved_row ]

   # Ignore reserved servers
   servers = list( set( rows ) - set( reserved_row ) )
   return servers

def fetch_status( fqdn ):
   splits = fqdn.split( ".", 1 )
   name = splits[ 0 ]
   domain = splits[ 1 ]

   db, cs = _connect()
      
   # Get status of server by given name and domain
   stmt = "SELECT status FROM servers WHERE name='%s' AND domain='%s'" % ( name,
                                                                           domain )
   cs.execute( stmt )
   row = cs.fetchone()
   return row[ 0 ]


def trigger( servers ):
   # Trigger Command
   prod = "http://gerrit/ardc-config"
   trigger_cmd = ( 'sudo -i flock -xn /tmp/ASB.lck env GIT_SSL_NO_VERIFY=true ansible-pull '
               '--url=%s --directory=/root/.ansible/pull ASB/ASB.yml '
               '--inventory=localhost, --purge' ) % ( prod )

   # PDSH Command
   pdsh_cmd = ( "PDSH_SSH_ARGS_APPEND='-o StrictHostKeyChecking=no' pdsh -w %s %s"
                " | dshbak -d /var/log/ASBTrigger/hosts/" )

   if servers:
      cmd = pdsh_cmd % ( ",".join( servers ), trigger_cmd )

      # XXX this is a security hole
      subprocess.call( cmd, shell=True )
