#!/usr/bin/env python
#Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

import collections
import subprocess
import sys
import smtplib
import datetime
import argparse
import shlex
import re
import csv
import os

from email import encoders
from email.mime.base import MIMEBase
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart


# ----------------------------------------------------------------------------------
# Constants
# ----------------------------------------------------------------------------------

# Email template for account reassignment
REASSIGN_HDR = "Time for a New Home! (new server assignment)"
REASSIGN_MSG = '''
Hi,

Your current server is suffering from btrfs related issues. 
We have noticed that you have been using your current server diligently, so you have 
been specially selected to transcend to a fresh new server for a much better 
experience!

Your current server: %s
Your new server: %s

Please follow these steps to transfer:   
   On your new server:
      1. Login as arastra: a4 ssh %s
      2. Add a new account for yourself: a4 account add <your username>
      3. Add yourself to /etc/motd: sudo vi /etc/motd
      4. Recreate workspaces and transfer over any files you kept

   On your current server:
      1. Login to %s using your username
      2. Submit all your pending changes: a4 submit
      3. Delete all workspaces: a4 deletetree -f <dir> (repeat for dirs)
      4. Delete any swi workspaces: swi workspace -d <dir>
      5. Copy over any other remaining files you wish to keep
      6. Logout
      7. Login to %s as arastra: a4 ssh %s
      8. Remove your home directory: sudo rm -rf /home/<your username>
      9. Release your account: a4 account release <your username>
  
Thanks,
Kind hearts of SW-INFRA
'''

# Email template for account delete
DELETE_HDR = "Please delete your infrequently used account on %s"
DELETE_MSG = '''
Hi,

There are too many users on your current server than what is preferable for the 
overall load for the server. We have noticed that you have not logged onto your 
current server for the past month and your account on there is looking abandoned. 
Please carefully delete the unused account on your current server, making any 
backups as needed.

Your current server: %s

Please follow these steps to delete:   
   On your current server:
      1. Login to %s using your username
      2. Submit all your pending changes: a4 submit
      3. Delete all workspaces: a4 deletetree -f <dir> (repeat for dirs)
      4. Delete any swi workspaces: swi workspace -d <dir>
      5. Copy over any other remaining files you wish to keep
      6. Logout
      7. Login to %s as arastra: a4 ssh %s
      8. Remove your home directory: sudo rm -rf /home/<your username>  
      9. Release your account: a4 account release <your username>

Thanks,
Kind hearts of SW-INFRA
'''

# Email template for reassignment report
REPORT_HDR = "Server Reassignment Report"

# Default Email Sender
DEFAULT_SENDER = "sw-infra-support@arista.com"

INFRA_MAIL_MSG = """
Attached is auto-generated reassignment report in csv format.
Please import file to a Google Drive Excel sheet. 
Users mentioned in the report have been noticed via email.
"""

INFRA_MAIL_REPORT_ONLY_MSG = """
Attached is auto-generated reassignment report in csv format.
Please import file to a Google Drive Excel sheet. 

NOTE: This is report only. No users have been notified of this reassignment.
"""


# ----------------------------------------------------------------------------------
# Logic
# ----------------------------------------------------------------------------------
def sendEmail( receiver, hdr, content, sender=DEFAULT_SENDER, attach="" ):
   '''
   Send email.

   Parameters
   ----------
   receiver: string
      Email of the user that is migrating

   hdr: string
      Email header

   content: string
      Body of email

   sender: string
      By default, send emails from sw-infra@arista.com

   attach: string
      Path to file attachment
   '''
   email = MIMEMultipart( )
   email[ 'Subject' ] = hdr
   email[ 'From' ] = sender
   email[ 'To' ] =  receiver

   email.attach( MIMEText( content, 'plain') )

   if attach:
      with open( attach ) as f:
         csv = MIMEBase( 'application', 'octet-stream' )
         csv.set_payload( f.read() )
         f.close()
         encoders.encode_base64( csv )
         csv.add_header( 'Content-Disposition', 'attachment', filename=attach )
         email.attach(csv)

   s = smtplib.SMTP()
   s.connect()
   s.sendmail( sender, receiver, email.as_string() )
   s.quit()
   print "Email sent to %s" % receiver



