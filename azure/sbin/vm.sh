#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

exists=`az group exists --name ${resource_group_name}`
if [ "$exists" == "false" ]; then
  log "Error: resource group ${resource_group_name} does not exist. Create it first."; exit 1
fi


log "Creating VMs in resource group ${resource_group_name}"
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
    log "Creating VM $vm_name at location ${dc_location} in resource group ${resource_group_name}"
    cmd="az vm create \
      --resource-group ${resource_group_name} \
      --location ${dc_location} \
      --vnet-name ${dc_vnet_name} \
      --subnet ${vnet_subnet_name} \
      --name ${vm_name} \
      --size ${vm_size} \
      --image ${vm_image} \
      --storage-sku ${vm_disk_type} \
      --admin-username ${vm_username} \
      --ssh-key-values ${vm_public_key} \
      "

    ## Optional configurations
    if [ "$vm_acc_network" == "true" ]; then
      cmd="$cmd --accelerated-networking ${vm_acc_network}"
    fi
    if [ "$vm_dns" == "true" ]; then
      vm_dns_name=`gen_dns_name ${vm_name}`
      cmd="$cmd --public-ip-address-dns-name ${vm_dns_name}"
    fi
    if [ "$vm_no_wait" == "true" ]; then
      cmd="$cmd --no-wait"
    fi

    log "Executing: $cmd"
    run_cmd $cmd

    ## Network information can be immediately fetched from the command output by using the following keys
    ##"fqdns": "vm1-eastus-dns.eastus.cloudapp.azure.com",
    ##"macAddress": "00-0D-3A-8E-A5-FB",
    ##"privateIpAddress": "10.0.0.4",
    ##"publicIpAddress": "13.68.138.224",
  done
done
