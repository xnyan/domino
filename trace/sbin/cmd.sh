generate_cmd() {
  _node_type_=$1
  _mode_=$2
  if [ "$_node_type_" == "s" ] && [ "$_mode_" == "-s" ]; then
    server_start
  elif [ "$_node_type_" == "s" ] && [ "$_mode_" == "-t" ]; then
    server_stop
  elif [ "$_node_type_" == "s" ] && [ "$_mode_" == "-c" ]; then
    check_process ${server_app}
  elif [ "$_node_type_" == "c" ] && [ "$_mode_" == "-s" ]; then
    client_start
  elif [ "$_node_type_" == "c" ] && [ "$_mode_" == "-t" ]; then
    client_stop
  elif [ "$_node_type_" == "c" ] && [ "$_mode_" == "-c" ]; then
    check_process ${client_app}
  else
    echo  $usage; exit -1
  fi  
}

server_start() {
  #cmd="${remote_setup}; cd $dir; ./${server_app} -i ${dc_id} -l ${location_file} $DEBUG > ${server_log_file} 2>&1 &" 
  cmd="cd $dir; ./${server_app} -i ${dc_id} -l ${location_file} $DEBUG > ${server_log_file} 2>&1 &" 
  echo "Starting ${server_app} at" $ip $dc_id
}

server_stop() {
  #cmd="pgrep -f 'server -i'"
  #ssh $ip 'pid=`pgrep ${server_app}`; echo $pid' #TODO ${server_app} should be materialized locally
  #cmd="killall ${server_app}; echo 'Stopped ${server_app} at $ip'" 
  cmd="pkill ${server_app}; echo 'Stopped ${server_app} at $ip'" 
  echo "Stopping ${server_app} at" $ip $dc_id
}

client_start() {
  cmd="cd $dir; ./${client_app} -i ${dc_id} -l ${location_file} -c ${config_file} $DEBUG > ${client_log_file} 2>&1 &"
  echo "Starting ${client_app} at" $ip $dc_id
}

client_stop() {
  #cmd="pkill ${client_app}; echo 'Stopped ${client_app} at $ip'" 
  cmd="killall ${client_app}; echo 'Stopped ${client_app} at $ip'" 
  echo "Stopping ${client_app} at" $ip $dc_id
}

check_process() {
  _app_=$1
  #cmd="pgrep -f '${_app_} -i'"
  cmd="ps -ef | grep ${_app_}" 
  echo "Checking ${_app_} process at" $ip $dc_id
}
