#!/usr/bin/python
import numpy as np
import math
import os
#from aenum import Enum
from enum import Enum

#Constants
warm_up_time = 15.0 #s
measure_time = 60.0 #s
cool_down_time = 15.0 #s
digits = 2 #decimal digits
regex = "," #output data sepeartor

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

class Metric(Enum):
  __order__ = 'LATENCY THROUGHPUT ACCEPT_RATE'
  LATENCY = 1
  THROUGHPUT = 2
  ACCEPT_RATE = 3
  ACCEPT_FAST_RATE = 4

Stat_List = [Stat.MEAN, Stat.MEDIAN, Stat.P95th, Stat.P99th, Stat.MAX]

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

#Filter data record from a file
def filter_data(data_file, warm_up_time, measure_time):
  data_record_list = list()
  w_time = warm_up_time * 1000 * 1000 * 1000 #ns
  m_time = measure_time * 1000 * 1000 * 1000 #ns
  line_count = 0
  exp_start_time = 0
  input_file = open(data_file, 'r')
  all_lines = input_file.readlines()

  for line in all_lines:
    if line.startswith("#"):
      continue
    line_count = line_count + 1
    stat_array = line.rstrip('\n').split(regex)
    if line_count == 1:
      exp_start_time = long(stat_array[4].strip())
      measure_start_time = exp_start_time + w_time
      measure_end_time = measure_start_time + m_time
    req_start_time = long(stat_array[4].strip())
    if req_start_time < measure_start_time:
      # Skips the requests in warm up time
      continue
    if req_start_time > measure_end_time:
      # Ignores the requests in the cool down time
      input_file.close()
      return data_record_list
    stat_array = line.rstrip('\n').split(',')
    is_accept = stat_array[1].strip()
    if is_accept == '1':
      data_record_list.append(line)
  print("Notice: No requests during the cool-down time, experiment:" + data_file)
  input_file.close()
  return data_record_list
  #exit()

#Calculates the average metrics given a list of data
def avg_metrics(data_record_list, measure_time):
  data_stat = {}
  #data_stat[Metric.ACCEPT_RATE] = 0.0
 
  accept_num = 0
  reject_num = 0
  accept_fast_num = 0
  lat_list = list()
  for line in data_record_list:
    stat_array = line.rstrip('\n').split(',')
    # Latency
    lat = float(stat_array[3].strip())
    lat_list.append(lat)
    is_accept = stat_array[1].strip()
    is_fast = stat_array[2].strip()
    if is_accept == '1':
      accept_num +=1 
      if is_fast == '1':
        accept_fast_num += 1
    else:
        reject_num += 1
  # Latency
  lat_stat = cal_stats(lat_list)
  data_stat[Metric.LATENCY] = lat_stat #lat_stat[Stat.MEAN]
  #data_stat[Metric.LATENCY] = lat_stat[Stat.MEAN]
  # Accept Throughput 
  data_stat[Metric.THROUGHPUT] = round(1.0 * accept_num / measure_time, digits)
  data_stat[Metric.ACCEPT_RATE] = round(1.0 * accept_num / (accept_num + reject_num), digits)
  data_stat[Metric.ACCEPT_FAST_RATE] = round(1.0 * accept_fast_num / accept_num, digits)
  
  return data_stat

def print_stat(stat):
  print(stat[Metric.LATENCY][Stat.MEAN], stat[Metric.THROUGHPUT], \
          stat[Metric.ACCEPT_RATE], stat[Metric.ACCEPT_FAST_RATE])

#Filters the raw latency for accepted operations in a dir
def get_exp_lat(exp_dir, num=5, warm_up_time=warm_up_time, measure_time=measure_time):
  lat_list = list()
  path = os.path.abspath(exp_dir)
  for i in range(1, num + 1):
    m_dir = os.path.join(path, str(i))
    #print(m_dir)
    for sub_root, sub_dirs, data_files in os.walk(m_dir):
      for data_file in data_files:
        if data_file.endswith('.log') and data_file.startswith('client-'):
          data_file = os.path.join(sub_root, data_file)
          #print(data_file)
          client_data_list = filter_data(data_file, warm_up_time, measure_time)
          #print("number of requests: " + str(len(client_data_list)))
          for line in client_data_list:
            stat_array = line.rstrip('\n').split(',')
            lat = float(stat_array[3].strip())
            lat_list.append(lat)
  return lat_list

#Get the latencies in one data file
def get_lat(data_file, warm_up_time=warm_up_time, measure_time=measure_time):
  lat_list = list()
  client_data_list = filter_data(data_file, warm_up_time, measure_time)
  for line in client_data_list:
    stat_array = line.rstrip('\n').split(',')
    lat = float(stat_array[3].strip())
    lat_list.append(lat)
  return lat_list

