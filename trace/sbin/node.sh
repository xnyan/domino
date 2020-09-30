#!/usr/bin/env bash

sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

source $sbin/common.sh
source $sbin/cmd.sh

usage="Usage: <s (server) | c (client)> \
  <-s (start) | -t (stop) | -c (check process)> \
  <setting_file>"
if [ -z $1 ] || [ -z $2 ] || [ -z $3 ]; then
  echo  $usage; exit -1
fi

node_type=$1
mode=$2
load_settings $3

if [ "$is_debug_mode" == "true" ]; then
  DEBUG="-d"
fi
dir="${remote_exec_path}"

#Loads the locations of servers and partitions
parse_config_file ${remote_location_file}
for loc in "${COMMON__CONFIG_LIST[@]}"
do
  loc=($loc)
  dc_id=${loc[0]}
  ip_port=(`echo ${loc[1]} | tr ":" " "`)
  ip=${ip_port[0]}
  server_log_file="${remote_log_dir}/server-${dc_id}.log"
  client_log_file="${remote_log_dir}/client-${dc_id}.log"

  #Generates cmd
  generate_cmd $node_type $mode

  cmd="ssh $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""

  echo "Executing command: $cmd"
  if [ "$mode" == "-c" ]; then
    run_cmd $cmd
  else
    run_cmd_in_background $cmd
  fi
  sleep 0.01
done
  
wait
date
