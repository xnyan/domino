#!/usr/bin/env bash

sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

source $sbin/common.sh

usage="Usage: <-d | -ds (download parallel | sequentially) | -r (remove) | -l (list)> <setting_file>"
if [ -z $1 ] || [ -z $2 ]; then
  echo  $usage; exit -1
fi

mode=$1
load_settings $2

mkdir -p $local_log_dir

#Loads the locations of servers and partitions
parse_config_file ${remote_location_file}
for loc in "${COMMON__CONFIG_LIST[@]}"
do
  loc=($loc)
  ip_port=(`echo ${loc[1]} | tr ":" " "`)
  ip=${ip_port[0]}

  if [ "$mode" == "-d" ] || [ "$mode" == "-ds" ]; then
    cmd="scp $SSH_OPTIONS ${USER_AT}$ip:${remote_log_dir}/*.log ${local_log_dir}/"
    echo "Downloading *.log files from $ip:${remote_log_dir} to ${local_log_dir}/"
  elif [ "$mode" == "-r" ]; then
    cmd="cd ${remote_log_dir}; rm *.log"
    cmd="ssh $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
    echo "Deleting *.log files at $ip:${remote_log_dir}"
  elif [ "$mode" == "-l" ]; then
    cmd="cd ${remote_log_dir}; ls -lh *.log"
    cmd="ssh $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
    echo "Listing *.log files at $ip:${remote_log_dir}"
  else
    echo  $usage; exit -1
  fi
  
  echo "Executing command: $cmd"
  if [ "$mode" == "-l" ] || [ "$mode" == "-ds" ]; then
    run_cmd $cmd
  else
    run_cmd_in_background $cmd
  fi
  sleep 0.01
done

wait
date
