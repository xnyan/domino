#!/usr/bin/python

import sys
import operator
import common as c
import owd as owd
    
#dat_f = "test.log.txt"
#dat_f = "trace-azure-globe-6dc-24h-202005170045-202005180045/eastus2-westus2.log.txt"
#predict_p, window_size = 95.0, 1000

## Command line args
if len(sys.argv) < 6:
    print "Usage: <input_data_dir> <input_data_file_name> <output_dir> <predict_delay_percentile> <window_size (ms)>"
    exit()
input_dir, dat_f, out_dir = sys.argv[1], sys.argv[2], sys.argv[3]
predict_p, window_size = float(sys.argv[4]), int(sys.argv[5])

out_f = dat_f+'-arrivaltime-pth' + str(predict_p) + '-window' + str(window_size) + 'ms.txt'
out=open(out_dir+'/'+out_f, 'w+')

print input_dir+'/'+dat_f

## Raw ns data (for Table 2)
ontime_l, expire_l = owd.owd_diff(input_dir+'/'+dat_f, predict_p, window_size)
## Roundup delay to 1ms 
#ontime_l, expire_l = owd.owd_diff_ms(input_dir+'/'+dat_f, predict_p, window_size)
## Raw ns Timeoffset (for Table 3 and Figure 3)
#ontime_l, expire_l = owd.owd_diff_timeoffset(input_dir+'/'+dat_f, predict_p, window_size)

ontime_stat, expire_stat = c.cal_stats(ontime_l), c.cal_stats(expire_l)
ontime_n, expire_n = len(ontime_l), len(expire_l)
expire_rate = round(expire_n * 1.0 / (expire_n + ontime_n) * 100, 2)

out.write("#Difference between the predicted arrival time and the actual arrival time\n")
out.write("#Mean (ms), error95, Median, 95Percentile, 99Percentile, min, max\n")
out.write("#Arrival on-time stats (predicted time >= arrival time stats\n")
out.write("#Arrival late stats (predicted time < arrival time stats\n")
out.write("#On-time count, expire count, expire rate (%)\n")
stat_l = [c.Stat.MEAN, c.Stat.ERR95, c.Stat.MEDIAN, c.Stat.P95th, c.Stat.P99th, c.Stat.MIN, c.Stat.MAX]
ms_dv = c.get_dv('ms')
#Output on-time stats
for s in stat_l:
  out.write(str(round(ontime_stat[s] / ms_dv, 2)) + ' ')
out.write('\n')
#Output expire stats
for s in stat_l:
  out.write(str(round(expire_stat[s] / ms_dv, 2)) + ' ')
out.write('\n')
#Output expire rate
out.write(str(ontime_n) + ' ' + str(expire_n) + ' ' +  str(expire_rate) + '\n')
out.close()
#print (ontime_n, expire_n, expire_rate)
#print (ontime_stat[c.Stat.MEAN] / ms_dv, ontime_stat[c.Stat.MEDIAN] / ms_dv, ontime_stat[c.Stat.P95th] / ms_dv)
#print (expire_stat[c.Stat.MEAN] / ms_dv, expire_stat[c.Stat.MEDIAN] / ms_dv, expire_stat[c.Stat.P95th] / ms_dv)
