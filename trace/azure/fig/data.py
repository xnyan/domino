#!/usr/bin/python

import common as c
 
def gen_rt(f, rt_unit = 'ms'):
  clock_l, rt_l, offset_l = c.gen_plot_data(f, False, clock_unit = 'min', rt_unit = rt_unit, offset_unit = 'ms')
  output_file = f+ "-rt.dat"
  dat_file = open(output_file, 'w+')
  for rt in rt_l:
    dat_file.write(str(rt) + "\n")
  dat_file.close()

