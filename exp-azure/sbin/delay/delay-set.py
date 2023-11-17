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
arg_parser.add_argument('-b', '--bandwidth', dest='bandwidth', nargs='?',
    help='bandwidth', default='1000Mbps')
arg_parser.add_argument('-d', '--dev', dest='dev', nargs='?',
    help='network interface device name', default='eno1')
arg_parser.add_argument('-p', '--parallel', action='store_true')

args = arg_parser.parse_args()

user = args.user
bandwidth = args.bandwidth
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
class_cmd = '%s class add dev %s parent' % (tc, dev)
delay_cmd = '%s qdisc add dev %s handle' % (tc, dev)
filter_cmd = '%s filter add dev %s pref' % (tc, dev)

#cmd
clean_cmd = '%s qdisc del dev %s root;' % (tc, dev)
setup_cmd = '%s qdisc add dev %s root handle 1: htb;' % (tc, dev)
setup_cmd += '%s 1: classid 1:1 htb rate %s;' % (class_cmd, bandwidth)

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
    shell_cmd = "ssh -n %s@%s \"%s %s" % (user, ip, clean_cmd, setup_cmd) 
    handle = 1
    dst_delay_list = dst_delay_table.keys()
    for dst_dc_id in dst_delay_list:
      delay = dst_delay_table[dst_dc_id]
      if dst_dc_id not in dc_ip_map:
          continue
      dst_ip_list = dc_ip_map[dst_dc_id]

      handle += 1
      shell_cmd += "%s 1:1 classid 1:%d htb rate %s;" % \
                  (class_cmd, handle, bandwidth)
      shell_cmd += "%s %d: parent 1:%d netem delay %s;" % \
                  (delay_cmd, handle, handle, delay)
      for dst_ip in dst_ip_list:
        shell_cmd += "%s %d protocol ip u32 match ip dst %s flowid 1:%d;" % \
                    (filter_cmd, handle, dst_ip, handle)
    shell_cmd += "\"" + suffix
    print "Executes: %s" % (shell_cmd)
    os.system(shell_cmd)
    print"Done"
    print""
  print""

