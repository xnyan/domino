#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

exists=`az group exists --name ${resource_group_name}`
if [ "$exists" == "false" ]; then
  log "Error: resource group ${resource_group_name} does not exist. Create it first."; exit 1
fi

log "Creating virtual networks (VNets) in resource group ${resource_group_name}"
dc_vnet_id_list=()
for c in "${cluster_config[@]}"
do
  log "Configuring: $c"
  c=($c)
  dc_location="${c[0]}"
  dc_vnet_name=`gen_vnet_name ${c[1]}`
  dc_vnet_ip=${c[2]}

  cmd="az network vnet create \
    --name ${dc_vnet_name} \
    --resource-group ${resource_group_name} \
    --address-prefixes ${dc_vnet_ip}/${vnet_mask} \
    --subnet-name ${vnet_subnet_name} \
    --subnet-prefix ${dc_vnet_ip}/${vnet_subnet_mask} \
    --location ${dc_location}"

  log "Executing: $cmd"
  run_cmd $cmd

  # Get the id for myVirtualNetwork1.
  cmd="az network vnet show \
    --resource-group ${resource_group_name} \
    --name ${dc_vnet_name} \
    --query id --out tsv"

  log "Executing: $cmd"
  dc_vnet_id=`run_cmd $cmd`
  log "VNet $dc_vnet_name ID: $dc_vnet_id"
  dc_vnet_id_list=("${dc_vnet_id_list[@]}" "$dc_vnet_id")
done

echo ""

log "Peering virtual networks in resource group ${resource_group_name}"
n=${#cluster_config[@]}
for i in `seq 0 $(($n-1))`
do
  src_c=(${cluster_config[$i]})
  src_vnet_name=`gen_vnet_name ${src_c[1]}`
  for j in `seq 0 $(($n-1))`
  do
    if [ "$i" == "$j" ]; then
      continue
    fi
    dst_vnet_id=${dc_vnet_id_list[$j]}
    dst_c=(${cluster_config[$j]})
    dst_vnet_name=`gen_vnet_name ${dst_c[1]}`
    peer_name=`gen_vnet_peer_name ${src_vnet_name} ${dst_vnet_name}`

    cmd="az network vnet peering create \
      --name ${peer_name} \
      --resource-group ${resource_group_name} \
      --vnet-name ${src_vnet_name} \
      --remote-vnet ${dst_vnet_id} \
      --allow-vnet-access"
    
    log "Executing: $cmd"
    run_cmd $cmd
  done
done

echo ""

log "Checking virtual network peering states in resource group $resource_group_name"
for i in `seq 0 $(($n-1))`
do
  src_c=(${cluster_config[$i]})
  src_vnet_name=`gen_vnet_name ${src_c[1]}`
  for j in `seq 0 $(($n-1))`
  do
    if [ "$i" == "$j" ]; then
      continue
    fi
    dst_vnet_id=${dc_vnet_id_list[$j]}
    dst_c=(${cluster_config[$j]})
    dst_vnet_name=`gen_vnet_name ${dst_c[1]}`
    peer_name=`gen_vnet_peer_name ${src_vnet_name} ${dst_vnet_name}`

    cmd="az network vnet peering show \
      --name ${peer_name} \
      --resource-group ${resource_group_name} \
      --vnet-name ${src_vnet_name} \
      --query peeringState"

    log "Checking: ${peer_name}"
    log "Executing: $cmd"
    peer_ret=`run_cmd $cmd`
    log "Peering ${peer_name} state: ${peer_ret}"
    if [ "${peer_ret}" != "\"Connected\"" ]; then
      log "Error: peering ${peer_name} failed. Stop scripts."; exit 1
    fi
  done
done

exit 1

#vnet_file=$1
#
#parse_list_file $vnet_file
#for c in "${__PARSE__RET__[@]}"
#do
#  c=($c)
#  dc=${c[0]}
#  log $dc
#done


