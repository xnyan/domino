#!/usr/bin/python

import common as c
import numpy as np

def load_predict_rate(f):
  ret = list()
  input_file = open(f, 'r')
  lines = input_file.readlines()
  input_file.close()
  for l in lines:
    if l.startswith("#"):
      continue
    data = l.rstrip('\n').split()
    ret.append((data[0], float(data[1]), float(data[2])))
  return ret

input_dir = 'roundtrip-predictrate-globe'
output_dir = input_dir
src_dc_l = [
        'eastus2',
        'westus2',
        'francecentral',
        'eastasia',
        'southeastasia',
        'australiaeast',
        ]
dst_dc_l = src_dc_l
suffix='-percentile-w1000.txt'

outf = open(output_dir+'/roundtrip-predictrate-table.csv', 'w+')
row = '\t'
for dst_dc in dst_dc_l:
  row = row + '\t,\t' + dst_dc
outf.write(row + '\n')

for src_dc in src_dc_l:
  row = src_dc
  for dst_dc in dst_dc_l:
    if src_dc == dst_dc:
      row = row + '\t,\t' + '-'
      continue

    f = input_dir+'/'+src_dc+'-'+dst_dc+suffix
    ret = load_predict_rate(f)
    rate = round(ret[0][2] * 100, 2)
    row = row + '\t,\t' + str(rate)
    print f
    print rate
  outf.write(row + '\n')