#Given the experiment dir, calculates the metrics of the given num measurements
def exp_metric(exp_dir, num=5, warm_up_time=warm_up_time, measure_time=measure_time):
  path = os.path.abspath(exp_dir)
  exp_stat = {}
  for metric in Metric:
    if metric is Metric.LATENCY:
      exp_stat[metric] = {}
      for stat in Stat_List:
        exp_stat[metric][stat] = list()
    else:
      exp_stat[metric] = list()
  for i in range(1, num + 1):
    m_dir = os.path.join(path, str(i))
    #print(m_dir)
    m_data_list = list()
    for sub_root, sub_dirs, data_files in os.walk(m_dir):
      for data_file in data_files:
        if data_file.endswith('.log') and data_file.startswith('client-'):
          data_file = os.path.join(sub_root, data_file)
          #print(data_file)
          client_data_list = filter_data(data_file, warm_up_time, measure_time)
          #print("number of requests: " + str(len(client_data_list)))
          m_data_list.extend(client_data_list)
    m_stat=avg_metrics(m_data_list, measure_time)
    #print_stat(m_stat)
    for metric in Metric:
      if metric is Metric.LATENCY:
        for stat in Stat_List:
          exp_stat[metric][stat].append(m_stat[metric][stat])
      else:
        exp_stat[metric].append(m_stat[metric])
  for metric in Metric:
    if metric is Metric.LATENCY:
      for stat in Stat_List:
        exp_stat[metric][stat] = cal_stats(exp_stat[metric][stat])
    else:
      exp_stat[metric] = cal_stats(exp_stat[metric])
  return exp_stat

#Get the metrics for a list of data files
def client_metric(data_file_list, num=5, warm_up_time=warm_up_time, measure_time=measure_time):
  exp_stat = {}
  for metric in Metric:
    if metric is Metric.LATENCY:
      exp_stat[metric] = {}
      for stat in Stat_List:
        exp_stat[metric][stat] = list()
    else:
      exp_stat[metric] = list()
  for data_file in data_file_list:
    print(data_file)
    client_data_list = filter_data(data_file, warm_up_time, measure_time)
    print("number of requests: " + str(len(client_data_list)))
    m_stat=avg_metrics(client_data_list, measure_time)
    print_stat(m_stat)
    for metric in Metric:
      if metric is Metric.LATENCY:
        for stat in Stat_List:
          exp_stat[metric][stat].append(m_stat[metric][stat])
      else:
        exp_stat[metric].append(m_stat[metric])
  for metric in Metric:
    if metric is Metric.LATENCY:
      for stat in Stat_List:
        exp_stat[metric][stat] = cal_stats(exp_stat[metric][stat])
    else:
      exp_stat[metric] = cal_stats(exp_stat[metric])
  return exp_stat

def dump_exp_metric(output_file, append, exp_dir, num=5, warm_up_time=warm_up_time, measure_time=measure_time):
  exp_stat = exp_metric(exp_dir, num, warm_up_time, measure_time)
  if append:
    f = open(output_file, "a")
  else:
    f = open(output_file, "w")
    f.write("#lat mean,err95,median, err95, 95th, err95, 99th, err95, thr mean, err95, accept rate, err95, fast_accept_rate, err95\n")

  for metric in Metric:
    if metric is Metric.LATENCY:
      for stat in Stat_List:
        f.write(str(exp_stat[metric][stat][Stat.MEAN]) + " ")
        f.write(str(exp_stat[metric][stat][Stat.ERR95]) + " ")
    else:
      f.write(str(exp_stat[metric][Stat.MEAN]) + " ")
      f.write(str(exp_stat[metric][Stat.ERR95]) + " ")
  f.write("\n")
  f.close()

def load_exp_metric(data_file):
  exp_stat = {}
  for metric in Metric:
    if metric is Metric.LATENCY:
      exp_stat[metric] = {}
      for stat in Stat_List:
        exp_stat[metric][stat] = {}
        exp_stat[metric][stat][Stat.MEAN] = list()
        exp_stat[metric][stat][Stat.ERR95] = list()
    else:
      exp_stat[metric] = {}
      exp_stat[metric][Stat.MEAN] = list()
      exp_stat[metric][Stat.ERR95] = list()
  f = open(data_file, "r")
  all_lines = f.readlines()
  for line in all_lines:
    if line.startswith("#"):
      continue
    stat_array = line.rstrip('\n').split()
    exp_stat[Metric.LATENCY][Stat.MEAN][Stat.MEAN].append(float(stat_array[0]))
    exp_stat[Metric.LATENCY][Stat.MEAN][Stat.ERR95].append(float(stat_array[1]))
    exp_stat[Metric.LATENCY][Stat.MEDIAN][Stat.MEAN].append(float(stat_array[2]))
    exp_stat[Metric.LATENCY][Stat.MEDIAN][Stat.ERR95].append(float(stat_array[3]))
    exp_stat[Metric.LATENCY][Stat.P95th][Stat.MEAN].append(float(stat_array[4]))
    exp_stat[Metric.LATENCY][Stat.P95th][Stat.ERR95].append(float(stat_array[5]))
    exp_stat[Metric.LATENCY][Stat.P99th][Stat.MEAN].append(float(stat_array[6]))
    exp_stat[Metric.LATENCY][Stat.P99th][Stat.ERR95].append(float(stat_array[7]))
    exp_stat[Metric.THROUGHPUT][Stat.MEAN].append(float(stat_array[8]))
    exp_stat[Metric.THROUGHPUT][Stat.ERR95].append(float(stat_array[9]))
    exp_stat[Metric.ACCEPT_RATE][Stat.MEAN].append(float(stat_array[10]))
    exp_stat[Metric.ACCEPT_RATE][Stat.ERR95].append(float(stat_array[11]))
    exp_stat[Metric.ACCEPT_FAST_RATE][Stat.MEAN].append(float(stat_array[12]))
    exp_stat[Metric.ACCEPT_FAST_RATE][Stat.ERR95].append(float(stat_array[13]))
  f.close()
  return exp_stat

#test_dir="../exp-data/s5-dc1-delay0ms/1m-8/"
#print(exp_metric(test_dir))
