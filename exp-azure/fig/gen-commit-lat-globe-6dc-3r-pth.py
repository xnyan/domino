#!/usr/bin/python
import azure_tail
import common as c
import label
import custom

def plot(output_file, idx, cl_list, y_label):
  exp_list = [
      'dt1s0ms'  , \
      'dt1s1ms'  , \
      'dt1s2ms'  , \
      'dt1s4ms'  , \
      'dt1s8ms'  , \
      'dt1s12ms' , \
      'dt1s16ms' , \
      ]
  net95th_exp_dir_map = {
      'dt1s0ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic', \
      'dt1s1ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.95-add1ms', \
      'dt1s2ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.95-add2ms', \
      'dt1s4ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.95-add4ms', \
      'dt1s8ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.95-add8ms', \
      'dt1s12ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.95-add12ms', \
      'dt1s16ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.95-add16ms', \
      }
  net99th_exp_dir_map = {
      'dt1s0ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add0ms', \
      'dt1s1ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add1ms', \
      'dt1s2ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add2ms', \
      'dt1s4ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add4ms', \
      'dt1s8ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add8ms', \
      'dt1s12ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add12ms', \
      'dt1s16ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.99-add16ms', \
      }
  net90th_exp_dir_map = {
      'dt1s0ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add0ms', \
      'dt1s1ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add1ms', \
      'dt1s2ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add2ms', \
      'dt1s4ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add4ms', \
      'dt1s8ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add8ms', \
      'dt1s12ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add12ms', \
      'dt1s16ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.9-add16ms', \
      }
  net75th_exp_dir_map = {
      'dt1s0ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add0ms', \
      'dt1s1ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add1ms', \
      'dt1s2ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add2ms', \
      'dt1s4ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add4ms', \
      'dt1s8ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add8ms', \
      'dt1s12ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add12ms', \
      'dt1s16ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.75-add16ms', \
      }
  net50th_exp_dir_map = {
      'dt1s0ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add0ms', \
      'dt1s1ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add1ms', \
      'dt1s2ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add2ms', \
      'dt1s4ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add4ms', \
      'dt1s8ms'  : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add8ms', \
      'dt1s12ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add12ms', \
      'dt1s16ms' : '../exp-data/azure-commit-lat-globe-6dc-3r/dynamic-pth0.5-add16ms', \
      }
  line_list=['m', 'ef', 'p']
  line_dir_map = {
      'm'  : '../exp-data/azure-commit-lat-globe-6dc-3r/mencius/', \
      'ef' : '../exp-data/azure-commit-lat-globe-6dc-3r/epaxos-thrifty/', \
      'p'  : '../exp-data/azure-commit-lat-globe-6dc-3r/paxos/', \
    }
  line_map = {
      'm'  : [0, 'Mencius', label.line_color['m'] , label.line_style['m']], \
      'ef' : [0, 'EPaxos' , label.line_color['ef'], label.line_style['ef']], \
      'p'  : [0, 'Multi-Paxos'  , label.line_color['p'] , label.line_style['p']], \
      }
  line_lat, line_err = azure_tail.get_metric(idx, line_dir_map, cl_list, line_list, exp_n=10)
  for i in range(len(line_list)):
    line_map[line_list[i]][0] = line_lat[i]
  net50_lat, net50_err = azure_tail.get_metric(idx, net50th_exp_dir_map, cl_list, exp_list, exp_n=10)
  net75_lat, net75_err = azure_tail.get_metric(idx, net75th_exp_dir_map, cl_list, exp_list, exp_n=10)
  net90_lat, net90_err = azure_tail.get_metric(idx, net90th_exp_dir_map, cl_list, exp_list, exp_n=10)
  net95_lat, net95_err = azure_tail.get_metric(idx, net95th_exp_dir_map, cl_list, exp_list, exp_n=10)
  net99_lat, net99_err = azure_tail.get_metric(idx, net99th_exp_dir_map, cl_list, exp_list, exp_n=10)
  x_tick_l=['0', '1', '2', '4', '8', '12', '16']
  y_list = ['net50', 'net75', 'net90', 'net95', 'net99']
  y_map = {
      'net50' : (net50_lat, net50_err, 'using p50th in network measurements', '\\', 'grey'), \
      'net75' : (net75_lat, net75_err, 'using p75th in network measurements', '/' , 'grey'), \
      'net90' : (net90_lat, net90_err, 'using p90th in network measurements', 'o' , 'white'), \
      'net95' : (net95_lat, net95_err, 'using p95th in network measurements', ''  , 'white'), \
      'net99' : (net99_lat, net99_err, 'using p99th in network measurements', 'x' , 'white'), \
      }
  azure_tail.plot_lat_bar_and_line(output_file, y_label, x_tick_l, y_map, y_list, line_list, line_map)

cl_list = [
    'client-1-1-australiaeast-australiaeast.log', \
    'client-2-1-eastus2-westus2.log', \
    'client-3-1-francecentral-francecentral.log', \
    'client-4-1-westus2-westus2.log', \
    'client-5-1-eastasia-australiaeast.log', \
    'client-6-1-southeastasia-australiaeast.log', \
    ]
y_label='99th Percentile Commit Latency (ms)'
output_file="azure-commit-lat-globe-6dc-3r-pth.pdf"
plot(output_file, custom.idx_lat_99, cl_list, y_label)
