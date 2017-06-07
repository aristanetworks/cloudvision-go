#!/bin/python

import json
import requests
import subprocess
import pprint
import argparse
import sys
from pprint import pprint
from random import randint

try:
    from tinydb import TinyDB, Query
except:
    print """Some dependencies are needed by this tool, install with:
    sudo yum install python-pip || sudo apt-get install python-pip
    sudo pip install tinydb
    """
    sys.exit(-1)

def get_range_for_ip_ks(keyspace, node_ip, nodetool_host='127.0.0.1', nodetool_port='7199', api_host='127.0.0.1', primary_range=True, nonprimary_range=True):
    cmds = []
    url = 'http://{}:10000/storage_service/describe_ring/{}'.format(api_host, keyspace)
    resp = requests.get(url=url)
    data = json.loads(resp.text)
    if 'code' in data and 'message' in data:
        print "ERROR: ", data['message']
        return ""
    print "====== Listing repair cmds for keyspace {} on node {} ======".format(keyspace, node_ip)
    for d in data:
        primary_ip = d['endpoints'][0]
        nonprimary_ips = d['endpoints'][1:]
        if primary_range:
            if node_ip == primary_ip:
                st = d['start_token']
                et = d['end_token']
                cmd = 'nodetool -h {} -p {} repair --start-token {} --end-token {} {}'.format(nodetool_host, nodetool_port, st, et, keyspace)
                print "    PRIMARY_RANGE: ", cmd
                cmds.append(cmd)
        if nonprimary_range:
            if node_ip in nonprimary_ips:
                st = d['start_token']
                et = d['end_token']
                cmd = 'nodetool -h {} -p {} repair --start-token {} --end-token {} {}'.format(nodetool_host, nodetool_port, st, et, keyspace)
                print "NON_PRIMARY_RANGE: ", cmd
                cmds.append(cmd)
    return cmds

def run_cmd(cmd):
    # emulate_err is for testing only
    emulate_err = False
    status = "FAIL"
    ret = subprocess.call(cmd, shell=True)
    if ret == 0:
        status = "SUCCEED"
    else:
        status = "FAIL"
    if emulate_err and randint(0,3) == 1:
        # emluate some failure ;-)
        status = "FAIL"
    return status

def show_succeed_fail(tag, cmds_ok, cmds_fail):
    print '{}: SUCCEED={}: '.format(tag, len(cmds_ok))
    print '{}: FAIL={}: '.format(tag, len(cmds_fail))
    if cmds_fail:
        pprint(cmds_fail)

def show_succeed_fail_nr(tag, cmds_ok, cmds_fail):
    print '{}: SUCCEED={}, FAIL={}'.format(tag, len(cmds_ok), len(cmds_fail))

def get_db_file_name(keyspace):
    return 'scylla_repair_results_{}.json'.format(keyspace)

def run_failed_repair(keyspace, dryrun):
    print "############  SCYLLA REPAIR: MODE=CONTINUE START #############"
    db = TinyDB(get_db_file_name(keyspace))
    query = Query()
    cmds_fail = db.search(query.status == 'FAIL')
    cmds_ok = db.search(query.status == 'SUCCEED')
    show_succeed_fail_nr("REPAIR SUMMARY BEFORE", cmds_ok, cmds_fail)
    cmds = cmds_fail
    if dryrun:
        print "====== Listing previously failed repair cmds for keyspace {} ======".format(keyspace)
        for cmd in cmds:
            print cmd['cmd']
    else:
        if cmds:
            print "====== Running previously failed repair cmds for keyspace {} ======".format(keyspace)
        for i in xrange(len(cmds)):
            db_row = cmds[i]
            cmd = db_row['cmd']
            cmd_id = db_row['cmd_id']
            node_ip = db_row['node_ip']
            print "Reparing for [{} / {}] failed sub ranges on {}".format(i+1, len(cmds), node_ip)
            status = run_cmd(cmd)
            if status == 'SUCCEED':
                db.update({'status' : 'SUCCEED'}, query.cmd_id == cmd_id)
        # show_succeed_fail(node_ip, cmds_ok, cmds_fail)
        cmds_fail = db.search(query.status == 'FAIL')
        cmds_ok = db.search(query.status == 'SUCCEED')
        show_succeed_fail_nr("REPAIR SUMMARY AFTER", cmds_ok, cmds_fail)
    print "############  SCYLLA REPAIR: MODE=CONTINUE END   #############"