def notify( new_assign, del_assign, skipped, assignOnly ):
   '''
   Send emails notifying users to reassign to a different server, or delete their
   old unused accounts.

   Parameters
   ----------
   new_assign: dictionary of the format -

      { newserver : [ ( user, oldserver, email ), ... ] }

      Where newserver is name of new server to migrate to, user is username,
      oldserver is name of old server to migrate from, and email is the user's listed
      email on a4 users. Represents users to be migrated.

   del_assign: list of the format -

      [ ( user, oldserver, email ), ... ]

      Items are similar to the description for new_assign. Represents users to be
      deleted.

   skipped: list of the format -

      [ ( user, oldserver ), ... ]

      Represents skipped users because their username does not exist in "a4 users"
      output.

   assignOnly: boolean
      If True, send assignment report only to SW-INFRA-SUPPORT.
      If False, send assignment report to SW-INFRA-SUPPORT and also email the 
         affected users.
   '''
   report = []

   # Email reassignment to affected users
   for newassignment in new_assign.items():
      newserver, users = newassignment
      report.append( "Accounts Reassigned to %s" % newserver )
      report.append( "User, Previous Server, New Server, Email" )
      for entry in users:
         user, oldserver, email = entry
         body = REASSIGN_MSG % ( oldserver, 
                                 newserver, 
                                 newserver,
                                 oldserver,
                                 oldserver,
                                 oldserver )
         if not assignOnly:
            sendEmail( email, REASSIGN_HDR, body )
         report.append( "%s, %s, %s, %s" % ( user, oldserver, newserver, email ) )
   report.append( "\n" )

   # Email deletions to affected users
   report.append( "Accounts Deleted" )
   report.append( "User, Previous Server, Email" )
   for entry in del_assign:
      user, oldserver, email = entry
      body = DELETE_MSG % ( oldserver, 
                            oldserver, 
                            oldserver, 
                            oldserver )
      if not assignOnly:
         sendEmail( email, DELETE_HDR % oldserver, body )
      report.append( "%s, %s, %s" % ( user, oldserver, email ) )
   report.append( "\n" )

   #Add skipped users to report
   report.append( "Skipped Users due to Missing User Info" )
   report.append( "User, Previous Server" )
   for entry in skipped:
      user, oldserver = entry
      report.append( "%s, %s" % ( user, oldserver ) )
   report.append( "\n" )

   # Write reassignment report CSV file
   report = "\n".join( report )
   csvpath = "/tmp/reassign_report.csv" 
   with open( csvpath, 'w' ) as f:
      f.write( report )
      f.flush()
      f.close()

   if assignOnly:
      sendEmail( DEFAULT_SENDER, REPORT_HDR, INFRA_MAIL_REPORT_ONLY_MSG, 
                 attach=csvpath )
   else:
      sendEmail( DEFAULT_SENDER, REPORT_HDR, INFRA_MAIL_MSG, attach=csvpath )
   os.remove( csvpath )


