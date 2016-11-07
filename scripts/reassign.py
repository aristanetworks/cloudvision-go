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
# Default Email Sender
DEFAULT_SENDER = "sw-infra-support@arista.com"


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




def notify( reassigns, dels, interns, skipped, assignOnly ):
   '''
   Send emails notifying users to reassign to a different server, or delete their
   old unused accounts.

   Parameters
   ----------
   reassigns: dictionary of the format -

      { newserver : [ ( user, oldserver, email ), ... ] }

      Where newserver is name of new server to migrate to, user is username,
      oldserver is name of old server to migrate from, and email is the user's listed
      email on a4 users. Represents users to be migrated.

   dels: list of the format -

      [ ( user, oldserver, email ), ... ]

      Items are similar to the description for new_assign. Represents users to be
      deleted.

   interns: list of the format - 
      [ ( user, oldserver, { name, manager, mentor } ), ... ]

      List of intern user names with a dictionary containing their full name,
      manager username, and mentor username.

   skipped: list of the format -

      [ ( user, oldserver ), ... ]

      Represents skipped users because their username does not exist in "a4 users"
      output.

   assignOnly: boolean
      If True, send assignment report only to SW-INFRA-SUPPORT.
      If False, send assignment report to SW-INFRA-SUPPORT and also email the 
         affected users.
   '''
   DEFAULT_STAT = "1st Email Sent"
   REPONLY_STAT = "Report Only - no emails sent"

   report = []
   report.append( "User Migration Report\n\n" )

   # Email Reassignments
   reassign_hdr = "Time for a New Home! (new server assignment)"
   reassign_sv = "Your current server: %s\nYour new server: %s"
   reassign_body = ""    
   with open( "./templates/account_reassign", 'r' ) as f:
      reassign_body = f.read()

   for newserver, users in reassigns.iteritems():
      report.append( "Accounts Notified to be Reassigned to %s" % newserver )
      report.append( "User,Previous Server,New Server,Email,Status" )
      for entry in users:
         user, oldserver, email = entry
         body = "%s\n\n%s" % ( reassign_sv % ( oldserver, newserver ), 
                               reassign_body )
         stat = REPONLY_STAT
         if not assignOnly:
            sendEmail( email, reassign_hdr, body )
            stat = DEFAULT_STAT

         report.append( "%s,%s,%s,%s,%s" %
                        ( user, oldserver, newserver, email, stat ) )
      report.append( "\n" )
   report.append( "\n" )


   # Email deletions to affected users
   delete_hdr = "Please delete your infrequently used account"
   delete_sv = "Your current server: %s"
   delete_body = "" 
   with open( "./templates/account_delete", 'r' ) as f:
      delete_body = f.read()

   report.append( "Inactive Accounts (no login in 30 days) Notified to be Deleted" )
   report.append( "User,Previous Server,Email,Status" )
   for entry in dels:
      user, oldserver, email = entry
      body = "%s\n\n%s" % ( delete_sv % ( oldserver ), delete_body )
      stat = REPONLY_STAT       
      if not assignOnly:
         sendEmail( email, delete_hdr, body )
         stat = DEFAULT_STAT

      report.append( "%s,%s,%s,%s" % ( user, oldserver, email, stat ) )
   report.append( "\n" )


   # Email Intern Accounts
   intern_hdr = "Please clean up your intern's account"
   intern_sv = "Server with Intern Account: %s\n"
   intern_inf = "Intern Name: %s\nIntern Username: %s\nMentor: %s\nManager: %s"
   intern_body = "" 
   with open( "./templates/intern_delete", 'r' ) as f:
      intern_body = f.read()

   report.append( "Intern Accounts" )
   report.append( "User,Server,Full Name,Mentor,Manager,Status" )
   for entry in interns:
      user, oldserver, info = entry
      intName = info[ 'name' ]
      intMentor = info[ 'mentor' ].replace( " ", "" ) 
      intManag = info[ 'manager' ].replace( " ", "" )
      email = "%s@arista.com" % ( intMentor if intMentor else intManag )
      body = "%s%s\n\n%s" % ( intern_sv % oldserver,
                              intern_inf % ( intName, 
                                             user,
                                             intMentor if intMentor else "Unknown",
                                             intManag if intManag else "Unknown" ), 
                              intern_body )
      stat = REPONLY_STAT
      if not assignOnly:
         sendEmail( email, intern_hdr, body )
         stat = DEFAULT_STAT
      
      report.append( "%s,%s,%s,%s,%s,%s" % 
                     ( user, oldserver, intName, intMentor, intManag, stat ) )
   report.append( "\n" )

   # Write reassignment report CSV file
   report = "\n".join( report )
   csvpath = "/tmp/reassign_report.csv" 
   with open( csvpath, 'w' ) as f:
      f.write( report )
      f.flush()

   report_hdr = "Server Reassignment Report"
   if assignOnly:
      #report_body_reponly = "" 
      with open( "./templates/report_only", 'r' ) as f:
         report_body_reponly = f.read()
      sendEmail( DEFAULT_SENDER, report_hdr, report_body_reponly, attach=csvpath )
   else:
      #report_body = ""    
      with open( "./templates/report_and_send", 'r' ) as f:
         report_body = f.read()
      sendEmail( DEFAULT_SENDER, report_hdr, report_body, attach=csvpath )

   os.remove( csvpath )



