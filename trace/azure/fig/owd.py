#!/usr/bin/python

import common as c
     
def owd_diff(data_file, clock_unit, rt_unit, offset_unit):
  # Load raw data
  data_l = c.load_data(f)
  # sending-clock-time, roundtrip-lat, server-clock-time
  send_clock_l, send_rt_l, rev_clock_l, rev_rt_l = c.gen_predict_data(data_file, 'ms', 'ms', 'ms')
