#!/usr/bin/python

import matplotlib.pyplot as plt
import numpy as np
import label as l

def load_predict_rate(f):
  ret = list()
  input_file = open(f, 'r')
  lines = input_file.readlines()
  input_file.close()
  for l in lines:
    if l.startswith("#"):
      continue
    data = l.rstrip('\n').split()
    ret.append((data[0], float(data[1]), float(data[2])))
  return ret

def predict_rate_window_size(data_dir, dc_l):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  for src_dc in dc_l:
    output_file = data_dir + '/predictionrate-' + src_dc + '-window.pdf'
    #output_file = src_dc + '-window.pdf'
    plt.figure(figsize=(8,4))
    for dst_dc in dc_l:
      if dst_dc is src_dc:
        continue
      df = data_dir + '/' + src_dc + '-' + dst_dc + '.txt'
      data = load_predict_rate(df)
      window_size_l = [d[0] for d in data]
      percentile = str(data[0][1])
      predict_rate_l = [100 * d[2] for d in data]
      plt.plot(window_size_l, predict_rate_l, l.line_style[dst_dc], label=dst_dc, \
          color=l.line_color[dst_dc], linewidth=lw, \
          marker=l.point_fmt[dst_dc], markerfacecolor=l.point_color[dst_dc], markersize=ms)
    plt.ylabel('Correct Prediction Rate (%)', fontsize=fs)
    plt.ylim(ymin=0)
    plt.xlabel('Window Size (ms)', fontsize=fs)
    plt.xlim(xmin=0)
    plt.xticks(np.arange(0, 2001, 200), fontsize=fs)
    plt.tick_params(axis='both',direction='in', labelsize=fs)
    plt.title('Host datacenter: ' + src_dc + ', using ' + percentile + 'th percentile latency as predicted latency' , fontsize=fs)
    plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
    plt.savefig(output_file, bbox_inches='tight')

def predict_rate_percentile(data_dir, dc_l):
  fs = 10; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  for src_dc in dc_l:
    output_file = data_dir + '/predictionrate-' + src_dc + '-percentile.pdf'
    #output_file = src_dc + '-percentile.pdf'
    plt.figure(figsize=(8,4))
    for dst_dc in dc_l:
      if dst_dc is src_dc:
        continue
      df = data_dir + '/' + src_dc + '-' + dst_dc + '.txt'
      data = load_predict_rate(df)
      window_size = str(data[0][0])
      percentile_l = [d[1] for d in data]
      predict_rate_l = [100 * d[2] for d in data]
      plt.plot(percentile_l, predict_rate_l, l.line_style[dst_dc], label=dst_dc, \
          color=l.line_color[dst_dc], linewidth=lw, \
          marker=l.point_fmt[dst_dc], markerfacecolor=l.point_color[dst_dc], markersize=ms)
    plt.ylabel('Correct Prediction Rate (%)', fontsize=fs)
    plt.ylim(ymin=0)
    plt.yticks(np.arange(0, 104, 10), fontsize=fs)
    plt.xlabel('Percentile Latency as Predicted Latency (th)', fontsize=fs)
    plt.xlim(xmin=0)
    plt.xticks(np.arange(0, 104, 5), fontsize=fs)
    plt.tick_params(axis='both', direction='in', labelsize=fs)
    plt.title('Host datacenter: ' + src_dc + ', using window size ' + window_size + 'ms', fontsize=fs)
    plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
    plt.savefig(output_file, bbox_inches='tight')

def predict_rate_percentile_multiWindow(output_file, data_dir, src_dc, dst_dc, w_l):
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  count=0
  for window_size in w_l:
    df = data_dir + '/' + src_dc + '-' + dst_dc + '-percentile-w' + str(window_size) + '.txt'
    data = load_predict_rate(df)
    percentile_l = [d[1] for d in data]
    predict_rate_l = [100 * d[2] for d in data]
    plt.plot(percentile_l, predict_rate_l, l.g_line_style[count%len(l.g_line_style)], \
        label='Window Size ' + str(window_size) + 'ms', \
        color=l.g_line_color[count%len(l.g_line_color)], linewidth=lw, \
        marker=l.g_point_fmt[count%len(l.g_point_fmt)], \
        markerfacecolor=l.g_point_color[count%len(l.g_point_color)], markersize=ms)
    count += 1
  plt.ylabel('Correct Prediction Rate (%)', fontsize=fs)
  plt.ylim(ymin=0)
  plt.yticks(np.arange(0, 104, 10), fontsize=fs)
  plt.xlabel('Percentile Delay in Network Measurements for Prediction (th)', fontsize=fs)
  plt.xlim(xmin=0)
  plt.xticks(np.arange(0, 104, 10), fontsize=fs)
  plt.tick_params(axis='both', direction='in', labelsize=fs)
  ##plt.title(src_dc + '-' + dst_dc, fontsize=fs)
  #plt.legend(loc='lower right', ncol=1, fontsize=(fs-1))
  plt.legend(loc='upper left', ncol=1, fontsize=(fs-1))
  plt.savefig(output_file, bbox_inches='tight')

data_dir='./'
w_l = [100, 200, 400, 600, 800, 1000]
src_dc, dst_dc = 'eastus2', 'westus2'
predict_rate_percentile_multiWindow(src_dc+'-'+dst_dc+'-predictrate.pdf', data_dir, src_dc, dst_dc, w_l) 
