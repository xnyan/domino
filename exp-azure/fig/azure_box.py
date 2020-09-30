#!/usr/bin/python

import label
import common as c
import matplotlib.pyplot as plt
import numpy as np

def lat_box(output_file, exp_dir_map, p_list, p_tick_label, p_label, whisker = [5, 95], ymin = -1, ymax = 0, exp_n = 10):
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
  plt.boxplot(lat_box, whis=whisker, showfliers=False) # No show of outliers
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.xlabel(p_label,fontsize=fs)
  x_tick_idx = range(1, len(p_tick_label)+1)
  plt.xticks(x_tick_idx, p_tick_label)
  plt.ylabel('Commit Latency (ms)', fontsize=fs)
  if ymin >= 0:
    plt.ylim(ymin=ymin)
  if ymax > 0:
    plt.ylim(ymax=ymax)
  plt.savefig(output_file, bbox_inches='tight')

def horizon_lat_box(output_file, exp_dir_map, p_list, p_tick_label, p_label, whisker = [5, 95], xmin = -1, xmax = 0, exp_n = 10):
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
  plt.xlabel('Commit Latency (ms)', fontsize=fs)
  if xmin >= 0:
    plt.xlim(xmin=xmin)
  if xmax > 0:
    plt.xlim(xmax=xmax)
  plt.savefig(output_file, bbox_inches='tight')
