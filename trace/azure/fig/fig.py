#!/usr/bin/python

import matplotlib.pyplot as plt
import numpy as np
import common as c

#Line following clock time
def rt_time_line(clock_l, clock_unit, rt_l, rt_unit, output_file):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  plt.plot(clock_l, rt_l, '--', linewidth=lw, label="Roundtrip Latency", color='black')
  plt.ylabel('Roundtrip Latency (' + rt_unit + ')', fontsize=fs)
  #plt.ylim(ymin=50, ymax=100)
  plt.ylim(ymin=0)
  plt.xlabel('Measuring Clock (' + clock_unit + ')', fontsize=fs)
  plt.xlim(xmin=0)
  #plt.xticks(np.arange(0, 700, 100), fontsize=fs)
  #plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')

def offset_time_line(clock_l, clock_unit, offset_l, offset_unit, output_file):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  plt.plot(clock_l, offset_l, '--', linewidth=lw, label="Timeoffset", color='black')
  plt.ylabel('Clock Time Offset (' + offset_unit + ')', fontsize=fs)
  plt.ylim(ymin=0)
  plt.xlabel('Measuring Clock (' + clock_unit + ')', fontsize=fs)
  plt.xlim(xmin=0)
  plt.savefig(output_file, bbox_inches='tight')


def time_line(f, sent_time_order = True, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, offset_l = c.gen_data_by_send_clock(f, sent_time_order, clock_unit, rt_unit, offset_unit)
  #Plot
  rt_time_line(clock_l, clock_unit, rt_l, rt_unit, f+"-rt.pdf")
  offset_time_line(clock_l, clock_unit, offset_l, offset_unit, f+"-offset.pdf")
  clock_l, rt_l, offset_l = list(), list(), list()

#Hisgtram as distribution
def rt_dist (rt_l, rt_unit, bin_n, output_file):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  if bin_n is 0:
    plt.hist(rt_l, color = 'blue', edgecolor = 'black')
  else:
    plt.hist(rt_l, color = 'blue', edgecolor = 'black', bins = bin_n)
  plt.xlabel('Roundtrip Latency (' + rt_unit + ')', fontsize=fs)
  plt.ylabel('Count', fontsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def dist(f, bin_n = 0, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, offset_l = c.gen_data_by_send_clock(f, False, clock_unit, rt_unit, offset_unit)
  rt_dist(rt_l, rt_unit, bin_n, f+"-rt-dist.pdf")
  ## Prints all of the hist data
  rt_h = c.hist(rt_l)
  for rt in rt_h:
    print rt
  clock_l, rt_l, offset_l = list(), list(), list()

#CDF
def rt_cdf(rt_l, rt_unit, output_file):
  lat_l = np.sort(np.asarray(rt_l))
  Y = (np.arange(len(lat_l)) + 1) / float(len(lat_l))
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  plt.plot(lat_l, Y, '--', linewidth=lw, label=output_file, color='black')
  plt.ylabel('CDF', fontsize=fs)
  plt.ylim(ymin=0, ymax=1.1)
  plt.yticks(np.arange(0, 1.01, 0.2), fontsize=fs)
  plt.xlabel('Roundtrip Latency (' + rt_unit + ')', fontsize=fs)
  plt.xlim(xmin=0)
  plt.tick_params(axis='both',direction='in')
  plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
  lb_color='blue'
  plt.axhline(y=0.5, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.5, "{:.1f}".format(0.5), color=lb_color, ha='right', va='center', fontsize=fs)
  plt.axhline(y=0.95, color=lb_color, linestyle='-.', linewidth=0.5)
  plt.text(0, 0.95, "{:.2f}".format(0.95), color=lb_color, ha='right', va='center', fontsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def cdf(f, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, offset_l = c.gen_data_by_send_clock(f, False, clock_unit, rt_unit, offset_unit)
  rt_cdf(rt_l, rt_unit, f+"-rt-cdf.pdf")
  clock_l, rt_l, offset_l = list(), list(), list()

