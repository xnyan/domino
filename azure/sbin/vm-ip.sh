#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

exists=`az group exists --name ${resource_group_name}`
if [ "$exists" == "false" ]; then
  log "Error: resource group ${resource_group_name} does not exist."; exit 1
fi
    
cmd="az vm list-ip-addresses \
  --resource-group ${resource_group_name} \
  --output table"

log "Executing: $cmd"
run_cmd $cmd

#for c in "${cluster_config[@]}"
#do
#  log "Configuring: $c"
#  c=($c)
#  vm_num=${c[3]}
#
#  for i in `seq $vm_num`
#  do
#    vm_name=`gen_vm_name $i ${c[1]}`
#
#    cmd="az vm list-ip-addresses \
#      --resource-group ${resource_group_name} \
#      --name ${vm_name} \
#      --output table"
#
#    log "Executing: $cmd"
#    run_cmd $cmd
#  done
#done
