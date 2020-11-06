#!/usr/bin/python

from collections import deque
import numpy as np
import common as c

# NOTE: window size must have the same unit as in clock_l
# is_percentile_scale = false to  use the percentile latency from the list
# is_percentile_scale = true to use the scaled percentile latency (may not be from the list)
def correctPredict(window_size, predict_p, is_percentile_scale, rev_clock_l, rev_rt_l, send_clock_l, send_rt_l):
  success, win_s, win_e = 0, 0, 0
  for i in range(0, len(send_clock_l)):
    # Sending time and the actual latency
    send_clock, send_rt = send_clock_l[i], send_rt_l[i]
    # Finds out the prediction window
    while win_s < len(rev_clock_l) and rev_clock_l[win_s] + window_size < send_clock:
      win_s += 1
    while win_e < len(rev_clock_l) and rev_clock_l[win_e] < send_clock:
      win_e += 1
    #if win_e >= len(rev_clock_l):
    #  print(send_clock, win_s, win_e, rev_clock_l[win_s], rev_clock_l[win_e - 1])
    #else:
    #  print(send_clock, win_s, win_e, rev_clock_l[win_s], rev_clock_l[win_e])
    # Nothing in the window
    if win_e <= win_s:
      continue
    # Calculates the predicted latency
    if is_percentile_scale:
      predict_rt = np.percentile(rev_rt_l[win_s : win_e], predict_p) 
    else:
      predict_rt = c.percentile(rev_rt_l[win_s : win_e], predict_p)
    if predict_rt > send_rt:
      success += 1
  return success, success * 1.0 / len(rev_clock_l)

def cal_predict_rate(f, window_size, percentile, is_percentile_scale):
  send_clock_l, send_rt_l, rev_clock_l, rev_rt_l = c.gen_predict_data(f, 'ms', 'ms', 'ms')
  _, rate = correctPredict(window_size, percentile, is_percentile_scale, rev_clock_l, rev_rt_l, send_clock_l, send_rt_l)
  return rate

def get_predict_rate(send_clock_l, send_rt_l, rev_clock_l, rev_rt_l, window_size, percentile, is_percentile_scale):
  _, rate = correctPredict(window_size, percentile, is_percentile_scale, rev_clock_l, rev_rt_l, send_clock_l, send_rt_l)
  return rate
