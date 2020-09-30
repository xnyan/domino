setting_file="settings.sh"
source ${setting_file}

./sbin/build.sh ${setting_file}
vm_ip_config='azure-exp-vm-ip.config'
dc_delay_file='azure-globe-delay.json'
key_file='key.dat'
leader_dc='westus2'
./sbin/gen-location.py -f ${vm_ip_config} -l ${leader_dc} -d ${dc_delay_file}
./deploy.sh ${dynamic_app} ${epaxos_app} ${client_app} ${config_file} ${replica_location_file} ${key_file}

# Experiment result dir
exp_data_dir="exp-data/azure-commit-lat-globe-6dc-3r"
wait_time=96

# Protocol
base_p="dynamic"
base_config_file="config/default.config.commit"

run_exp() {
    # server
    ./sbin/server.sh ${setting_file} start dynamic
    sleep 5
    # client
    ./sbin/client.sh ${setting_file} start d
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
    ./sbin/server.sh ${setting_file} stop dynamic

    # delete logs
    ./sbin/log.sh ${setting_file} delete
}

# Vary the percentile delay from network measurements for delay prediction
# No need to run the 95th percentile delay again since it is the default value
# when comparing Domino with other protocols
pth="0.95"
for pth in 0.99 0.9 0.75 0.5
do
  sed "s/dynamic.lat.predict.percentile = 0.95/dynamic.lat.predict.percentile = ${pth}/g" ${base_config_file} > ${config_file}
  ./deploy.sh ${config_file}
  p="${base_p}-pth${pth}"
  for i in {1..10}
  do
    echo "`date` Experiment: ${base_p} ${base_config_file} ${pth} $i"
    run_exp
  done
done #pth
