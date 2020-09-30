#!/usr/bin/python

import os
import json

# df file is gnerated by gen-lat-stat.py
df='lat_stat_global_8dc.txt'
outf='azure-globe-delay.json'
dc_l = ['westus2', 'eastus2', 'australiaeast', 'southeastasia', 'eastasia', 'francecentral']

df='lat_stat_na_9dc.txt'
outf='azure-na-delay.json'
dc_l = ['westus', 'northcentralus', 'southcentralus', 'centralus', 'eastus2', 'canadacentral', 'canadaeast', 'westcentralus', 'westus2']

avg_map = {}
stat_f = open(df, 'r')
all_lines = stat_f.readlines()
for line in all_lines:
  if line.startswith("#"):
    continue
  data = line.rstrip('\n').split(',')
  src_dc = data[0][2:-1]
  dst_dc = data[1][2:-1]
  avg = data[2].strip()
  print src_dc + '->' + dst_dc + "=" + avg
  if src_dc not in avg_map.keys():
    avg_map[src_dc] = {}
  avg_map[src_dc][dst_dc] = avg
stat_f.close()

print avg_map

delay_map = {}
delay_map['datacenter'] = {}
delay_map['oneway-delay'] = {}
for src_dc in dc_l:
  delay_map['datacenter'][src_dc] = []
  delay_map['oneway-delay'][src_dc] = {}
  for dst_dc in dc_l:
    if src_dc is dst_dc:
      continue
    print src_dc + '->' + dst_dc
    owd = round(float(avg_map[src_dc][dst_dc])/2.0, 2)
    delay_map['oneway-delay'][src_dc][dst_dc] = str(owd)+'ms'

out = open(outf, 'w')
out.write(json.dumps(delay_map, indent=2, sort_keys=True) + "\n")
out.close()
