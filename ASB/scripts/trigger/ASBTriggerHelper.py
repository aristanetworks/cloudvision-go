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
ASB_STATUS_REIMG_STANDBY = "reimage_standby"
ASB_STATUS_REIMG_SANITY = "reimage_sanity"
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
                         host="mysql-b1.cs2.aristanetworks.com" )
   cs = db.cursor()
   return db, cs

def _disconnect( cursor, dbcon ):
   cursor.close()
   dbcon.close()


def fetch_servers( condition, testmode=False ):
   db, cs = _connect()

   # Select from datacenter.servers given condition
   stmt = "select name, domain from servers where %s" % ( condition )
   cs.execute( stmt )
   rows = list( cs.fetchall() )

   # Select all servers from datacenter.servers that is marked for test use
   testcond = 'propertyName="misc" and propertyValue="TESTONLY"'
   stmt2 = "select serverName, domain from server_property where %s" % ( testcond ) 
   cs.execute( stmt2 )
   test_rows = list( cs.fetchall() )

   _disconnect( cs, db )

   # Combine name with domain and create tuple to create full FQDN
   rows = [ "%s.%s" % ( n, d ) for n, d in rows ]
   test_rows =  [ "%s.%s" % ( n, d ) for n, d in test_rows ]

   servers = list( set( rows ) - set( test_rows ) )
   if testmode:
      servers = list( set( rows ) & set( test_rows ) )

   return servers


def trigger( servers ):
   # Trigger Command
   prod = "http://gerrit/ardc-config"
   trigger_cmd = ( 'sudo -i flock -xn /tmp/ASB.lck env GIT_SSL_NO_VERIFY=true ansible-pull '
               '--url=%s --directory=/root/.ansible/pull ASB/ASB.yml '
               '--inventory=localhost,' ) % ( prod )

   # PDSH Command
   pdsh_cmd = ( "PDSH_SSH_ARGS_APPEND='-o StrictHostKeyChecking=no' pdsh -w %s %s"
                " | dshbak -d /var/log/ASBTrigger-hosts/" )

   if servers:
      cmd = pdsh_cmd % ( ",".join( servers ), trigger_cmd )	

      # XXX this is a security hole
      subprocess.call( cmd, shell=True )

