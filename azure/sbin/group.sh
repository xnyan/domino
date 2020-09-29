#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

exists=`az group exists --name $resource_group_name`
if [ "$exists" == "true" ]; then
  echo "Error: resource group $resource_group_name already exists. Change to a different one."; exit 1
fi

cmd="az group create --name $resource_group_name --location $resource_group_location"

log "Creating resource group ${resource_group_name} at location $resource_group_location"
log "Executing: $cmd"
run_cmd $cmd

exit 1

## Deletes a resource group
cmd="az group delete --name $resource_group_name"
#--no-wait
#--yes

## Other commands for resource groups
#cmd="az group list [--name]"
#cmd="az group export --name $resource_group_name"

