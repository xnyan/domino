#!/usr/bin/python

import argparse
import os

arg_parser = argparse.ArgumentParser(description="Splits Azure VM public and private IPs into .pub and .pri files")

arg_parser.add_argument('-f', '--ipfile', dest='ip_file', nargs='?', 
    help='Azure VM public and private IP list file', required=True)
arg_parser.add_argument('-p', '--port', dest='port', nargs='?', default='',
    help='generates the given port after each private IP')
  
args = arg_parser.parse_args()

port = ''
if args.port != '':
  port = ':' + args.port


ip_file = open(args.ip_file, 'r')
all_lines = ip_file.readlines()
ip_file.close()

pub_file = open(args.ip_file+".public", 'w+')
pri_file = open(args.ip_file+".private", 'w+')

ip_start = False
for line in all_lines:
  if line.startswith("#"):
    continue
  if line.startswith("--------"):
    ip_start = True
    continue
  if ip_start:
    ip_list = line.rstrip('\n').split()
    vm_name, ip_pub, ip_pri = ip_list[0], ip_list[1], ip_list[2]
    pub_file.write(vm_name + "\t" + ip_pub + "\n")
    ip_pri = ip_pri + port 
    pri_file.write(vm_name + "\t" + ip_pri + "\n")

pub_file.close()
pri_file.close()
