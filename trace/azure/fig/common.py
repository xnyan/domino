#!/usr/bin/python

import math
import os
import operator
import numpy as np
from enum import Enum

class Stat(Enum):
  __order__ = 'COUNT MEAN STDEV ERR95 P95th P99th MEDIAN MAX MIN'
  COUNT = 1
  MEAN = 2
  STDEV = 3
  ERR95 = 4 #95% confidence interval
  P95th = 5
  P99th = 6
  MEDIAN = 7
  MAX = 8
  MIN = 9

regex = "" #output data sepeartor

ns_table = {
    'ns'  : 1.0, \
    'us'  : 1000.0, \
    'ms'  : 1000000.0, \
    's'   : 1000000000.0, \
    'sec' : 1000000000.0, \
    'm'   : 60000000000.0, \
    'min' : 60000000000.0, \
    'h'   : 3600000000000.0, \
    'hour': 3600000000000.0, \
  }
def get_dv(unit):
  return ns_table[unit]

#calculate statistics
def cal_stats(data_list, ndigits=2):
  data_set = np.array(data_list, dtype=np.float64)
  stat_table = {}
  stat_table[Stat.COUNT] = len(data_list)
  stat_table[Stat.MEAN] = round(np.mean(data_set), ndigits)
  stdev = np.std(data_set)
  stat_table[Stat.STDEV] = round(stdev, ndigits)
  #95% percentile confidence inverval
  stat_table[Stat.ERR95] = round(1.96 * (stdev / math.sqrt(len(data_set))), ndigits)
  stat_table[Stat.P95th] = round(np.percentile(data_set, 95), ndigits)
  stat_table[Stat.P99th] = round(np.percentile(data_set, 99), ndigits)
  stat_table[Stat.MEDIAN] = round(np.median(data_set), ndigits)
  stat_table[Stat.MAX] = round(np.amax(data_set), ndigits)
  stat_table[Stat.MIN] = round(np.amin(data_set), ndigits)
  return stat_table
#end of statistics

# Percentile data from the given data list
def percentile(data, percentile):
  size = len(data)
  return sorted(data)[int(math.ceil((size * percentile) / 100.0)) - 1]

# File format: num1 num2 num3
def load_data(f):
  ret_list = list()
  trace_file = open(f, 'r')
  all_lines = trace_file.readlines()
  for line in all_lines:
    if line.startswith("#"):
      continue
    data = line.rstrip('\n').split()
    ret_list.append((long(data[0]), long(data[1]), long(data[2])))
  trace_file.close()
  return ret_list

def hist(d_list):
  h = {}
  for d in d_list:
    i = int(d)
    if i in h:
      h[i] += 1
    else:
      h[i] = 1
  r = sorted(h.items(), key = operator.itemgetter(0)) 
  return r

def gen_clock_rt_offset(data_l, clock_idx, is_clock_order, clock_unit, rt_unit, offset_unit):
  clock_l, rt_l, offset_l = list(), list(), list()
  clock_dv, rt_dv, offset_dv = get_dv(clock_unit), get_dv(rt_unit), get_dv(offset_unit)
  if is_clock_order:
    data_l.sort(key = operator.itemgetter(clock_idx))
  base_c = data_l[0][clock_idx]
  for i, d in enumerate(data_l):
    clock_l.append(float(d[clock_idx] - base_c) / clock_dv) # change unit
    rt_l.append(float(d[1] - d[0]) / rt_dv)
    offset_l.append(float(d[2] - d[0]) / offset_dv)
  return clock_l, rt_l, offset_l

# File format: Send_Clock Rev_Clock Server_Clock
def gen_predict_data(f, clock_unit, rt_unit, offset_unit):
  data_l = load_data(f)
  send_clock_l, send_rt_l, _ = gen_clock_rt_offset(data_l, 0, True, clock_unit, rt_unit, offset_unit)
  rev_clock_l, rev_rt_l, _ = gen_clock_rt_offset(data_l, 1, True, clock_unit, rt_unit, offset_unit)
  return send_clock_l, send_rt_l, rev_clock_l, rev_rt_l
  
# File format: Send_Clock Rev_Clock Server_Clock
def gen_data_by_rev_clock(f, is_clock_order, clock_unit, rt_unit, offset_unit):
  data_l = load_data(f)
  return gen_clock_rt_offset(data_l, 1, is_clock_order, clock_unit, rt_unit, offset_unit)

# File format: Send_Clock Rev_Clock Server_Clock
def gen_data_by_send_clock(f, is_clock_order, clock_unit, rt_unit, offset_unit):
  data_l = load_data(f)
  return gen_clock_rt_offset(data_l, 0, is_clock_order, clock_unit, rt_unit, offset_unit)

# File format: Send_Clock Rev_Clock Server_Clock
# Only uses sending time to filter the range
# clock offset and length must be in clock_unit
def gen_data_range(f, clock_offset, clock_length, clock_unit, rt_unit, offset_unit):
  clock_l, rt_l, _ = gen_data_by_send_clock(f, True, clock_unit, rt_unit, offset_unit)
  lower, upper = clock_l[0]+ clock_offset, clock_l[0] + clock_offset + clock_length
  ret_clock_l, ret_rt_l = list(), list()
  for i, cl in enumerate(clock_l):
    if cl >= lower and cl <= upper: 
      ret_clock_l.append(cl)
      ret_rt_l.append(rt_l[i])
    elif cl > upper:
      break
  return ret_clock_l, ret_rt_l
      

   






