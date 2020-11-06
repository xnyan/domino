#!/usr/bin/python

import common as c
import operator
    
def gen_owd_dat(data_l):
  # Delays that are orderd based on sending timestamps, considered as the actual delay
  sending_clock_l, arrival_clock_l = list(), list()
  for i, d in enumerate(data_l):
    sending_clock_l.append(d[0])
    arrival_clock_l.append(d[2])
  # The delays that are in the receiving order, which are used for predictions.
  # This is because at a sending time, the immediate previous probing result
  # may not be available yet due to network delays.
  rev_clock_l, rev_rt_l, rev_offset_l = list(), list(), list()
  data_l.sort(key = operator.itemgetter(1))
  for i, d in enumerate(data_l):
    rev_clock_l.append(d[1])
    rev_rt_l.append(d[1] - d[0])
    rev_offset_l.append(d[2] - d[0])
  return sending_clock_l, arrival_clock_l, rev_clock_l, rev_rt_l, rev_offset_l

def owd_diff(data_file, predict_p = 95.0, window_size = 1000):
  # Load raw data
  data_l = c.load_data(data_file)
  send_clock_l, arrival_clock_l, rev_clock_l, rev_rt_l, rev_offset_l = gen_owd_dat(data_l)
  w_size = window_size * 1000000 # ms to ns
  #for i in range(0, 10):
  # print (send_clock_l[i], arrival_clock_l[i], rev_clock_l[i], rev_rt_l[i], rev_offset_l[i])
  # Calculates time diffs
  ontime_diff_l, expire_diff_l, win_s, win_e = list(), list(), 0, 0
  for i in range(0, len(send_clock_l)):
    # Sending time and the actual roundtrip latency
    send_clock, arrival_clock = send_clock_l[i], arrival_clock_l[i]
    #print (i, send_clock, arrival_clock)
    # Finds out the prediction window
    while win_s < len(rev_clock_l) and rev_clock_l[win_s] + w_size < send_clock:
      win_s += 1
    while win_e < len(rev_clock_l) and rev_clock_l[win_e] < send_clock:
      win_e += 1
    #print (win_s, win_e)
    if win_e <= win_s:
      continue # Nothing in the window
    # Calculates the predicted latency
    predict_rt = c.percentile(rev_rt_l[win_s : win_e], predict_p)
    predict_owd = predict_rt / 2 + 1
    predict_arrival_time = send_clock + predict_owd 
    if predict_arrival_time >= arrival_clock:
      ontime_diff_l.append(predict_arrival_time - arrival_clock)
    else:
      expire_diff_l.append(arrival_clock - predict_arrival_time)
  return ontime_diff_l, expire_diff_l

# Roundup half delay to 1ms
def owd_diff_ms(data_file, predict_p = 95.0, window_size = 1000):
  # Load raw data
  data_l = c.load_data(data_file)
  send_clock_l, arrival_clock_l, rev_clock_l, rev_rt_l, rev_offset_l = gen_owd_dat(data_l)
  w_size = window_size * 1000000 # ms to ns
  #for i in range(0, 10):
  # print (send_clock_l[i], arrival_clock_l[i], rev_clock_l[i], rev_rt_l[i], rev_offset_l[i])
  # Calculates time diffs
  ontime_diff_l, expire_diff_l, win_s, win_e = list(), list(), 0, 0
  for i in range(0, len(send_clock_l)):
    # Sending time and the actual roundtrip latency
    send_clock, arrival_clock = send_clock_l[i], arrival_clock_l[i]
    #print (i, send_clock, arrival_clock)
    # Finds out the prediction window
    while win_s < len(rev_clock_l) and rev_clock_l[win_s] + w_size < send_clock:
      win_s += 1
    while win_e < len(rev_clock_l) and rev_clock_l[win_e] < send_clock:
      win_e += 1
    #print (win_s, win_e)
    if win_e <= win_s:
      continue # Nothing in the window
    # Calculates the predicted latency
    predict_rt = c.percentile(rev_rt_l[win_s : win_e], predict_p)
    predict_owd = (predict_rt / 2 / 1000000 + 1) * 1000000 # roundup to 1ms
    predict_arrival_time = send_clock + predict_owd 
    if predict_arrival_time >= arrival_clock:
      ontime_diff_l.append(predict_arrival_time - arrival_clock)
    else:
      expire_diff_l.append(arrival_clock - predict_arrival_time)
  return ontime_diff_l, expire_diff_l

# Roundup half delay to 1ms
def owd_diff_timeoffset(data_file, predict_p = 95.0, window_size = 1000):
  # Load raw data
  data_l = c.load_data(data_file)
  send_clock_l, arrival_clock_l, rev_clock_l, rev_rt_l, rev_offset_l = gen_owd_dat(data_l)
  w_size = window_size * 1000000 # ms to ns
  #for i in range(0, 10):
  # print (send_clock_l[i], arrival_clock_l[i], rev_clock_l[i], rev_rt_l[i], rev_offset_l[i])
  # Calculates time diffs
  ontime_diff_l, expire_diff_l, win_s, win_e = list(), list(), 0, 0
  for i in range(0, len(send_clock_l)):
    # Sending time and the actual roundtrip latency
    send_clock, arrival_clock = send_clock_l[i], arrival_clock_l[i]
    #print (i, send_clock, arrival_clock)
    # Finds out the prediction window
    while win_s < len(rev_clock_l) and rev_clock_l[win_s] + w_size < send_clock:
      win_s += 1
    while win_e < len(rev_clock_l) and rev_clock_l[win_e] < send_clock:
      win_e += 1
    #print (win_s, win_e)
    if win_e <= win_s:
      continue # Nothing in the window
    # Calculates the predicted latency by using the timestamp offset
    predict_offset = c.percentile(rev_offset_l[win_s : win_e], predict_p)
    predict_arrival_time = send_clock + predict_offset
    if predict_arrival_time >= arrival_clock:
      ontime_diff_l.append(predict_arrival_time - arrival_clock)
    else:
      expire_diff_l.append(arrival_clock - predict_arrival_time)
  return ontime_diff_l, expire_diff_l



