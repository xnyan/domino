#!/usr/bin/python

import common as c

idx_lat_mean = 0
idx_lat_median = 1
idx_lat_95 = 2
idx_lat_99 = 3
idx_lat_min = 4
idx_lat_max = 5
idx_fpath = 6
idx_fpaxos = 7
idx_thr = 8

metric_map = {
  idx_lat_mean    : ('lat-mean'     , 'Mean Commit Latency (ms)'),\
  idx_lat_median  : ('lat-median'   , 'Median Commit Latency (ms)'),\
  idx_lat_95      : ('lat-95th'     , '95th Percentile Commit Latency (ms)'),\
  idx_lat_99      : ('lat-99th'     , '99th Percentile Commit Latency (ms)'),\
  idx_lat_min     : ('lat-min'      , 'Min Commit Latency (ms)'),\
  idx_lat_max     : ('lat-max'      , 'Max Commit Latency (ms)'),\
  idx_fpath       : ('FastPathRate' , 'Fast Path Success Rate (%)'),\
  idx_fpaxos      : ('FastPaxosRate', 'Requests using Fast Paxos (%)'),\
  idx_thr         : ('Throughput'   , 'Throughput (# of requests per second'),\
  }

# Reteurns: p --> [raw data in each measurements]
def filter_exp_dat(exp_dir_map, cl_list, p_list, exp_n=10):
  exp_dat_map = {}
  for p in p_list:
    exp_dir = exp_dir_map[p]
    exp_dat_map[p] = list()
    for i in range(1, exp_n+1):
      dat_dir = exp_dir + '/' + str(i) + '/'
      exp_dat = list()
      for cl in cl_list:
        cl_file = dat_dir + '/' + cl
        print cl_file
        exp_dat.extend(c.filter_data(cl_file, c.warm_up_time, c.measure_time))
      exp_dat_map[p].append(exp_dat)
  return exp_dat_map

# Returns a list of measurements' Latency (Mean, Median, 95th, 99th, Min, Max),
# FastPath rate, FastPaxos rate, and Throughput
def get_exp_metrics(exp_dat_map, p_list):
  exp_stat_map = {}
  for p in p_list:
    exp_dat = exp_dat_map[p] # one experiment consisting of multiple measurements
    exp_stat = {}
    for i in range(0, idx_thr + 1):
      exp_stat[i] = list()
    for dat_l in exp_dat: # one measurement
      lat_l, fpfc, fpsc, fpr, pc, pr, n = list(), 0, 0, 0, 0, 0, 0
      for dat in dat_l:
        n += 1
        stat = dat.rstrip('\n').split(',')
        is_use_fp, is_accept, is_fast = bool(int(stat[0].strip())), bool(int(stat[1].strip())), bool(int(stat[2].strip()))
        lat = float(stat[3].strip())
        lat_l.append(lat)
        fpfc += int(is_use_fp) & int(is_accept) & int(is_fast)
        fpsc += int(is_use_fp) & int(is_accept) & int((not is_fast))
        fpr  += int(is_use_fp) & int((not is_accept))
        pc  += int((not is_use_fp)) & int(is_accept)
        pr  += int((not is_use_fp)) & int((not is_accept))
      fpn, pn = fpfc + fpsc + fpr, pc + pr
      if n != (fpn + pn):
        print ("Error: request numbers", n, fpn, pn, fpfc, fpsc, fpr, pc, pr)
        exit()
      lat_stat = c.cal_stats(lat_l)
      exp_stat[idx_lat_mean].append(lat_stat[c.Stat.MEAN] / 1000000.0) # from ns to ms
      exp_stat[idx_lat_median].append(lat_stat[c.Stat.MEDIAN] / 1000000.0)
      exp_stat[idx_lat_95].append(lat_stat[c.Stat.P95th] / 1000000.0)
      exp_stat[idx_lat_99].append(lat_stat[c.Stat.P99th] / 1000000.0)
      exp_stat[idx_lat_min].append(lat_stat[c.Stat.MIN] / 1000000.0)
      exp_stat[idx_lat_max].append(lat_stat[c.Stat.MAX] / 1000000.0)
      if fpn == 0:
        exp_stat[idx_fpath].append(0.0)
      else:
        exp_stat[idx_fpath].append(float(fpfc)/float(fpn) * 100.0)
      exp_stat[idx_fpaxos].append(float(fpn)/float(n) * 100.0)
      exp_stat[idx_thr].append(float(n) / c.measure_time)
    exp_stat_map[p] = exp_stat
    #print (p, exp_stat)
  return exp_stat_map

# Returns the average of each metric and its 95% confidence interval
def cal_avg_metrics(exp_stat_map, p_list):
  exp_avg_map = {}
  for p in p_list:
    exp_stat = exp_stat_map[p]
    exp_avg = {}
    for k in exp_stat.keys():
      stat = c.cal_stats(exp_stat[k])
      exp_avg[k] = (stat[c.Stat.MEAN], stat[c.Stat.ERR95])
    exp_avg_map[p] = exp_avg
  return exp_avg_map

