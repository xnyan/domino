#!/usr/bin/python
import azure_cdf
import common as c

## Commilt Latency on Azure
p_list = ['dt', 'm', 'ef', 'p']

## NA 9dc 3r Exp ##
exp_dir_map = {
    'dt' : '../exp-data/azure-commit-lat-na-9dc-3r/dynamic/', \
    'm'  : '../exp-data/azure-commit-lat-na-9dc-3r/mencius/', \
    'ef' : '../exp-data/azure-commit-lat-na-9dc-3r/epaxos-thrifty/', \
    'p'  : '../exp-data/azure-commit-lat-na-9dc-3r/paxos/', \
  }
output_file="azure-commit-lat-na-9dc-3r.pdf"
azure_cdf.lat_cdf(output_file, exp_dir_map, p_list, xmax=230, exp_n=10)
