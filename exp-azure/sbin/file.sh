#!/usr/bin/env bash

sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_settings $1
# clean, list
mode=$2

clean_file_list="${dynamic_app} ${epaxos_app} ${fastpaxos_app} ${client_app} ${config_file} ${replica_location_file} key.dat"

## Servers
#Loads the locations of servers
parse_config_file ${remote_server_location_file}
for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  id=${machine_config[0]}
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"

  if [ "$mode" == "clean" ]; then
    op="cd ${remote_exec_path}; rm ${clean_file_list}"
    op="${op}; cd ${remote_server_log_dir}; rm server*.log kv-*.log stable*"
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$op\""
    log "Deleting files on server $id at $dc_id ($ip)"
  elif [ "$mode" == "list" ]; then
    if [ "${remote_exec_path}" == "${remote_server_log_dir}" ]; then
      op="ls -lh ${remote_exec_path}"
    else
      op="ls -lh ${remote_exec_path} ${remote_server_log_dir}"
    fi
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$op\""
    log "Listing files on server $id at $dc_id ($ip)"
  fi
  log "Executing command: $cmd"
  
  run_cmd $cmd
  #run_cmd_in_background $cmd
  sleep 0.01
done
wait

## Clients
#Loads the locations of clients
parse_config_file ${remote_client_location_file}
for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  id=${machine_config[0]}
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"

  if [ "$mode" == "clean" ]; then
    op="cd ${remote_exec_path}; rm ${clean_file_list}"
    op="${op}; cd ${remote_client_log_dir}; rm client*.log"
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$op\""
    log "Deleting files on client $id at $dc_id ($ip)"
  elif [ "$mode" == "list" ]; then
    if [ "${remote_exec_path}" == "${remote_client_log_dir}" ]; then
      op="ls -lh ${remote_exec_path}"
    else
      op="ls -lh ${remote_exec_path} ${remote_client_log_dir}"
    fi
    cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$op\""
    log "Listing files on client $id at $dc_id ($ip)"
  fi
  log "Executing command: $cmd"
  
  run_cmd $cmd
  #run_cmd_in_background $cmd
  sleep 0.01
done
wait
