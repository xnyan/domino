#!/usr/bin/python
import azure_exec
import common as c

p_list = ['dt', 'm', 'ef', 'p']

## Exec latency zipf 0.95 (dynamic-95th-8ms)
# Text : (x, y)
txt_map = {
    '(4)' : (200, 0.42), \
    '(5)' : (460, 0.7), \
    }
exp_dir_map = {
    'dt':('Domino-8ms' , '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.95/dynamic/'),\
    'm' :('Mencius'    , '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.95/mencius/'), \
    'ef':('EPaxos'     , '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.95/epaxos-thrifty/'), \
    'p' :('Multi-Paxos', '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.95/paxos/'), \
  }
output_file="azure-exec-lat-globe-6dc-3r-zipf0.95.pdf"
azure_exec.lat_cdf(output_file, exp_dir_map, p_list, xmax=910, exp_n=10, txt_map=txt_map)
