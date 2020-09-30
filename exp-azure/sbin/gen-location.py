#!/usr/bin/python

import argparse
import os
import json

def LoadDcDelayMap(config_json_file):
  config_file = open(config_json_file, "r")
  config = json.load(config_file)
  config_file.close()
  dc_delay_map = config["oneway-delay"]
  dc_list = dc_delay_map.keys() # DC Id list
  dc_list.sort()
  for dc in dc_list:
    for i in dc_delay_map[dc].keys():
      s = dc_delay_map[dc][i]
      dc_delay_map[dc][i] = float(s[0:len(s)-2])
    dc_delay_map[dc][dc] = 0.0
  return dc_delay_map

def GetClosestReplica(c_dc, replica_dc_l, dc_delay_map):
  r = replica_dc_l[0]
  for r_dc in replica_dc_l:
    if dc_delay_map[c_dc][r_dc] < dc_delay_map[c_dc][r]:
      r = r_dc
  return r 

def maxLatInClosestQuorum(dc_delay_map, src_dc, dst_dc_set, quorum):
  lat = list()
  for i in dst_dc_set:
    lat.append(dc_delay_map[src_dc][i])
  lat.sort()
  return lat[quorum-1]

def GetLowestLatReplica(c_dc, replica_dc_l, dc_delay_map, quorum):
  r = replica_dc_l[0]
  lat = dc_delay_map[c_dc][r] + maxLatInClosestQuorum(dc_delay_map, r, replica_dc_l, quorum)
  for r_dc in replica_dc_l:
    l = dc_delay_map[c_dc][r_dc] + maxLatInClosestQuorum(dc_delay_map, r_dc, replica_dc_l, quorum)
    if l < lat:
      r = r_dc
      lat = l
  return r

# Parsing arguments
arg_parser = argparse.ArgumentParser(description="Generates location files based on the Azure VM IP files")

arg_parser.add_argument('-f', '--ipfile', dest='ip_file', nargs='?', 
    help='Azure VM public and private IP list file', required=True)
arg_parser.add_argument('-p', '--port', dest='port', nargs='?', default=10011, type=int,
    help='base port for servers')
arg_parser.add_argument('-l', '--leader', dest='leader', nargs='?', default='',
    help='DC id for the Multi-Paxos leader or Fast Paxos coordinator', required=True)
arg_parser.add_argument('-d', '--delay', dest='delay_file', nargs='?', default='',
    help='Delay configuration json file', required=True)
arg_parser.add_argument('-e', '--lowest', action='store_true',
    help='Sets a client target as the replica that achieves the lowest commit latency')
  
args = arg_parser.parse_args()
port = args.port
fp_leader_dc = args.leader
delay_config = args.delay_file
is_lowest = args.lowest

ip_file = open(args.ip_file, 'r')
all_lines = ip_file.readlines()
ip_file.close()

server_file = open("server-location.config", 'w+')
replica_file = open("replica-location.config", 'w+')
client_file = open("client-location.config", 'w+')

client_dir, server_dir = {}, {}
start = False
for line in all_lines:
  if line.startswith("#"):
    continue
  if line.startswith("--------"):
    start = True
    continue
  if start:
    info = line.rstrip('\n').split()
    vm_name, ip_pub, ip_pri = info[0], info[1], info[2]
    vm_dc = vm_name.split('-')
    dc = vm_dc[1]
    if vm_dc[0] == 'vm1': #client
      client_dir[dc] = (ip_pub, ip_pri)
    elif vm_dc[0] == 'vm2': #server
      server_dir[dc] = (ip_pub, ip_pri)
    else:
      print "Invalid vm tag " + vm_dc[0]

#Generate the leader in both server and replica location files
count = 1
for s_dc in server_dir:
  if s_dc == fp_leader_dc:
    pub_ip, pri_ip = server_dir[s_dc][0], server_dir[s_dc][1]
    replica_file.write(str(count) + ' ' + s_dc + ' ' + pri_ip + ' ' + str(port) + ' L L\n')
    server_file.write(str(count) + ' ' + s_dc + ' ' + pub_ip + '\n')
    break
#Generate server and replica location files
for s_dc in server_dir:
  print s_dc
  print server_dir[s_dc]
  if s_dc == fp_leader_dc:
    continue
  pub_ip, pri_ip = server_dir[s_dc][0], server_dir[s_dc][1]
  port += 1
  count += 1
  replica_file.write(str(count) + ' ' + s_dc + ' ' + pri_ip + ' ' + str(port) + ' L F\n')
  server_file.write(str(count) + ' ' + s_dc + ' ' + pub_ip + '\n')
replica_file.close()
server_file.close()

# Generate client location file
dc_delay_map = LoadDcDelayMap(delay_config)
s_dc_l = server_dir.keys()
m = (len(s_dc_l) - 1) / 2 + 1# majority
count = 1
for c_dc in client_dir:
  print c_dc
  print client_dir[c_dc]
  pub_ip = client_dir[c_dc][0] 
  if is_lowest:
    target_dc = GetLowestLatReplica(c_dc, s_dc_l, dc_delay_map, m)
  else:
    target_dc = GetClosestReplica(c_dc, s_dc_l, dc_delay_map)
  client_file.write(str(count) + ' ' + c_dc + ' ' + pub_ip + ' 1 ' + target_dc + '\n')
  count += 1
client_file.close()
