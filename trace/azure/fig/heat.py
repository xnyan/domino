#!/usr/bin/python

from matplotlib.pyplot import cm
import matplotlib.pyplot as plt
import matplotlib.colors as colors
import numpy as np
import common as c
import operator


# Row: clock, Column: latency
# NOTE: Deprecated
def heat_map_clock_lat(clock_l, rt_l, interval, lat_step, lat_num):
  #rt_hist_l = c.hist(rt_l)
  lat_min = int(min(rt_l))
  hmap, clock_label, k, lat_label = list(), list(), 1, list()
  limit = clock_l[0] + interval
  rt_hist = [0] * (lat_num+1)
  for i, c in enumerate(clock_l):
    if c >= limit:
      hmap.append(rt_hist)
      clock_label.append(k)
      k += 1
      limit = c + interval
      rt_hist = [0] * (lat_num+1)
    lat = int(rt_l[i])
    for j in range(1, lat_num+1):
      if lat < lat_min + lat_step * j:
        rt_hist[j-1] += 1
        break
    if lat >= lat_min + lat_step * lat_num:
      rt_hist[lat_num] += 1
  hmap.append(rt_hist)
  clock_label.append(k)
  # Latency label
  for i in range(0, lat_num):
    lat_label.append(lat_min + lat_step * i)
  lat_label.append('>' + str(lat_min + lat_step * lat_num))
  hmap.reverse()
  clock_label.reverse()
  return hmap, clock_label, lat_label

# Row: latency, Column: clock
# interval: time interval
# interval_overlap: percentage of overlaps between overlaps
# lat_step: latency interval
# lat_up_num: max number of lat_steps above the min latency
# lat_down_num: number of lat_steps below the min latency
def heat_map_lat_clock(clock_l, rt_l, interval, lat_step, lat_up_num, lat_down_num):
  #rt_hist_l = c.hist(rt_l)
  lat_min = int(min(rt_l))
  lat_max = int(lat_min + lat_step * lat_up_num)
  hmap, lat_label = {}, list()
  for i in range(-1 * lat_down_num, lat_up_num + 1):
    lat = int(lat_min + lat_step * i)
    hmap[lat] = list()
    lat_label.append(lat)
  lat_label[0] = '<' + str(lat_label[0])
  lat_label[-1] = '>' + str(lat_label[-1])
  limit = clock_l[0] + interval
  inv_s, inv_e = 0, 0
  for i, cl in enumerate(clock_l):
    if cl >= limit:
      rt_h = c.hist(rt_l[inv_s:i])
      for h_l in hmap.values():
        h_l.append(0)
      for rt_e in rt_h:
        if rt_e[0] >= lat_max:
          hmap[lat_max][-1] += rt_e[1]
        else:
          hmap[rt_e[0]][-1] += rt_e[1]
      inv_s = i
      limit = cl + interval
  if inv_s < len(clock_l):
    rt_h = c.hist(rt_l[inv_s:len(clock_l)])
    for h_l in hmap.values():
      h_l.append(0)
    for rt_e in rt_h:
      if rt_e[0] >= lat_max:
        hmap[lat_max][-1] += rt_e[1]
      else:
        hmap[rt_e[0]][-1] += rt_e[1]
  hmap_l = list()
  for e in sorted(hmap.items(), key = operator.itemgetter(0)):
    hmap_l.append(e[1])
  hmap_l.reverse()
  lat_label.reverse()
  return hmap_l, lat_label

