#!/usr/bin/python3

import argparse
import json
import random
import string
import sys

argParser = argparse.ArgumentParser("")

argParser.add_argument(
    '-k', '--key', 
    dest='key', 
    nargs='?', 
    help='total number of keys',
    type=int,
    required=True)

argParser.add_argument(
    '-l', '--length', 
    dest='key_length', 
    nargs='?', 
    help='the length of a key in bytes; default 64',
    type=int,
    default=64)

argParser.add_argument(
    '-s', '--seed', 
    dest='seed', 
    nargs='?', 
    help='random seed; 0 for dynamic; default 1',
    type=int,
    default=1)

args = argParser.parse_args()
if args.seed != 0:
  random.seed(args.seed)

#key generation is cloned from TAPIR
charset = string.ascii_uppercase + string.ascii_lowercase + string.digits
keyDict = {}

for i in range(args.key):
  rkey = "".join(random.choice(charset) for j in range(args.key_length))
  while rkey in keyDict:
    rkey = "".join(random.choice(charset) for j in range(args.key_length))
  keyDict[rkey] = True
  print(rkey)

