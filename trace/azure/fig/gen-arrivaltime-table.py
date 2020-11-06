#!/usr/bin/python

import common as c
import numpy as np

def load_owd_stat(f):
  ret = list()
  input_file = open(f, 'r')
  lines = input_file.readlines()
  input_file.close()
  c = 1
  for l in lines:
    if l.startswith("#"):
      continue
    data = l.rstrip('\n').split()
    if c == 1:
        # On time
        ret.append(data[0]) # Mean diff
        ret.append(data[2]) # Median diff
        ret.append(data[3]) # P95 diff
        ret.append(data[4]) # P99 diff
    if c == 2:
        # Expire
        ret.append(data[0]) # Mean diff
        ret.append(data[2]) # Median diff
        ret.append(data[3]) # P95 diff
        ret.append(data[4]) # P99 diff
    if c == 3:
        # Expire rate
        ret.append(data[2]) #%
    c = c+1
  return ret

def gen_table(f, idx):
    outf = open(f, 'w+')
    row = '-'
    for dst_dc in dst_dc_l:
      row = row + '\t' + dst_dc
    outf.write(row + '\n')
    for src_dc in src_dc_l:
      row = src_dc
      for dst_dc in dst_dc_l:
        if src_dc == dst_dc:
          row = row + '\t' + '-'
          continue
        f = input_dir+'/'+src_dc+'-'+dst_dc+suffix
        ret = load_owd_stat(f)
        rate = ret[idx]
        row = row + '\t' + str(rate)
        print f
        print rate
      outf.write(row + '\n')
    outf.close()

#input_dir = 'arrivaltime-globe-raw'
input_dir = 'arrivaltime-globe'
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
suffix='.log.txt-arrivaltime-pth95.0-window1000ms.txt'

gen_table(output_dir+'/arrivaltime-table-expire-rate.csv', 8)
gen_table(output_dir+'/arrivaltime-table-expire-mean.csv', 4)
gen_table(output_dir+'/arrivaltime-table-expire-median.csv', 5)
gen_table(output_dir+'/arrivaltime-table-expire-p95.csv', 6)
gen_table(output_dir+'/arrivaltime-table-expire-p99.csv', 7)


