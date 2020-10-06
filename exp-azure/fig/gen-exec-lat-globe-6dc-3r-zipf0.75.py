#!/usr/bin/python
import azure_exec
import common as c

p_list = ['dt', 'm', 'ef', 'p']

## Exec latency zipf 0.75 (dynamic-95th-8ms)
# Text : (x, y)
txt_map = {
    '(1)' : (260, 0.2), \
    '(2)' : (170, 0.54), \
    '(3)' : (230, 0.88), \
    }
exp_dir_map = {
    'dt':('Domino-8ms' , '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic/'),\
    'm' :('Mencius'    , '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/mencius/'), \
    'ef':('EPaxos'     , '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/epaxos-thrifty/'), \
    'p' :('Multi-Paxos', '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/paxos/'), \
  }
output_file="azure-exec-lat-globe-6dc-3r-zipf0.75.pdf"
azure_exec.lat_cdf(output_file, exp_dir_map, p_list, exp_n=10, txt_map=txt_map)
