#!/usr/bin/python

import matplotlib.pyplot as plt
import numpy as np
import common as c

def window(window_size, clock_l, rt_l):
  rt_box = list()
  s = clock_l[0]
  rt_window = list()
  for i, c in enumerate(clock_l):
    if c <= s + window_size:
      rt_window.append(rt_l[i])
    else:
      rt_box.append(rt_window)
      s = c
      rt_window = list()
      rt_window.append(rt_l[i])
  if len(rt_window) is not 0:
    rt_box.append(rt_window)
  return rt_box

def slide_window(window_size, slide_size, clock_l, rt_l):
  s, rt_box, rt_window = 0, list(), list()
  while s < len(clock_l):
    lower, upper, slide_clock, slide_clock_set = clock_l[s], clock_l[s] + window_size, clock_l[s] + slide_size, False
    for i in range(s, len(clock_l)):
      if not slide_clock_set and clock_l[i] > slide_clock: 
        s, slide_clock_set = i, True
      if clock_l[i] > upper:
        rt_box.append(rt_window)
        rt_window = list()
        if not slide_clock_set:
          s, slide_clock_set = i, True
        break
      else:
        rt_window.append(rt_l[i])
    if i == len(clock_l)-1 or not slide_clock_set: # Last window
      rt_box.append(rt_window)
      s = len(clock_l)
  return rt_box

# Whisker box for different time windows
def window_box_plot(output_file, window_size, clock_l, rt_l, clock_unit, rt_unit):
  rt_box = window(window_size, clock_l, rt_l)
  #Plot
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  plt.boxplot(rt_box, whis=[5, 95], showfliers=False) # No show of outliers
  #plt.boxplot(rt_box, whis=[5, 95])
  #plt.boxplot(rt_box, whis=[0, 100])
  plt.xlabel('Window Number (No Overlap). Window size:' + str(window_size) + '(' + clock_unit + ')', fontsize=fs)
  #plt.ylim(ymin=0)
  plt.ylabel('Network Roundtrip Delay (' + rt_unit + ')', fontsize=fs)
  plt.savefig(output_file, bbox_inches='tight')
  return rt_box

def window_box(output_file, f, window_size, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, offset_l = c.gen_data_by_send_clock(f, True, clock_unit, rt_unit, offset_unit)
  window_box_plot(output_file, window_size, clock_l, rt_l, clock_unit, rt_unit)

def window_box_range(output_file, f, range_offset, range_length, window_size, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l = c.gen_data_range(f, range_offset, range_length, clock_unit, rt_unit, offset_unit)
  window_box_plot(output_file, window_size, clock_l, rt_l, clock_unit, rt_unit)

# Slinding window Whisker box
def slide_window_box_plot(output_file, window_size, slide_size, clock_l, rt_l, clock_unit, rt_unit):
  rt_box = slide_window(window_size, slide_size, clock_l, rt_l)
  #Plot
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  plt.tick_params(axis='both', labelsize=fs)
  plt.boxplot(rt_box, whis=[5, 95], showfliers=False) # No show of outliers
  #plt.boxplot(rt_box, whis=[5, 95], sym='') # No show of outliers for some older versions
  #plt.boxplot(rt_box, whis=[5, 95])
  #plt.boxplot(rt_box, whis=[0, 100])
  plt.xlabel('Window Number. Window size:'+str(window_size)+'('+clock_unit+')' + 'Sliding szie:'+str(slide_size)+'('+clock_unit+')', fontsize=fs)
  #plt.ylim(ymin=0)
  plt.ylabel('Network Roundtrip Delay (' + rt_unit + ')', fontsize=fs)
  plt.savefig(output_file, bbox_inches='tight')
  return rt_box

def slide_window_box(output_file, f, window_size, slide_size, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, offset_l = c.gen_data_by_send_clock(f, True, clock_unit, rt_unit, offset_unit)
  slide_window_box_plot(output_file, window_size, slide_size, clock_l, rt_l, clock_unit, rt_unit)

def slide_window_box_range(output_file, f, range_offset, range_length, window_size, slide_size, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l = c.gen_data_range(f, range_offset, range_length, clock_unit, rt_unit, offset_unit)
  slide_window_box_plot(output_file, window_size, slide_size, clock_l, rt_l, clock_unit, rt_unit)

def custom_slide_window_box_range(output_file, f, range_offset, range_length, window_size, slide_size, clock_unit = 'min', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l = c.gen_data_range(f, range_offset, range_length, clock_unit, rt_unit, offset_unit)
  rt_box = slide_window(window_size, slide_size, clock_l, rt_l)
  #Plot
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  plt.tick_params(axis='both', labelsize=fs)
  plt.boxplot(rt_box, whis=[5, 95], showfliers=False) # No show of outliers
  #plt.boxplot(rt_box, whis=[5, 95], sym='') # No show of outliers for some older versions
  #plt.boxplot(rt_box, whis=[5, 95])
  #plt.boxplot(rt_box, whis=[0, 100])
  #x_tick_idx = range(1, len(rt_box)+1, (len(rt_box)+1)/17)
  x_tick_idx = range(1, len(rt_box)+1, (len(rt_box)+1)/10)
  x_tick_label = list()
  for x in x_tick_idx:
    if x%2 == 0:
      x_tick_label.append(float(x)/2.0+0.5)
    else:
      x_tick_label.append(x/2+1)
  print x_tick_idx
  print x_tick_label
  plt.xticks(x_tick_idx, x_tick_label)
  #plt.xlabel('Window Number. Window size:'+str(window_size)+'('+clock_unit+')' + 'Sliding szie:'+str(slide_size)+'('+clock_unit+')', fontsize=fs)
  #plt.xlabel('Window Number. (Window size '+str(window_size)+' '+clock_unit+'. Two adjacent windows overlap in '+str(slide_size)+' '+clock_unit+')', fontsize=fs)
  plt.xlabel('Time ('+clock_unit+')', fontsize=fs+1)
  #plt.ylim(ymin=0)
  plt.ylabel('Network Roundtrip Delay (' + rt_unit + ')', fontsize=fs+1)
  plt.savefig(output_file, bbox_inches='tight')

