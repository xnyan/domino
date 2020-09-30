#!/usr/bin/env bash

app=$1

if [ -z $app ]; then
  echo "Usage: <dynamic | fastpaxos | epaxos | mencius | paxos>"
  echo "dynamic is an alias name for Domino, and paxos is an alias name for Multi-Paxos"
  echo "For example, ./client.sh dynamic"
  exit 1
fi

if [ "${app}" == "dynamic" ] || [ "${app}" == "domino" ]; then 
  # Domino
  CMD_OPTIONS="-p d"
elif [ "${app}" == "fastpaxos" ]; then 
  # Fast Paxos
  CMD_OPTIONS="-p fp"
elif [ "${app}" == "epaxos" ]; then
  # EPaxos
  CMD_OPTIONS="-p e"
elif [ "${app}" == "mencius" ]; then
  # Mencius
  CMD_OPTIONS="-p m"
elif [ "${app}" == "paxos" ] || [ "${app}" == "multipaxos" ]; then
  # Multi-Paxos
  CMD_OPTIONS="-p p"
else
  echo "Invalid protocol $app"
fi

./client -i 1 -dc dc1 -c test.config -r replica-location.config -t dc1 ${CMD_OPTIONS}
