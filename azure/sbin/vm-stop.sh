#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

exists=`az group exists --name ${resource_group_name}`
if [ "$exists" == "false" ]; then
  log "Error: resource group ${resource_group_name} does not exist. Check the group name."; exit 1
fi


log "Deallocating VMs in resource group ${resource_group_name}"
for c in "${cluster_config[@]}"
do
  log "Configuring: $c"
  c=($c)
  dc_location="${c[0]}"
  dc_vnet_name=`gen_vnet_name ${c[1]}`
  vm_num=${c[3]}

  for i in `seq $vm_num`
  do
    vm_name=`gen_vm_name $i ${c[1]}`
    log "Deallocating VM $vm_name at location ${dc_location} in resource group ${resource_group_name}"

    cmd="az vm deallocate \
      --name ${vm_name} \
      --resource-group ${resource_group_name} \
      --no-wait \
      "

    log "Executing: $cmd"
    run_cmd $cmd
  done
done
