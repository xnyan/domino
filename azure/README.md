# Scripts for creating a private cluster across datacenters on Azure

This folder has a set of scripts to create a cluster of identical virtual
machines (VMs) in different datacenters on Microsoft Azure and to connect
these VMs through Azure's vitual network (VNets) peering.

./sbin: scripts to use Azure CLI

./utils: utility scripts

./exp: settings for Domino experiments

## Prerequisites

Install Azure CLI

NOTE: Azure has limitations on VMs and the access to datacenters for
different types of subscriptions or users.
Make sure that your Azure subscription has access to target datacenters and
quotas to create the expected VMs in these datacenters. 
One way to check the availability (without generating bills) is to use the web UI
trhough Azure Portal to reserve a VM but not finalize the reservation.
There could also be limitations on the total number of CPUs and other
resources.
Quotas can be increased by contacting Azure support.

## Quick start

### Azure CLI login

$ az login

### Configuration

$ cp sbin/settings-default.sh ./settings.sh

Edits settings.sh to configure a cluster, like VM locations and sizes.

NOTE: Make sure to correctly specifiy the "vm_public_key" which should be the
public key file for the private key that has been configured on Azure for VM
logins. 

### Cluster creation

$ ./sbin/cluster.sh settings.sh

The immediate output of cluster.sh may not be able to list all of the required
VMs because it will take some time to have all the VMs start running.
It would be better to check the VM status via Azure Portal to make sure that
all VMs are successfully created and running before using them.

After all VMs are running, use ./sbin/vm-ip.sh settings.sh to list all of the
VMs and their IPs.

NOTE: Once the cluster is not needed, currently you have to stop VMs and delete
resources via Azure Portal to avoid unnecessary bills. Make sure that needed
data are copied out before removing storage resources.

## Experimental settings for Domino experiments

The exp folder contains the cluster settings and scripts for the experiments
for the Domino paper.

For example, the following commands will create a cluster with 9 datacenters and 3 replicas in North America.

cd exp

./azure-cluster-setup.sh settings-Azure-NA-3r.sh

# Script usage samples

sbin/settings-default.sh specifies the settings of a cluster to create

sbin/cluster.sh creates a cluster of VMs

sbin/groups.sh creates a resource group

sbin/vnet.sh establishs and peers VNets

sbin/vm.sh creates VMs

sbin/vm-ip.sh obtains the public and private IPs of the created VMs

sbin/location.sh lists the locations of available datacenters in Azure

sbin/vnet-del.sh deletes VNets

sbin/sample.sh consists of several commonly used Azure CLI commands for information fetching

utils/ip.py splits the output of sbin/vm-ip.sh (as a file) into
.public and .private files for public and private IPs, respecitively.
The .private file will associate a port for each VM as a configuration property
for any future use.

For example:

sbin/vm-ip.sh settings.sh > vm.ip

utils/ip.py -f vm.ip -p 10001

# TODO 

Add scripts for releasing cluster resources on Azure.
