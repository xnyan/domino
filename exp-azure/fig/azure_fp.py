#!/usr/bin/python

import label
import common as c
import matplotlib.pyplot as plt
import numpy as np

def lat_cdf(output_file, exp_dir_map, p_list, xmax=0, exp_n = 10):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  for p in p_list:
    exp_dir = exp_dir_map[p]
    lat_list = c.get_exp_lat(exp_dir, exp_n)
    lat_list = np.sort(np.asarray(lat_list)/1000000.0)
    print (p, len(lat_list))#, lat_list[0:10])
    Y = (np.arange(len(lat_list)) + 1) / float(len(lat_list))
    plt.plot(lat_list, Y, label.line_style[p], linewidth=lw, \
      label=label.protocol_label[p], color=label.line_color[p])
  plt.ylabel('CDF', fontsize=fs)
  plt.ylim(ymin=0, ymax=1.2)
  plt.yticks(np.arange(0, 1.01, 0.2), fontsize=fs)
  plt.xlabel('Commit Latency (ms)', fontsize=fs)
  plt.xlim(xmin=0)
  if xmax > 0:
    plt.xlim(xmax=xmax)
  #plt.xticks(np.arange(0, 700, 100), fontsize=fs)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.legend(loc='upper left', ncol=2, fontsize=(fs-1))
  lb_color='blue'
  plt.axhline(y=0.5, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.5, "{:.1f}".format(0.5), color=lb_color, ha='right', va='center', fontsize=fs)
  plt.axhline(y=0.95, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.95, "{:.2f}".format(0.95), color=lb_color, ha='right', va='center', fontsize=fs)
  #plt.text(10, 1.01, dc, color=lb_color, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')

def custom_lat_cdf(output_file, exp_dir_map, exp_list, exp_label, xmax=0, exp_n = 10):
  print output_file
  fs = 15; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  for exp in exp_list:
    exp_dir = exp_dir_map[exp]
    lat_list = c.get_exp_lat(exp_dir, exp_n)
    lat_list = np.sort(np.asarray(lat_list)/1000000.0)
    stat = c.cal_stats(lat_list)
    print (exp, stat[c.Stat.MEDIAN], stat[c.Stat.P95th], stat[c.Stat.P99th], stat[c.Stat.MEAN], stat[c.Stat.ERR95])
    print (exp, len(lat_list))#, lat_list[0:10])
    Y = (np.arange(len(lat_list)) + 1) / float(len(lat_list))
    plt.plot(lat_list, Y, exp_label[exp][0], linewidth=lw, \
      label=exp_label[exp][1], color=exp_label[exp][2])
  plt.ylabel('CDF', fontsize=fs)
  plt.ylim(ymin=0, ymax=1.35)
  plt.yticks(np.arange(0, 1.01, 0.2), fontsize=fs)
  plt.xlabel('Commit Latency (ms)', fontsize=fs)
  plt.xlim(xmin=0)
  if xmax > 0:
    plt.xlim(xmax=xmax)
  plt.xticks(fontsize=fs)
  #plt.xticks(np.arange(0, 700, 100), fontsize=fs)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.legend(loc='upper left', ncol=2, fontsize=(fs-1))
  lb_color='blue'
  plt.axhline(y=0.5, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.5, "{:.1f}".format(0.5), color=lb_color, ha='right', va='center', fontsize=fs)
  plt.axhline(y=0.95, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.95, "{:.2f}".format(0.95), color=lb_color, ha='right', va='center', fontsize=fs)
  #plt.text(10, 1.01, dc, color=lb_color, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')
