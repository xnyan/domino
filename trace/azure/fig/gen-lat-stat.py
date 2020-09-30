#!/usr/bin/python

import common as c

dat_dir='network-opt-D4v3/az-na-10dc-24h-202005071450-202005081450'
##dc_l = ['westus', 'eastus', 'northcentralus', 'southcentralus', 'centralus', 'eastus2', 'canadacentral', 'canadaeast', 'westcentralus', 'westus2']
dc_l = ['westus', 'northcentralus', 'southcentralus', 'centralus', 'eastus2', 'canadacentral', 'canadaeast', 'westcentralus', 'westus2']
#dc_l = ['eastus2', 'westus2']

dat_dir='network-opt-D4v3/az-global-8dc-na-eu-ap-24h-202005170045-202005180045'
#dc_l = ['westus2', 'eastus2', 'switzerlandnorth', 'australiaeast', 'southeastasia', 'eastasia', 'francecentral', 'koreasouth']
dc_l = ['westus2', 'eastus2', 'australiaeast', 'southeastasia', 'eastasia', 'francecentral']

print "#srcDC dstDC avg err95 std median p95th p99th min max" 
stat_list = [c.Stat.MEAN, c.Stat.ERR95, c.Stat.STDEV, c.Stat.MEDIAN, c.Stat.P95th, c.Stat.P99th, c.Stat.MIN, c.Stat.MAX]
for src in dc_l:
  for dst in dc_l:
    if dst is src:
      continue
    df=dat_dir+'/' + src+'-'+dst+'.log.txt'
    #print "#"+str(df)
    _, rt_l, _ = c.gen_data_by_send_clock(df, False, 'ms', 'ms', 'ms')
    rt_stat = c.cal_stats(rt_l)
    output = [src, dst]
    for stat in stat_list:
      output.append(rt_stat[stat])
    print output
