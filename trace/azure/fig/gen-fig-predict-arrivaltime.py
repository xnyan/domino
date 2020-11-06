#!/usr/bin/python

import matplotlib.pyplot as plt
import numpy as np
import label as l
import common

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
    if c == 3:
        # On-time rate
        return 100.0 - float(data[2]) #%
    c = c+1
  return 0

def arrival_time_predict_rate_percentile_multiWindow(output_file, data_dir, src_dc, dst_dc, p_l, w_l):
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  count=0
  for w_s in w_l:
    predict_rate_l = list()
    for pth in p_l:
      #eastus2-westus2.log.txt.owd-pth99.0-window1000ms
      df=data_dir+'/'+src_dc+'-'+dst_dc+'.log.txt.owd-pth'+str(pth)+".0-window"+str(w_s)+'ms'
      predict_rate_l.append(load_owd_stat(df))
    #percentile_l = [d[1] for d in data]
    #predict_rate_l = [100 * d[2] for d in data]
    plt.plot(p_l, predict_rate_l, l.g_line_style[count%len(l.g_line_style)], \
        label='Window Size ' + str(w_s) + 'ms', \
        color=l.g_line_color[count%len(l.g_line_color)], linewidth=lw, \
        marker=l.g_point_fmt[count%len(l.g_point_fmt)], \
        markerfacecolor=l.g_point_color[count%len(l.g_point_color)], markersize=ms)
    count += 1
  plt.ylabel('Correct Prediction Rate (%)', fontsize=fs)
  plt.ylim(ymin=0,ymax=108)
  plt.yticks(np.arange(0, 104, 10), fontsize=fs)
  #plt.xlabel('Percentile Delay in Network Measurements for Prediction (th)', fontsize=fs)
  #plt.xlabel('Percentile Value from Measurement Data (th)', fontsize=fs)
  plt.xlabel('n-th Percentile Value in Network Measurements', fontsize=fs)
  plt.xlim(xmin=0)
  plt.xticks(np.arange(0, 104, 10), fontsize=fs)
  plt.tick_params(axis='both', direction='in', labelsize=fs)
  ##plt.title(src_dc + '-' + dst_dc, fontsize=fs)
  #plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
  plt.legend(loc='upper left', ncol=1, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')

data_dir='./owd-globe'
pth_l = np.arange(5,100, 5)
p_l=list()
for p in pth_l:
  p_l.append(p)
p_l.append(99)
print p_l
w_l = [100, 200, 400, 600, 800, 1000]
src_dc, dst_dc = 'eastus2', 'westus2'
arrival_time_predict_rate_percentile_multiWindow(src_dc+'-'+dst_dc+'-arrivaltime-predictrate.pdf', data_dir, src_dc, dst_dc, p_l, w_l) 
