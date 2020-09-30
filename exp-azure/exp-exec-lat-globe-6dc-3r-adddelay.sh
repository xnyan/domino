setting_file="settings.sh"
source ${setting_file}

./sbin/build.sh ${setting_file}
vm_ip_config='azure-exp-vm-ip.config'
dc_delay_file='azure-globe-delay.json'
key_file='key.dat'
leader_dc='westus2'
./sbin/gen-location.py -f ${vm_ip_config} -l ${leader_dc} -d ${dc_delay_file}
./deploy.sh ${dynamic_app} ${epaxos_app} ${client_app} ${config_file} ${replica_location_file} ${key_file}

# Protocol
base_p="dynamic"
base_config_file="config/default.config.exec"

zipf="0.75"
exp_data_dir="exp-data/azure-exec-lat-globe-6dc-3r-zipf${zipf}"
wait_time=96

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

# Vary the additional delay for increasing DFP request timestamps
# No need to run the 8ms addition delay since it is the default value when
# comparing Domino with other protocols
add_delay="8ms"
for add_delay in 0ms 1ms 2ms 4ms 12ms 16ms 24ms 36ms
do
  sed "s/client.add.delay = 8ms/client.add.delay = ${add_delay}/g" ${base_config_file} > ${config_file}
  ./deploy.sh ${config_file}

  p="${base_p}-add${add_delay}"
  for i in {1..10}
  do
    echo "`date` Experiment: ${base_p} ${base_config_file} ${add_delay} $i"
    run_exp
  done
done #add_delay
