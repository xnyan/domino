#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

source $sbin/common.sh

setting_file=$1
# start, stop, or list
mode=$2
# d (Dynamic), e (EPaxos), m (Mencius), or p (Multi-Paxos)
protocol_type=$3

load_settings ${setting_file}

if [ "$is_debug_mode" == "true" ]; then
  DEBUG="-d"
fi

if [ ! -z ${protocol_type} ]; then
  CMD_OPTIONS="-p ${protocol_type}"
fi

#Loads the locations of clients 
parse_config_file ${remote_client_location_file}

for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  m_id=${machine_config[0]}
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"
  i_num="${machine_config[3]}"
  target_dc="${machine_config[4]}"

  if [ "$mode" == "start" ]; then 
    cmd="cd ${remote_exec_path};" 
    cN=`seq $i_num` 
    for i in $cN
    do
      cId="${m_id}-$i"
      client_log_file="${remote_client_log_dir}/client-${cId}-${dc_id}-${target_dc}.log"
      cmd="$cmd ./${client_app} -i ${cId} -dc ${dc_id} -c ${config_file} -r ${replica_location_file} -t ${target_dc} ${CMD_OPTIONS} $DEBUG > ${client_log_file} 2>&1 &"
    done
    log "Starting ${i_num} ${client_app} at ${dc_id} $ip"
  elif [ "$mode" == "stop" ]; then
    cmd="killall ${client_app}; echo 'Stopped ${client_app} at ${dc_id} $ip'" 
  elif [ "$mode" == "list" ]; then
    cmd="ps -ef | grep ${client_app}"
  fi #mode

  cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  
  log "$mode ${client_app} at $dc_id $ip"
  log "Executing command: $cmd"
  run_cmd $cmd
  run_cmd_in_background $cmd
  sleep 0.01
done
wait
log "Clients $mode done."
