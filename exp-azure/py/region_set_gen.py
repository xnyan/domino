import sys
import itertools
import random
from pprint import pprint
# python3 py/region_set_gen.py 3(rep_num) 200(sample_num) japaneast,southeastasia `cat vm-list.config | cut -d " " -f 1`

args = sys.argv
n = int(args[1])
sample_num = int(args[2])
clients = args[3]
regions = args[4:]
sets = []
for s in itertools.combinations(regions, n):
    sets.append(list(s))
replica_sets = random.sample(sets, sample_num)

id=1
for sets in replica_sets:
    print(id, end=" ")
    id += 1
    print(clients, end=" ")
    for s in random.sample(sets, len(sets)):
        print(s, end=",")
    print()