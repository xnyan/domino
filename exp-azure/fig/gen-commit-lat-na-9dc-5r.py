#!/usr/bin/python
import azure_cdf
import common as c

## Commilt Latency on Azure
p_list = ['dt', 'm', 'ef', 'p']

## NA 9dc 5r Exp ##
exp_dir_map = {
    'dt' : '../exp-data/azure-commit-lat-na-9dc-5r/dynamic/', \
    'm'  : '../exp-data/azure-commit-lat-na-9dc-5r/mencius/', \
    'ef' : '../exp-data/azure-commit-lat-na-9dc-5r/epaxos-thrifty/', \
    'p'  : '../exp-data/azure-commit-lat-na-9dc-5r/paxos/', \
  }
output_file="azure-commit-lat-na-9dc-5r.pdf"
azure_cdf.lat_cdf(output_file, exp_dir_map, p_list, xmax=230, exp_n=10)