def run_repair(keyspace, hosts_list, ports_list, api_host, primary_range, nonprimary_range, dryrun):
    db = TinyDB(get_db_file_name(keyspace))
    db.purge()
    cmd_id = 0
    print "############  SCYLLA REPAIR: MODE=NORMAL     START #############"
    print "HOSTS = ", hosts_list
    print "PORTS = ", ports_list
    for x in xrange(len(hosts_list)):
        node_ip = hosts_list[x]
        nodetool_port = ports_list[x]
        nodetool_host = '127.0.0.1'
        cmds = get_range_for_ip_ks(keyspace, node_ip, nodetool_host, nodetool_port, api_host, primary_range, nonprimary_range)
        cmds_ok = []
        cmds_fail = []
        if dryrun:
            continue
        print "====== Running repair cmds for keyspace {} on node {} ======".format(keyspace, node_ip)
        for i in xrange(len(cmds)):
            cmd = cmds[i]
            print "Repairing for [{} / {}] sub ranges on {}".format(i+1, len(cmds), node_ip)
            status = run_cmd(cmd)
            cmd_id = cmd_id + 1
            if status == 'SUCCEED':
                cmds_ok.append(cmd)
            elif status == 'FAIL':
                cmds_fail.append(cmd)
            else:
                cmds_dryrun.append(cmd)
            db_row = {'node_ip' : node_ip, 'status' : status,  "cmd" : cmd, "idx" : i, 'cmd_id' : cmd_id}
            db.insert(db_row)
        show_succeed_fail(node_ip, cmds_ok, cmds_fail)
    if not dryrun:
        query = Query()
        cmds_fail = db.search(query.status == 'FAIL')
        cmds_ok = db.search(query.status == 'SUCCEED')
        show_succeed_fail_nr("REPAIR SUMMARY", cmds_ok, cmds_fail)
    print "############  SCYLLA REPAIR: MODE=NORMAL     END   #############"

if __name__ == '__main__':
    eg = '''
    Examples: 
    1) Run repair on node 127.0.0.{1,2,3} for both primary range and non-primary range
    ./scylla_repair.py --keyspace ks3 --hosts 127.0.0.1,127.0.0.2,127.0.0.3 --ports 7199,7200,7300 --pr --npr
    2) Run with --cont option to rerun repair for the failed range
    ./scylla_repair.py --keyspace ks3 --hosts 127.0.0.1,127.0.0.2,127.0.0.3 --ports 7199,7200,7300 --cont
    3) Print the repair cmd for primary range for node 127.0.0.3
    ./scylla_repair.py --keyspace ks3 --hosts 127.0.0.3 --ports 7300 --dryrun --pr
    '''
    parser = argparse.ArgumentParser(description='Scylla Repair' + eg , formatter_class=argparse.RawTextHelpFormatter)
    parser.add_argument('--cont', help='Rerun repair for the failed sub subranges', action='store_true')
    parser.add_argument('--keyspace', help='Keyspace to repair', required=True)
    parser.add_argument('--apihost', help='HTTP API HOST', required=False)
    parser.add_argument('--hosts', help='Hosts to repair', required=True)
    parser.add_argument('--ports', help='Ports for nodetool')
    parser.add_argument('--pr', help='Repair primary ranges', action='store_true')
    parser.add_argument('--npr', help='Repair non primary ranges', action='store_true')
    parser.add_argument('--dryrun', help='Do not run nodetool repair, print repair cmd only', action='store_true')
    args = parser.parse_args()
    keyspace = args.keyspace
    hosts_list = args.hosts.split(',');
    ports_list = ['7199'] * len(hosts_list)
    primary_range = args.pr
    nonprimary_range = args.npr
    api_host = '127.0.0.1'
    if primary_range == False and nonprimary_range == False:
        print "Specify --pr (PrimaryRange) and/or --npr (NonPrimaryRange) to repair"
        sys.exit(-1)
        
    if args.apihost:
        api_host = args.apihost
    if args.ports:
        ports_list = args.ports.split(',');
    if len(hosts_list) == len(ports_list):
        if args.cont:
            run_failed_repair(keyspace, args.dryrun)
        else:
            run_repair(keyspace, hosts_list, ports_list, api_host, primary_range, nonprimary_range, args.dryrun)
    else:
        print "hosts_list = ", hosts_list
        print "ports_list = ", ports_list
        print "Number of hosts and ports do not match"
