#!/usr/bin/env bash

sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_settings $1
# collect, delete, list
mode=$2

mkdir -p ${local_log_dir}

## Server logs
#Loads the locations of servers
parse_config_file ${remote_server_location_file}
for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  id=${machine_config[0]}
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"

  server_log_file="${remote_server_log_dir}/server*.log"
  if [ "$mode" == "collect" ]; then
    cmd="scp $SSH_OPTIONS ${USER_AT}$ip:${server_log_file} ${local_log_dir}/"
    log "Copying server log files ${server_log_file} from $dc_id ($ip) to ${local_log_dir}/"
  elif [ "$mode" == "delete" ]; then
    cmd="cd ${remote_server_log_dir}; rm server*.log kv-*.log stable*"
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
    log "Deleting server log files on $dc_id ($ip)"
  elif [ "$mode" == "list" ]; then
    cmd="ls -lh ${server_log_file}"
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
    log "Listing server log files ${server_log_file} on $dc_id ($ip)"
  fi
  log "Executing command: $cmd"
  
  run_cmd $cmd
  #run_cmd_in_background $cmd
  sleep 0.01
done
wait
log "Server log files are stored in ${local_log_dir}/"

## Client logs
#Loads the locations of clients
parse_config_file ${remote_client_location_file}
for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  id=${machine_config[0]}
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"

  client_log_file="${remote_client_log_dir}/client*.log"
  if [ "$mode" == "collect" ]; then
    cmd="scp $SSH_OPTIONS ${USER_AT}$ip:${client_log_file} ${local_log_dir}/"
    log "Copying client log files ${client_log_file} from $dc_id ($ip) to ${local_log_dir}/"
  elif [ "$mode" == "delete" ]; then
    cmd="cd ${remote_client_log_dir}; rm client*.log"
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
    log "Deleting client log files on $dc_id ($ip)"
  elif [ "$mode" == "list" ]; then
    cmd="ls -lh ${client_log_file}"
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
    log "Listing client log files ${client_log_file} on $dc_id ($ip)"
  fi
  log "Executing command: $cmd"
  
  run_cmd $cmd
  #run_cmd_in_background $cmd
  sleep 0.01
done
wait
log "Client log files are stored in ${local_log_dir}/"
