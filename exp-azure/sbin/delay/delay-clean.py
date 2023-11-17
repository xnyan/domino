#!/usr/bin/python
import json
import argparse
import os

default_user = os.environ['USER']

arg_parser = argparse.ArgumentParser(description="Set delays among datacenters.")

# Cluster configuration file
arg_parser.add_argument('-c', '--config', dest='config', nargs='?', 
    help='cluster configuration file', required=True)
arg_parser.add_argument('-u', '--user', dest='user', nargs='?',
    help='username', default=default_user)
arg_parser.add_argument('-d', '--dev', dest='dev', nargs='?',
    help='network interface device name', default='eno1')
arg_parser.add_argument('-p', '--parallel', action='store_true')

args = arg_parser.parse_args()

user = args.user
dev = args.dev
suffix = ''
if args.parallel:
  suffix = " &"

#Reads configurations
config_file = open(args.config, "r")
config = json.load(config_file)
config_file.close()

dc_ip_map = config["datacenter"]
dc_delay_map = config["oneway-delay"]

#cmd prefix
tc = 'sudo tc'

#cmd
clean_cmd = '%s qdisc del dev %s root;' % (tc, dev)

#Sets up delays among different DCs
dc_ip_list = dc_ip_map.keys()
dc_ip_list.sort()
for dc_id in dc_ip_list:
  print "Datacenter: %s" % (dc_id)
  ip_list = dc_ip_map[dc_id]
  dst_delay_table = dc_delay_map[dc_id]
 
  if not dst_delay_table:
    continue
  
  for ip in ip_list:
    shell_cmd = "ssh -n %s@%s \"%s" % (user, ip, clean_cmd) 
    shell_cmd += "\"" + suffix
    print "Executes: %s" % (shell_cmd)
    os.system(shell_cmd)
    print"Done"
    print""
  print""