# Referece: https://matplotlib.org/3.1.1/gallery/images_contours_and_fields/image_annotated_heatmap.html
def lat_heat_map(f, interval, lat_step, lat_up_num, lat_down_num, clock_unit = 's', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, _ = c.gen_data_by_send_clock(f, True, clock_unit, rt_unit, offset_unit)
  #hmap, clock_label, lat_label = heat_map_clock_lat(clock_l, rt_l, interval, lat_step, lat_up_num)
  hmap, lat_label = heat_map_lat_clock(clock_l, rt_l, interval, lat_step, lat_up_num, lat_down_num)
  for h in hmap:
    print h
  output_file = f + "-heat.pdf"
  fs = 11; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  fig, ax = plt.subplots()
  ax.tick_params(axis='both', labelsize=fs)
  #plt.matshow(hmap)
  #im = ax.imshow(hmap)
  #im = ax.imshow(hmap, interpolation='nearest', cmap=cm.Greys)
  #im = ax.imshow(hmap, interpolation='nearest', cmap=cm.Reds)
  im = ax.imshow(hmap, aspect='auto', interpolation='nearest', cmap=cm.Reds)
  #im = ax.imshow(hmap, aspect=2, interpolation='nearest', cmap=cm.Greys)
  #im = ax.imshow(hmap, aspect='auto', cmap=cm.Greys)
  plt.xlabel('Time ' + '(' + str(interval) + ' ' + clock_unit + ')', fontsize=fs+2)
  #ax.set_xticks(np.arange(len(clock_label)))
  #ax.set_xticklabels(clock_label)
  ## Rotate the tick labels and set their alignment.
  #plt.setp(ax.get_xticklabels(), rotation=45, ha="right", rotation_mode="anchor")
  plt.ylabel('Network Roundtrip Delay ' + '(' + rt_unit + ')', fontsize=fs+2)
  ax.set_yticks(np.arange(len(lat_label)))
  ax.set_yticklabels(lat_label, fontsize=fs)
  # Loop over data dimensions and create text annotations.
  #for i in range(len(vegetables)):
  #  for j in range(len(farmers)):
  #      text = ax.text(j, i, harvest[i, j], ha="center", va="center", color="w")
  # Creates color bar
  cbarlabel = 'Number of Measured Delays per ' + str(interval) + ' ' +clock_unit
  #cbar = ax.figure.colorbar(im, ax=ax, orientation = 'vertical') # vertical color bar
  #cbar.ax.set_ylabel(cbarlabel, rotation=-90, va="bottom", fontsize=fs+2)
  ##cbar = ax.figure.colorbar(im, ax=ax, orientation = 'horizontal', fraction=0.045) # horizontal color bar
  cbar = ax.figure.colorbar(im, ax=ax, orientation = 'horizontal')
  cbar.ax.set_xlabel(cbarlabel, va="top", fontsize=fs+2)
  cbar.ax.tick_params(axis='both', labelsize=fs)
  #ax.set_title(f, fontsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def vertical_lat_heat_map(f, interval, lat_step, lat_up_num, lat_down_num, clock_unit = 's', rt_unit = 'ms', offset_unit = 'ms'):
  clock_l, rt_l, _ = c.gen_data_by_send_clock(f, True, clock_unit, rt_unit, offset_unit)
  hmap, lat_label = heat_map_lat_clock(clock_l, rt_l, interval, lat_step, lat_up_num, lat_down_num)
  for h in hmap:
    print h
  output_file = f + "-heat.pdf"
  fs = 13; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,3.5))
  fig, ax = plt.subplots(figsize=(8,3.5))
  ax.tick_params(axis='both', labelsize=fs)
  im = ax.imshow(hmap, aspect='auto', interpolation='nearest', cmap=cm.Reds)
  plt.xlabel('Time ' + '(' + str(interval) + ' ' + clock_unit + ')', fontsize=fs)
  plt.ylabel('Network Roundtrip Delay ' + '(' + rt_unit + ')', fontsize=fs)
  ax.set_yticks(np.arange(len(lat_label)))
  ax.set_yticklabels(lat_label, fontsize=fs)
  cbarlabel = '# of Measured Delays per ' + str(interval) + ' ' +clock_unit
  cbar = ax.figure.colorbar(im, ax=ax, orientation = 'vertical') # vertical color bar
  cbar.ax.set_ylabel(cbarlabel, rotation=-90, va="bottom", fontsize=fs)
  cbar.ax.tick_params(axis='both', labelsize=fs)
  plt.savefig(output_file, bbox_inches='tight')

def test():
  hmap = list()
  for i in range(0, 3):
    hmap.append(list())
    for j in range(0, 5):
      hmap[i].append(i*j)
  for h in hmap:
    print h
  hmap = heat_map(0, 0, 0, 0, 0, 0)
  plt.figure(figsize=(8,4))
  plt.matshow(hmap)
  plt.savefig("test-heat-map.pdf", bbox_inches='tight')

