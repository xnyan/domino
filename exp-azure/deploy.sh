files=$@
if [ -z "$files" ]; then
  echo "usage: <files>"; exit 1
fi 

source ./sbin/common.sh
load_settings settings.sh
dir="${remote_exec_path}"

log "==== Deploying files ===="
parse_config_file ${remote_client_location_file}
for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"
  
  #log "Creating deployment directory at $dc_id $ip:$dir"
  #cmd="sudo mkdir -p $dir; sudo chmod 777 $dir"
  #cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  #log "Executing: $cmd"
  #run_cmd $cmd

  log "Deploying $files to $dc_id $ip:$dir"
  cmd="scp $SSH_OPTIONS $files ${USER_AT}$ip:$dir/"
  log "Executing: $cmd"
  run_cmd $cmd
  
  log "Checking deployment files at $dc_id $ip:$dir"
  cmd="cd $dir; ls -lh"
  cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  log "Executing: $cmd"
  run_cmd $cmd
  
  sleep 0.01
done

parse_config_file ${remote_server_location_file}
for machine_config in "${COMMON__CONFIG_LIST[@]}"
do
  machine_config=($machine_config)
  dc_id=${machine_config[1]}
  ip="${machine_config[2]}"
  
  #log "Creating deployment directory at $dc_id $ip:$dir"
  #cmd="sudo mkdir -p $dir; sudo chmod 777 $dir"
  #cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  #log "Executing: $cmd"
  #run_cmd $cmd

  log "Deploying $files to $dc_id $ip:$dir"
  cmd="scp $SSH_OPTIONS $files ${USER_AT}$ip:$dir/"
  log "Executing: $cmd"
  run_cmd $cmd
  
  log "Checking deployment files at $dc_id $ip:$dir"
  cmd="cd $dir; ls -lh"
  cmd="ssh -n $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  log "Executing: $cmd"
  run_cmd $cmd
  
  sleep 0.01
done
wait
