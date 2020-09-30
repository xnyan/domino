setting_file="settings.sh"
source ${setting_file}

# Experiment configuration file
cp config/default.config.exec ${config_file}

./sbin/build.sh ${setting_file}
vm_ip_config='azure-exp-vm-ip.config'
dc_delay_file='azure-globe-delay.json'
key_file='key.dat'
leader_dc='westus2'
./sbin/gen-location.py -f ${vm_ip_config} -l ${leader_dc} -d ${dc_delay_file}
./deploy.sh ${dynamic_app} ${epaxos_app} ${client_app} ${config_file} ${replica_location_file} ${key_file}

base_config_file="config/default.config.exec"
for zipf in 0.75 0.95
do
sed "s/workload.zipf.alpha = 0.75/workload.zipf.alpha = ${zipf}/g" ${base_config_file} > ${config_file}
./deploy.sh ${config_file}

exp_data_dir="exp-data/azure-exec-lat-globe-6dc-3r-zipf${zipf}"
wait_time=96

for i in {1..10}
do
  ## Protocols:
  # paxos: Multi-Paxos
  # fastpaxos: Fast Paxos
  # dynamic: Domino
  # mencius: Mencius
  # epaxos-thrifty: EPaxos with its thrifty optimization
  # epaxos: EPaxos without its thrifty optimization
  for p in paxos dynamic mencius epaxos-thrifty
  do
    echo "`date` Experiment: exec $zipf $i $p"
    # server
    if [ "$p" == "dynamic" ]; then
      ./sbin/server.sh ${setting_file} start dynamic
    elif [ "$p" == "epaxos" ]; then
      ./sbin/server.sh ${setting_file} start epaxos e false
    elif [ "$p" == "epaxos-thrifty" ]; then
      ./sbin/server.sh ${setting_file} start epaxos e true
    elif [ "$p" == "mencius" ]; then
      ./sbin/server.sh ${setting_file} start epaxos m false 
    elif [ "$p" == "paxos" ]; then
      ./sbin/server.sh ${setting_file} start epaxos p false
    fi
    sleep 5

    #client
    if [ "$p" == "dynamic" ]; then
      ./sbin/client.sh ${setting_file} start d
    elif [ "$p" == "epaxos" ] || [ "$p" == "epaxos-thrifty" ]; then
      ./sbin/client.sh ${setting_file} start e
    elif [ "$p" == "mencius" ]; then
      ./sbin/client.sh ${setting_file} start m
    elif [ "$p" == "paxos" ]; then
      ./sbin/client.sh ${setting_file} start p
    fi
    sleep ${wait_time}

    # Log files
    ./sbin/log.sh ${setting_file} collect
    ret_dir=${exp_data_dir}/$p/$i
    mkdir -p ${ret_dir}
    mv ${local_log_dir}/*.log ${ret_dir}/ 
    # Copy config files as well
    cp ${config_file} ${remote_client_location_file} ${remote_server_location_file} ${replica_location_file} ${ret_dir}/ 
   
    # stop clients and servers 
    ./sbin/client.sh ${setting_file} stop
    if [ "$p" == "dynamic" ]; then
      ./sbin/server.sh ${setting_file} stop dynamic
    else
      ./sbin/server.sh ${setting_file} stop epaxos
    fi

    # delete logs
    ./sbin/log.sh ${setting_file} delete
  done #p
done #i
done #zipf
