#!/usr/bin/python

import label
import common as c
import custom
import matplotlib.pyplot as plt
import numpy as np

def get_metric(idx, exp_dir_map, cl_list, exp_list, exp_n = 10):
  exp_dat_map = custom.filter_exp_dat(exp_dir_map, cl_list, exp_list, exp_n=10)
  exp_stat_map = custom.get_exp_metrics(exp_dat_map, exp_list)
  exp_avg_map = custom.cal_avg_metrics(exp_stat_map, exp_list)
  tail_list, err_list = list(), list()
  for e in exp_list:
    tail_list.append(exp_avg_map[e][idx][0])
    err_list.append(exp_avg_map[e][idx][1])
  return tail_list, err_list
  #tail_list, err_list = list(), list()
  #for e in exp_list:
  #  exp_dir = exp_dir_map[e]
  #  stat = c.exp_metric(exp_dir, exp_n)
  #  tail_list.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.MEAN]/1000000.0)
  #  err_list.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.ERR95]/1000000.0)
  #return tail_list, err_list

def plot_lat_bar_and_line(output_file, y_label, x_tick_l, y_map, y_list, line_list, line_map):
  fs = 14; bw = 0.7; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  fig, ax= plt.subplots(figsize=(8,4))
  x_pos = np.arange(0, len(x_tick_l) * len(y_list), len(y_list))
  count = -1 * len(y_list) / 2
  for y in y_list:
    ax.bar(x_pos + bw * count, y_map[y][0], bw, yerr=y_map[y][1], align='edge', edgecolor='black',\
        label=y_map[y][2], hatch=y_map[y][3], color=y_map[y][4], capsize=cs)
    count += 1
  ax.set_xlabel('Additional Delay (ms)', fontsize=fs)
  ax.set_xticks(x_pos)
  ax.set_xticklabels(x_tick_l)
  ax.set_ylabel(y_label, fontsize=fs)
  #ymax = 0
  for l in line_list:
    plt.axhline(y=line_map[l][0], color=line_map[l][2], linestyle=line_map[l][3], linewidth=0.5)
    plt.text(4.5, line_map[l][0]+7, line_map[l][1], color=line_map[l][2], ha='center', va='center', fontsize=(fs-3))
  #  if line_map[l][0] > ymax:
  #    ymax = line_map[l][0]
  #ymax += 30
  #ax.set_ylim(ymin=0, ymax=ymax)
  ax.set_ylim(ymin=0)
  plt.legend(loc='upper right', ncol=1, fontsize=(fs-3))
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def lat_99th_bar(output_file, exp_dir_map, exp_list, exp_n = 10):
  print output_file
  print exp_n
  label_list, tail_list, err_list = list(), list(), list()
  for e in exp_list:
    exp_label, exp_dir = exp_dir_map[e][0], exp_dir_map[e][1]
    stat = c.exp_metric(exp_dir, exp_n)
    tail_list.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.MEAN]/1000000.0)
    err_list.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.ERR95]/1000000.0)
    label_list.append(exp_label)
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  fig, ax= plt.subplots(figsize=(8,4))
  y_pos = np.arange(len(label_list))
  ax.barh(y_pos, tail_list, xerr=err_list, align='center', edgecolor='black', capsize=cs)
  ax.set_yticks(y_pos)
  ax.set_yticklabels(label_list)
  ax.invert_yaxis()  # labels read top-to-bottom
  ax.set_xlabel('99th Percentile Commit Latency (ms)')
  ax.set_xlim(xmin=0)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def lat_95th_bar(output_file, exp_dir_map, exp_list, exp_n = 10):
  print output_file
  print exp_n
  label_list, tail_list, err_list = list(), list(), list()
  for e in exp_list:
    exp_label, exp_dir = exp_dir_map[e][0], exp_dir_map[e][1]
    stat = c.exp_metric(exp_dir, exp_n)
    tail_list.append(stat[c.Metric.LATENCY][c.Stat.P95th][c.Stat.MEAN]/1000000.0)
    err_list.append(stat[c.Metric.LATENCY][c.Stat.P95th][c.Stat.ERR95]/1000000.0)
    label_list.append(exp_label)
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  fig, ax= plt.subplots(figsize=(8,4))
  y_pos = np.arange(len(label_list))
  ax.barh(y_pos, tail_list, xerr=err_list, align='center', edgecolor='black', capsize=cs)
  ax.set_yticks(y_pos)
  ax.set_yticklabels(label_list)
  ax.invert_yaxis()  # labels read top-to-bottom
  ax.set_xlabel('95th Percentile Commit Latency (ms)')
  ax.set_xlim(xmin=0)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def lat_max_bar(output_file, exp_dir_map, exp_list, exp_n = 10):
  print output_file
  print exp_n
  label_list, tail_list, err_list = list(), list(), list()
  for e in exp_list:
    exp_label, exp_dir = exp_dir_map[e][0], exp_dir_map[e][1]
    stat = c.exp_metric(exp_dir, exp_n)
    tail_list.append(stat[c.Metric.LATENCY][c.Stat.MAX][c.Stat.MEAN]/1000000.0)
    err_list.append(stat[c.Metric.LATENCY][c.Stat.MAX][c.Stat.ERR95]/1000000.0)
    label_list.append(exp_label)
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  fig, ax= plt.subplots(figsize=(8,4))
  y_pos = np.arange(len(label_list))
  ax.barh(y_pos, tail_list, xerr=err_list, align='center', edgecolor='black', capsize=cs)
  ax.set_yticks(y_pos)
  ax.set_yticklabels(label_list)
  ax.invert_yaxis()  # labels read top-to-bottom
  ax.set_xlabel('Max Commit Latency (ms)')
  ax.set_xlim(xmin=0)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def lat_95th_99th_bar(output_file, exp_dir_map, exp_list, exp_n = 10):
  print output_file
  print exp_n
  label_list, p99_l, p99_err, p95_l, p95_err, = list(), list(), list(), list(), list()
  for e in exp_list:
    exp_label, exp_dir = exp_dir_map[e][0], exp_dir_map[e][1]
    stat = c.exp_metric(exp_dir, exp_n)
    p99_l.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.MEAN]/1000000.0)
    p99_err.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.ERR95]/1000000.0)
    p95_l.append(stat[c.Metric.LATENCY][c.Stat.P95th][c.Stat.MEAN]/1000000.0)
    p95_err.append(stat[c.Metric.LATENCY][c.Stat.P95th][c.Stat.ERR95]/1000000.0)
    label_list.append(exp_label)
  fs = 10; bw = 1; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  fig, ax= plt.subplots(figsize=(8,4))
  y_pos = np.arange(0, len(label_list)*3, 3)
  ax.barh(y_pos + bw * -1, p99_l, bw, xerr=p99_err, align='edge', \
      hatch='' , edgecolor='black', capsize=cs, label='99th percentile')
  ax.barh(y_pos + bw * 0, p95_l, bw, xerr=p95_err, align='edge', \
      color='white', hatch='x', edgecolor='black', capsize=cs, label='95th percentile')
  ax.set_yticks(y_pos)
  ax.set_yticklabels(label_list)
  ax.invert_yaxis()  # labels read top-to-bottom
  ax.set_xlabel('Commit Latency (ms)')
  ax.set_xlim(xmin=0)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.legend(loc='upper right', ncol=1, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')

def lat_median_99th_bar(output_file, exp_dir_map, exp_list, exp_n = 10):
  print output_file
  print exp_n
  label_list, p99_l, p99_err, median_l, median_err, = list(), list(), list(), list(), list()
  for e in exp_list:
    exp_label, exp_dir = exp_dir_map[e][0], exp_dir_map[e][1]
    stat = c.exp_metric(exp_dir, exp_n)
    p99_l.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.MEAN]/1000000.0)
    p99_err.append(stat[c.Metric.LATENCY][c.Stat.P99th][c.Stat.ERR95]/1000000.0)
    median_l.append(stat[c.Metric.LATENCY][c.Stat.MEDIAN][c.Stat.MEAN]/1000000.0)
    median_err.append(stat[c.Metric.LATENCY][c.Stat.MEDIAN][c.Stat.ERR95]/1000000.0)
    label_list.append(exp_label)
  fs = 10; bw = 1; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  fig, ax= plt.subplots(figsize=(8,4))
  y_pos = np.arange(0, len(label_list)*3, 3)
  ax.barh(y_pos + bw * -1, p99_l, bw, xerr=p99_err, align='edge', \
      hatch='' , edgecolor='black', capsize=cs, label='99th percentile')
  ax.barh(y_pos + bw * 0, median_l, bw, xerr=median_err, align='edge', \
      color='white', hatch='x', edgecolor='black', capsize=cs, label='Median')
  ax.set_yticks(y_pos)
  ax.set_yticklabels(label_list)
  ax.invert_yaxis()  # labels read top-to-bottom
  ax.set_xlabel('Commit Latency (ms)')
  ax.set_xlim(xmin=0)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.legend(loc='upper right', ncol=1, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')