def reassign( old, new, maxusers ):
   '''
   Given list of servers to move users away from, list of servers to migrate
   users to, and maximum number of users to be put on each new server, figure out
   reassignment for the users, or deletion of the user accounts.

   Parameters
   ----------
   old: list of strings
      List of servers to move users off of

   new: list of strings
      List of new servers to move users to

   maxusers: int
      Maximum number of users on each server
   '''

   # pdsh command timeout
   timeout = 30

   reassign = {}
   delete = {}
   unknown = {}
   for server in old:
      # List of usernames by server for accounts sitting on /home
      timeout = 30
      cmd = "pdsh -t %d -N -w %s ls /home" % ( timeout, server ) 
      accts = subprocess.check_output( shlex.split( cmd ) ).strip().split( '\n' )
      accts.remove( 'arastra' )

      acct_emails = []
      unknown[ server ] = []
      for acct in accts:
         cmd = "a4 users %s" % acct
         info = subprocess.check_output( shlex.split( cmd ) )

         if info:
            email = re.split( '<|>', info )[ 1 ]          
            acct_emails.append( email )
         else:
            unknown[ server ].append( acct )

      # List of user emails by server for accounts that have logged in the past
      # 30 days
      days = 30
      cmd = "arventory email --days %s %s" % ( days, server )
      out = subprocess.check_output( shlex.split( cmd ) ).strip().split( '\n' )

      # Remove first item because it contains "Users using ..." information line
      out.pop( 0 )
      
      # Remove commas from end of line
      emails = [ user.strip( ',' ) for user in out ]

      reassign[ server ] = []
      delete[ server ] = []
      reassign[ server ] = set( acct_emails ) & set( emails )
      delete[ server ] = set( acct_emails ) - set( emails )

      msg = """
Status of %s:
   # of accounts in /home: %d
   # of accounts with emails: %d
   # of accounts logged in past 30 days: %d
      - # of accounts reassigned: %d
      - # of accounts deleted: %d
      - # of accounts with no user info: %d
""" 
      print msg % ( server,
                    len( accts ),
                    len( acct_emails ),
                    len( emails ),
                    len( reassign[ server ] ),
                    len( delete[ server ] ),
                    len( unknown[ server ] ) )


   # Check there are enough servers for users reassigned
   maxnum = maxusers * len( new )
   totalusers = sum( len( l ) for _,l in reassign.iteritems() )
   if totalusers > maxnum:
      sys.exit( "There are more users to reassign than there is room" )

   # Reassign
   new_assign = {}
   for server in new:
      new_assign[ server ] = []
   new_assign = collections.OrderedDict( sorted( new_assign.items(), 
                                                 key=lambda t:t[ 0 ] ) )
   def _nextAvailable():
      '''
      Find next available new server
      '''
      for name, lst in new_assign.iteritems():
         if len( lst ) < maxusers:
            return name

   # Pack information into useful units
   for oldserver, emails in reassign.iteritems():
      for email in emails:
         user = email.split( '@' )[ 0 ]
         new_assign[ _nextAvailable() ].append( ( user, oldserver, email ) )

   del_assign = []
   for oldserver, emails in delete.iteritems():
      for email in emails:
         user = email.split( '@' )[ 0 ]
         del_assign.append( ( user, oldserver, email ) )

   skipped = []
   for oldserver, logins in unknown.iteritems():
      for login in logins:
         skipped.append( ( login, oldserver ) )

   return new_assign, del_assign, skipped


# ----------------------------------------------------------------------------------
# Main
# ----------------------------------------------------------------------------------

def main():
   parser = argparse.ArgumentParser( description=( "Reassign users from a list of "
                                     "servers to new list of servers" ) )
   parser.add_argument( 'from_servers', action="store", 
       help=( "Comma separated list of server name(s) you wish to migrate people " 
              "away FROM ( i.e: 'us001,us002,us003' )" ) )
   parser.add_argument( 'to_servers', action="store",
       help=( "Comma separated list of server name(s) you wish to migrate people "
              "away TO ( i.e.: 'us999,us998,us998' )" ) )
   parser.add_argument( 'max_users', action="store", type=int,
       help="Max number of users expected one each server" )
   parser.add_argument( '--assignOnly', action="store_true",
       help=( "If flag is set, the tool will only generate the CSV containing server"
              " assignments to email SW-INFRA_SUPPORT and skip emailing the "
              "affected users." ) )
   args = parser.parse_args()

   old = args.from_servers.split( ',' )
   new = args.to_servers.split( ',' )
   maxnum = args.max_users
   assignOnly = args.assignOnly

   # Figure out reassignments
   new_assign, del_assign, skipped = reassign( old, new, maxnum )

   # Notify users
   notify( new_assign, del_assign, skipped, assignOnly )


if __name__ == "__main__":
   main()

