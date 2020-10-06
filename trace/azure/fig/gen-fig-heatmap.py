#!/usr/bin/python

import heat as h

dir='trace-azure-globe-6dc-24h-202005170045-202005180045'
src_dc_l = ['eastus2']
dst_dc_l = ['westus2', 'francecentral', 'australiaeast']
for src in src_dc_l:
  for dst in dst_dc_l:
    if dst is src:
      continue
    df=dir+'/' + src+'-'+dst+'.log.txt'
    print df
    ##h.lat_heat_map(df, 1, 1, 12, 1, 'min', 'ms', 'ms')
    #h.vertical_lat_heat_map(df, 1, 1, 12, 1, 'min', 'ms', 'ms')
    #h.lat_heat_map(df, 1, 1, 15, 1, 'min', 'ms', 'ms')
    h.lat_heat_map(df, 1, 1, 11, 1, 'min', 'ms', 'ms')

