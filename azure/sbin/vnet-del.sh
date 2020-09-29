#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

for c in "${cluster_config[@]}"
do
  c=($c)
  dc_vnet_name=`gen_vnet_name ${c[1]}`
  cmd="az network vnet delete -g ${resource_group_name} -n ${dc_vnet_name}"

  log "Deleting VNet $dc_vnet_name in resource group ${resource_group_name}"
  log "Executing: $cmd"
  run_cmd $cmd
done

log "Done deleting VNets in resource group ${resource_group_name}"