def reassign( old, new, all_interns, maxusers ):
   '''
   Given list of servers to move users away from, list of servers to migrate
   users to, and maximum number of users to be put on each new server, figure out
   reassignment for the users, or deletion of the user accounts.

   Also takes into account intern accounts given a dictionary of most up-to-date
   intern information (ask itthichok@ for most recent info).

   Parameters
   ----------
   old: list of strings
      List of servers to move users off of

   new: list of strings
      List of new servers to move users to

   all_interns: dict     
      Contains intern account, mentor, manager info
      Depends on information from itthichok@

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


   # Figure out if unknown accounts are intern accounts
   interns = []
   skipped = []
   for oldserver, logins in unknown.iteritems():
      for login in logins:
         if login in all_interns:
            interns.append( ( login, oldserver, all_interns[ login ] ) )
         else:
            skipped.append( ( login, oldserver ) )

   return new_assign, del_assign, interns, skipped



# ----------------------------------------------------------------------------------
# Main
# ----------------------------------------------------------------------------------

def main():

   parser = argparse.ArgumentParser( description=( "Reassign users from a list of "
      "servers to new list of servers. ***NOTE*** Expects email templates to be "
      "under templates/ in current working directory." ) )
   parser.add_argument( 'from_servers', action="store", 
      help=( "Comma separated list of server name(s) you wish to migrate people " 
             "away FROM ( i.e: 'us001,us002,us003' )" ) )
   parser.add_argument( 'to_servers', action="store",
      help=( "Comma separated list of server name(s) you wish to migrate people "
             "away TO ( i.e.: 'us999,us998,us998' )" ) )
   parser.add_argument( 'max_users', action="store", type=int,
      help="Max number of users expected one each server" )
   parser.add_argument( 'interns', action="store",
      help=( "Path containing CSV containing intern information gathered by "
             "itthichok@ that will be used to cross check accounts with no "
             "user information" ) )
   parser.add_argument( '--assignOnly', action="store_true",
      help=( "If flag is set, the tool will only generate the CSV containing server"
              " assignments to email SW-INFRA_SUPPORT and skip emailing the "
              "affected users." ) )
   args = parser.parse_args()

   old = args.from_servers.split( ',' )
   new = args.to_servers.split( ',' )
   maxnum = args.max_users
   interns_filepath = args.interns
   assignOnly = args.assignOnly
   
   print "\nReassigning users from [%s] to [%s]..." % ( args.from_servers, 
                                                    args.to_servers )

   # Ask user if they checked email templates to be what they want
   ans = ""
   while ans != "yes" and ans != "y":
      ans = raw_input( ( 'Did you double check the email templates under /templates?'
                         ' They will be sent as is. [Yes/Y/No/N]: ' ) ).lower()
      if ans == "no" or ans == "n":
         print >> sys.stderr, "Aborting - please double check the email templates!"
         sys.exit( 1 )
      elif ans == "yes" or ans == "y":
         print "Great, continuing..."
         break
      else:
         print "Please answer [Yes/Y/No/N]. I believe in you."


   def _interns( fpath ):
      # Retrieve intern information from CSV file provided by itthichok@
      all_interns={}
      with open( interns_filepath, 'r' ) as f:
         for line in f:
            if "firstName" in line:
               # Skip the first line containing format of this file
               # XXX: We're going to assume format of the file is in expected
               # firstName,lastName,userName,mentor,manager
               continue
            else:
               l = line.strip( '\n' ).split( ',' )
               all_interns[ l[ 2 ] ] = { "name": "%s %s" % ( l[0], l[1] ),
                                         "mentor": l[3],
                                         "manager": l[4] }
      return all_interns

   # Figure out reassignments
   reassigns, dels, interns, skipped = reassign( old, 
                                                 new, 
                                                 _interns( interns_filepath ), 
                                                 maxnum )

   # Notify users
   notify( reassigns, dels, interns, skipped, assignOnly )


if __name__ == "__main__":
   main()

