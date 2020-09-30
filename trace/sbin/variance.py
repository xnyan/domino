#!/usr/bin/python
import sys
import json
import argparse

#delay = float(sys.argv[1])
#std_p = float(sys.argv[2])
#
#std = delay * std_p / 100.0
#
#p95 = delay + std * 2
#print delay
#print std
#print p95

arg_parser = argparse.ArgumentParser(description="Generate delay-conf.json")

arg_parser.add_argument('-c', '--config', dest='config', nargs='?', 
    help='delay configuration file', required=True)
arg_parser.add_argument('-j', '--jitter', dest='jitter', nargs='?',
    help='jitter %', required=True)
args = arg_parser.parse_args()

#f = sys.argv[1]
#std_p = float(sys.argv[2])
f = args.config
std_p = float(args.jitter)

config_file = open(f, "r")
config = json.load(config_file)
config_file.close()

dc_delay_map = config["oneway-delay"]
variance_map = {}#config["oneway-delay-variance"]
p95th_map = {}#config["dc-delay"]

k_list = dc_delay_map.keys()
for k in k_list:
    v_map = dc_delay_map[k]
    v_list = v_map.keys()
    variance_map[k] = {}
    p95th_map[k] = {}
    for v in v_list:
        s = dc_delay_map[k][v]
        lat = float(s[:len(s) - 2])
        std = lat * std_p / 100.0
        variance_map[k][v] = str(std) + "ms"
        p95th_map[k][v] = str(lat + std * 2) + "ms"

config["oneway-delay-variance"] = variance_map
config["dc-delay"] = p95th_map

out = open("delay-conf.json", "w")
out.write(json.dumps(config, indent=2, sort_keys=True) + "\n")
out.close()

