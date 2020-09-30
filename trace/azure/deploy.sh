source ../sbin/common.sh

load_settings settings.sh

files=$@
if [ -z "$files" ]; then
  echo "usage: <files>"; exit 1
fi 

echo ""
echo "`date` ==== Deploying binaries and configuration files ===="
echo ""

dir="${remote_exec_path}"

#Loads the locations of servers and partitions
parse_config_file ${remote_location_file}
for loc in "${COMMON__CONFIG_LIST[@]}"
do
  loc=($loc)
  dc_id=${loc[0]}
  ip_port=(`echo ${loc[1]} | tr ":" " "`)
  ip=${ip_port[0]}

  #echo "Creating deployment directory: $ip:$dir"
  #cmd="sudo mkdir -p $dir; sudo chmod 777 $dir"
  #cmd="ssh $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  #echo "Executing: $cmd"
  #run_cmd $cmd

  echo "Deploying $files to $ip:$dir"
  cmd="scp $SSH_OPTIONS $files ${USER_AT}$ip:$dir/"
  echo "Executing: $cmd"
  run_cmd $cmd

  echo "Checking deployment @ $ip:$dir"
  cmd="cd $dir; ls"
  cmd="ssh $SSH_OPTIONS ${USER_AT}$ip \"$cmd\""
  echo "Executing: $cmd"
  run_cmd $cmd

  echo ""
  sleep 0.01
done
  
wait
date
