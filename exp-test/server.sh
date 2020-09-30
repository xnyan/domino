#!/usr/bin/env bash

mode=$1
app=$2

if [ -z $mode ] || [ -z $app ]; then
  echo "Usage: <start | stop | list> <dynamic | fastpaxos | epaxos | mencius | paxos>"
  echo "dynamic is an alias name for Domino, and paxos is an alias name for Multi-Paxos"
  echo "For example, ./server.sh start dynamic"
  exit 1
fi
    
if [ "${app}" == "dynamic" ] || [ "${app}" == "domino" ]; then 
  # Domino (dynamic is an alias name for Domino)
  binary="dynamic"
elif [ "${app}" == "fastpaxos" ]; then 
  # Fast Paxos
  binary="fastpaxos"
elif [ "${app}" == "epaxos" ] || [ "${app}" == "mencius" ] || [ "${app}" == "paxos" ] || [ "${app}" == "multipaxos" ]; then
  # EPaxos, Mencius, and Multi-Paxos (paxos is an alias name for Multi-Paxos)
  binary="epaxos"
else
  echo "Invalid protocol type = ${app}"
fi

if [ "$mode" == "start" ]; then 
  for id in {1..3}
  do
    if [ "${app}" == "dynamic" ] || [ "${app}" == "domino" ]; then 
      # Domino
      echo "Starting Domino replica server $id"
      ./dynamic -i $id -c test.config -r replica-location.config > server-$id.log 2>&1 &
    elif [ "${app}" == "fastpaxos" ]; then 
      # Fast Paxos
      echo "Starting Fast Paxos replica server $id"
      ./fastpaxos -i ${id} -c test.config -r replica-location.config > server-$id.log 2>&1 &
    elif [ "${app}" == "epaxos" ]; then
      # EPaxos
      echo "Starting EPaxos replica server $id"
      ./epaxos -i ${id} -c test.config -r replica-location.config -p e -t true > server-$id.log 2>&1 &
    elif [ "${app}" == "mencius" ]; then
      # Mencius
      echo "Starting Mencius replica server $id"
      ./epaxos -i ${id} -c test.config -r replica-location.config -p m -m false > server-$id.log 2>&1 &
    elif [ "${app}" == "paxos" ] || [ "${app}" == "multipaxos" ]; then
      # Multi-Paxos
      echo "Starting Multi-Paxos replica server $id"
      ./epaxos -i ${id} -c test.config -r replica-location.config -p p -t false > server-$id.log 2>&1 &
    fi #app
    sleep 0.1
  done
elif [ "$mode" == "stop" ]; then 
  killall ${binary}; echo 'Stopped replica servers'
elif [ "$mode" == "list" ]; then
  ps -ef | grep ${binary}
else
  echo "Invalid mode = ${mode}"
fi #mode

sleep 2
echo "$mode servers done"
