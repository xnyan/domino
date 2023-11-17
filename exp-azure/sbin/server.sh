#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

source $sbin/common.sh

setting_file=$1
# start, stop, or list
mode=$2
# dynamic, epaxos, or fastpaxos
server_app=$3
# only for epaxos: e (EPaxos), m (Mencius), or p (Multi-Paxos)
epaxos_type=$4
# true, false or empty. For Mencius: early commit; For EPaxos and Multi-Paxos: thrifty
epaxos_opt=$5

load_settings ${setting_file}

if [ "$is_debug_mode" == "true" ]; then
  DEBUG="-d"
fi

if [ "$epaxos_type" == "m" ] && [ ! -z ${epaxos_opt} ]; then
  CMD_OPTIONS="-m ${epaxos_opt}"
elif [ "$epaxos_type" == "e" ] && [ ! -z ${epaxos_opt} ]; then
  CMD_OPTIONS="-t ${epaxos_opt}"
elif [ "$epaxos_type" == "p" ] && [ ! -z ${epaxos_opt} ]; then
  CMD_OPTIONS="-t ${epaxos_opt}"
fi

#Loads the locations of servers and partitions
parse_config_file ${remote_server_location_file}

for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  server_id=${machine_config[0]}
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"
  server_log_file="${remote_server_log_dir}/server-${server_id}-${dc_id}.log"

  if [ "$mode" == "start" ]; then 
    if [ "${server_app}" == "${dynamic_app}" ]; then 
      cmd="cd ${remote_exec_path}; ./${server_app} -i ${server_id} -c ${config_file} -r ${replica_location_file} $DEBUG > ${server_log_file} 2>&1 &" 
    elif [ "${server_app}" == "${fastpaxos_app}" ]; then 
      cmd="cd ${remote_exec_path}; ./${server_app} -i ${server_id} -c ${config_file} -r ${replica_location_file} $DEBUG > ${server_log_file} 2>&1 &" 
    elif [ "${server_app}" == "${epaxos_app}" ]; then
      cmd="cd ${remote_exec_path}; ./${server_app} -i ${server_id} -c ${config_file} -r ${replica_location_file} -p ${epaxos_type} ${CMD_OPTIONS} $DEBUG > ${server_log_file} 2>&1 &" 
    else
      log "Invalid server app = ${server_app}"
    fi #server_app
  elif [ "$mode" == "stop" ]; then 
    cmd="pkill ${server_app}; echo 'Stopped ${server_app} at $dc_id $ip'" 
  elif [ "$mode" == "list" ]; then
    cmd="ps -ef | grep ${server_app}" 
  fi #mode

  cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""

  log "$mode ${server_app} at $dc_id $ip"
  log "Executing command: $cmd"
  
  run_cmd $cmd
  run_cmd_in_background $cmd
  sleep 0.01
done
wait
log "Servers $mode done."
