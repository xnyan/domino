#!/usr/bin/python

import common as c
import roundtrip as p
import numpy as np
# Fixed Percentile, varies window size
def predict_window(input_dir, dc_l, output_dir, percentile, window_l, is_percentile_scale):
  for src in dc_l:
    for dst in dc_l:
      if dst is src:
        continue
      data_file=input_dir+'/' + src+'-'+dst+'.log.txt'
      send_clock_l, send_rt_l, rev_clock_l, rev_rt_l = c.gen_predict_data(data_file, 'ms', 'ms', 'ms')
      out=open(output_dir+'/' + src+'-'+dst+'-window-p'+str(percentile)+'.txt', 'w')
      out.write('#window_size(ms) percentile prediction_rate\n')
      for window_size in window_l:
        #rate=p.cal_predict_rate(data_file, window_size, percentile, is_percentile_scale)
        rate=p.get_predict_rate(send_clock_l, send_rt_l, rev_clock_l, rev_rt_l, window_size, percentile, is_percentile_scale)
        print (data_file, window_size, percentile, rate)
        out.write(str(window_size) + ' ' + str(percentile) + ' ' + str(rate) + '\n')
        out.flush()
      out.close()

# Fixed window size, vaires percentile
def predict_percentile(input_dir, src_dc_l, dst_dc_l, output_dir, window_size, p_l, is_percentile_scale):
  for src in src_dc_l:
    for dst in dst_dc_l:
      if dst is src:
        continue
      data_file=input_dir+'/' + src+'-'+dst+'.log.txt'
      send_clock_l, send_rt_l, rev_clock_l, rev_rt_l = c.gen_predict_data(data_file, 'ms', 'ms', 'ms')
      out=open(output_dir+'/' + src+'-'+dst+'-percentile-w'+str(window_size)+'.txt', 'w')
      out.write('#window_size(ms) percentile prediction_rate\n')
      for percentile in p_l:
        #rate=p.cal_predict_rate(data_file, window_size, percentile, is_percentile_scale)
        rate=p.get_predict_rate(send_clock_l, send_rt_l, rev_clock_l, rev_rt_l, window_size, percentile, is_percentile_scale)
        print (data_file, window_size, percentile, rate)
        out.write(str(window_size) + ' ' + str(percentile) + ' ' + str(rate) + '\n')
        out.flush()
      out.close()

##### Only for eastus2 to westus2
window_size, percentile, is_percentile_scale = 1000, 95.0, False
p_l = range(5, 100, 5)
p_l.append(99)
p_l.append(100)
w_l = [100, 200, 400, 500, 600, 800, 1000]

input_dir='trace-azure-globe-6dc-24h-202005170045-202005180045'
output_dir = "./"
src_dc_l = ['eastus2']
dst_dc_l = ['westus2']

for window_size in w_l: 
  predict_percentile(input_dir, src_dc_l, dst_dc_l, output_dir, window_size, p_l, is_percentile_scale)
