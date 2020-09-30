#!/usr/bin/python
import azure_fp
import common as c

# Combined Basic Fast Paxos and Paxos with different client (locations) in one CDF figure
exp_dir_map = {
  'fp-1c' : '../exp-data/azure-commit-lat-fp-na-4dc-3r-1c/fastpaxos', \
  'p-1c'  : '../exp-data/azure-commit-lat-fp-na-4dc-3r-1c/paxos/', \
  'fp-2c' : '../exp-data/azure-commit-lat-fp-na-4dc-3r-2c/fastpaxos', \
  'p-2c'  : '../exp-data/azure-commit-lat-fp-na-4dc-3r-2c/paxos/', \
  }
exp_list = ['fp-1c', 'fp-2c', 'p-1c', 'p-2c']
exp_label = {
  'fp-1c' : ['-' , 'Fast Paxos 1 client'   , 'grey'],\
  'fp-2c' : ['--', 'Fast Paxos 2 clients'  , 'grey'],\
  'p-1c'  : [':' , 'Multi-Paxos 1 client'  , 'darkblue'],\
  'p-2c'  : ['-.', 'Multi-Paxos 2 clients' , 'darkblue'],\
  }
output_file="azure-commit-lat-fastpaxos.pdf"
azure_fp.custom_lat_cdf(output_file, exp_dir_map, exp_list, exp_label, xmax=152, exp_n=10)
