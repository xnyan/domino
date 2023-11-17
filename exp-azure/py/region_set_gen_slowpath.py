import sys
import itertools
import random
from pprint import pprint

#python3 region_set_gen_slowpath.py 3 10 `cat ../vm-list.config | cut -d " " -f 1`

args = sys.argv
n = int(args[1])
sample_num = int(args[2])
regions = args[3:]

sets = []
for s in itertools.combinations(regions, n):
    sets.append(list(s))
replica_sets = random.sample(sets, sample_num)

id=1
for sets in replica_sets:
    print(id, end=" ")
    cnt = 0
    id += 1
    l = []
    for s in random.sample(sets, len(sets)):
        l.append(s)
        
    cnt = 0 
    for s in l:
        if cnt == len(sets)-1:
            print(s, end=" ")
        else :
            print(s, end=",")
        cnt += 1
        
    for s in l: 
        print(s, end=",")
    print()
    