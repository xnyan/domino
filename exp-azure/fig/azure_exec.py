#!/usr/bin/python

import label
import common as c
import matplotlib.pyplot as plt
import numpy as np

def lat_cdf(output_file, p_exp_map, p_list, xmax=0, exp_n = 10, txt_map = None):
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  for p in p_list:
    exp_dir = p_exp_map[p][1]
    lat_list = c.get_exp_lat(exp_dir, exp_n)
    lat_list = np.sort(np.asarray(lat_list)/1000000.0)
    print (p, len(lat_list))#, lat_list[0:10])
    Y = (np.arange(len(lat_list)) + 1) / float(len(lat_list))
    plt.plot(lat_list, Y, label.line_style[p], linewidth=lw, \
      label=p_exp_map[p][0], color=label.line_color[p])
      #label=label.protocol_label[p], color=label.line_color[p])
  plt.ylabel('CDF', fontsize=fs)
  plt.ylim(ymin=0, ymax=1.1)
  plt.yticks(np.arange(0, 1.01, 0.2), fontsize=fs)
  plt.xlabel('Execution Latency (ms)', fontsize=fs)
  plt.xlim(xmin=0)
  if xmax > 0:
    plt.xlim(xmax=xmax)
  #plt.xticks(np.arange(0, 700, 100), fontsize=fs)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
  lb_color='blue'
  plt.axhline(y=0.5, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.5, "{:.1f}".format(0.5), color=lb_color, ha='right', va='center', fontsize=fs)
  plt.axhline(y=0.95, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.95, "{:.2f}".format(0.95), color=lb_color, ha='right', va='center', fontsize=fs)
  if txt_map != None:
    for txt in txt_map.keys():
      plt.text(txt_map[txt][0], txt_map[txt][1], txt, color='black', ha='center', va='center', fontsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def custom_lat_cdf(output_file, exp_dir_map, p_list, format_table, xmax=0, exp_n = 10):
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  for p in p_list:
    exp_dir = exp_dir_map[p]
    lat_list = c.get_exp_lat(exp_dir, exp_n)
    lat_list = np.sort(np.asarray(lat_list)/1000000.0)
    print (p, len(lat_list))#, lat_list[0:10])
    Y = (np.arange(len(lat_list)) + 1) / float(len(lat_list))
    plt.plot(lat_list, Y, format_table[p][1], linewidth=lw, \
      label=format_table[p][0], color=format_table[p][2])
  plt.ylabel('CDF', fontsize=fs)
  plt.ylim(ymin=0, ymax=1.1)
  plt.yticks(np.arange(0, 1.01, 0.2), fontsize=fs)
  plt.xlabel('Execution Latency (ms)', fontsize=fs)
  plt.xlim(xmin=0)
  if xmax > 0:
    plt.xlim(xmax=xmax)
  #plt.xticks(np.arange(0, 700, 100), fontsize=fs)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.legend(loc='upper left', ncol=1, fontsize=(fs-1))
  lb_color='blue'
  plt.axhline(y=0.5, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.5, "{:.1f}".format(0.5), color=lb_color, ha='right', va='center', fontsize=fs)
  plt.axhline(y=0.95, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.95, "{:.2f}".format(0.95), color=lb_color, ha='right', va='center', fontsize=fs)
  #plt.text(10, 1.01, dc, color=lb_color, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')

def lat_box(output_file, exp_dir_map, p_list, p_tick_label, x_label, whisker = [5, 95], ymin = -1, ymax=0, exp_n = 10):
  print output_file
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  lat_box = list()
  for p in p_list:
    exp_dir = exp_dir_map[p]
    lat_list = c.get_exp_lat(exp_dir, exp_n)
    lat_list = np.sort(np.asarray(lat_list)/1000000.0)
    stat = c.cal_stats(lat_list)
    print (p, stat[c.Stat.MEDIAN], stat[c.Stat.P95th], stat[c.Stat.P99th], stat[c.Stat.MEAN], stat[c.Stat.ERR95])
    print (p, len(lat_list))#, lat_list[0:10])
    lat_box.append(lat_list)
  #plt.boxplot(lat_box, whis=[5, 95], showfliers=False) # No show of outliers
  plt.boxplot(lat_box, whis=whisker, showfliers=False) # No show of outliers
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  if x_label != "":
    plt.xlabel(x_label,fontsize=fs)
  x_tick_idx = range(1, len(p_tick_label)+1)
  plt.xticks(x_tick_idx, p_tick_label)
  #plt.xticks(x_tick_idx, p_tick_label, rotation=30)
  plt.ylabel('Execution Latency (ms)', fontsize=fs)
  if ymin >= 0:
    plt.ylim(ymin=ymin)
  if ymax > 0:
    plt.ylim(ymax=ymax)
  plt.savefig(output_file, bbox_inches='tight')

def horizon_lat_box(output_file, exp_dir_map, p_list, p_tick_label, p_label, whisker = [5, 95], xmin = -1, xmax=0, exp_n = 10):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  #plot
  plt.figure(figsize=(8,4))
  lat_box = list()
  for p in p_list:
    exp_dir = exp_dir_map[p]
    lat_list = c.get_exp_lat(exp_dir, exp_n)
    lat_list = np.sort(np.asarray(lat_list)/1000000.0)
    print (p, len(lat_list))#, lat_list[0:10])
    lat_box.append(lat_list)
  plt.boxplot(lat_box, whis=whisker, showfliers=False, vert=False) # No show of outliers
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  if p_label != "":
    plt.ylabel(p_label,fontsize=fs)
  y_tick_idx = range(1, len(p_tick_label)+1)
  plt.yticks(y_tick_idx, p_tick_label)
  plt.gca().invert_yaxis() # invert y axis
  plt.xlabel('Execution Latency (ms)', fontsize=fs)
  if xmin >= 0:
    plt.xlim(xmin=xmin)
  if xmax > 0:
    plt.xlim(xmax=xmax)
  plt.savefig(output_file, bbox_inches='tight')
