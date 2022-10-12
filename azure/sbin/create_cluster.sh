#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

load_setting $1

# Creates a resource group
$sbin/group.sh $1

# Creates virtual networks
$sbin/vnet.sh $1

# Creates VMs
$sbin/vm.sh $1

# Creates Rule
$sbin/nsg.sh $1

# Lists VM IPs
$sbin/vm-ip.sh $1
