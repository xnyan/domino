#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1


#nsg_names=$(az network nsg list --resource-group ${resource_group_name} --query "[].name" --output tsv)
#
#for nsg_name in $nsg_names; do
#    echo "Adding rule to NSG: $nsg_name"
#    az network nsg rule create --resource-group ${resource_group_name} --nsg-name $nsg_name --name Allow_ICMP --protocol ICMP --direction Inbound --priority 100 --source-address-prefix '*' --destination-address-prefix '*' --access Allow
#done

nsg_names=$(az network nsg list --resource-group "${resource_group_name}" --query "[].name" --output tsv)

for nsg_name in $nsg_names; do
    echo "Updating rule in NSG: $nsg_name"
    az network nsg rule update --resource-group "${resource_group_name}" --nsg-name "$nsg_name" --name Allow_ICMP --access Allow --protocol ICMP --source-address-prefix '*' --destination-address-prefix '*' --priority 100
done
