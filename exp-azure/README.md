# Build binaries

./build.sh settings.sh

If the building fails due to missing dependencies, use govendor to fetch the dependency libs.

After building, there will be 4 binaries as follows:

dynamic: the executbale server for Domino. dynamic is a alias name for Domino in current repo. TODO rename dynamic to domino

epaxos: the executable server for EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

fastpaxos: the executable server for Fast Paxos.

client: the benchmark client.

# Domino Experiments on Azure

In this document, $DOMINO represents the path to this repo's root directory.

A typical $DOMINO will be $GOPATH/src/domino/

## Creating a cluster on Azure

The scripts for creating a cluster of VMs on Azure are in the "azure" folder
under this repo's root directory.
Follow the README in the "azure" folder to use the scripts.
There have been pre-defined scripts to create clusters that are required to
repeat the experiments on Azure in the Domino paper.

NOTE: You should first make sure that your subscription of Auzre has enough
quotas to create cluster resources that are required by the experiments.
Also, Azure may experience high demand usage, and it may not be able to allocate
VMs as required. Before actually running any experiment, double check that you
have successfully started the required VMs in all target datacenters. 
One way to do so is to use the web UI via Azure Portal.

## Configurations for all experiments

cd $DOMINO/exp-azure

Generate data key file:

./sbin/key-gen.py -k 1000000 -l 8 > key.dat

Update settings.sh to replace $USER and $HOME with the corresponding directory
path on Azure VMs.

## Commit Latency with 9 datacenters and 3 replicas in North America (NA)

Create the target cluster on Azure:

cd $DOMINO/azure/exp

./azure-cluster-setup.sh settings-Azure-NA-3r.sh

This will generate a "azure-exp-vm-ip.config" file. Makes sure that this file has all of
the required VMs' information.
Since Azure may take some time to start up VMs, or it may fail to do so, it would
be better to check VMs via Azure Portal to make sure all the required VMs
(defined in settings-Azure-NA-3r.sh) successfully start before starting experiments.

Copy the created VM information to the azure experiment folder:

cp azure-exp-vm-ip.config $DOMINO/exp-azure/

Start running the experiment:

cd $DOMINO/exp-azure

./exp-commit-lat-na-9dc-3r.sh

The experiment result data files will be collected from Azure VMs and stored in
a local folder, which is pre-defined in the script.

Draw experimental result figure:

cd $DOMINO/exp-azure/fig

./gen-commit-lat-na-9dc-3r.py

A pdf file azure-commit-lat-na-9dc-3r.pdf will be generated in the current directory.

Stop the VMs and delete the cluster's resources via Azure Portal to avoid
unecessary charges.
Also, current scripts do NOT directly support re-using this cluster for
generating other experimental result figures. Manual operations are required to
do so. 
Without manual operations, it is suggested to delete the clueter resources
before switching to generate  another experiment result figure.

NOTE: stopping VMs on Azure will cause the IP addresses to change next time the
VMs restart.  As a result restarting VMs will require generating a new
azure-exp-vm-ip.config file and re-deloying location configuraion files.

## Commit Latency with 9 datacenters and 5 replicas in North America (NA)

cd $DOMINO/azure/exp

./azure-cluster-setup.sh settings-Azure-NA-5r.sh

Makes sure the "azure-exp-vm-ip.config" file has all VMs' IP addresses.

cp azure-exp-vm-ip.config $DOMINO/exp-azure/

cd $DOMINO/exp-azure 

./exp-commit-lat-na-9dc-5r.sh

cd $DOMINO/exp-azure/fig

./gen-commit-lat-na-9dc-5r.py

A pdf file azure-commit-lat-na-9dc-5r.pdf will be generated.

## Commit Latency with 6 datacenters and 3 replicas in the Globe setting

cd $DOMINO/azure/exp

./azure-cluster-setup.sh settings-Azure-Globe-3r.sh

Makes sure the "azure-exp-vm-ip.config" file has all VMs' IP addresses.

cp azure-exp-vm-ip.config $DOMINO/exp-azure/

cd $DOMINO/exp-azure 

./exp-commit-lat-globe-6dc-3r.sh

cd $DOMINO/exp-azure/fig

./gen-commit-lat-globe-6dc-3r.py

A pdf file azure-commit-lat-globe-6dc-3r.pdf will be generated.

NOTE: If continues to do the following experiments that use 6 datacenters and 3
replicas in the Globe setting, there is no need to stop the VMs or detele the
cluster resources.
The following steps assume that the experiments will continue to use this Globe setting.

## Impact of the percentile delay and the additional delays
This experiment examines the impact of the percentile delay in network
measurements for delay predictions and the additional delays on Domino's commit
Latency with 6 datacenters and 3 replicas in the Globe setting.

cd $DOMINO/exp-azure

./exp-commit-lat-globe-6dc-3r-dynamic-pth.sh

cd $DOMINO/exp-azure/fig

./gen-commit-lat-globe-6dc-3r-pth.py

A pdf file azure-commit-lat-globe-6dc-3r-pth.pdf will be generated.

## Execution Latency with 6 datacenters and 3 replicas in the Globe setting

cd $DOMINO/exp-azure

./exp-exec-lat-globe-6dc-3r.sh

cd $DOMINO/exp-azure/fig

./gen-exec-lat-globe-6dc-3r-zipf0.75.py

A pdf file azure-exec-lat-globe-6dc-3r-zipf0.75.pdf will be generated.

./gen-exec-lat-globe-6dc-3r-zipf0.95.py

A pdf file azure-exec-lat-globe-6dc-3r-zipf0.95.pdf will be generated.

## Impact of additional delays on execution latency

This experiment uses 6 datacenters and 3 replicas in the Globe setting.

cd $DOMINO/exp-azure

./exp-exec-lat-globe-6dc-3r-adddelay.sh

cd $DOMINO/exp-azure/fig

./gen-exec-lat-globe-6dc-3r-zipf0.75-adddelay.py

A pdf file azure-exec-lat-globe-6dc-3r-zipf0.75-adddelay.pdf will be generated.

## Comparing Commit Latency of Fast Paxos and Multi-Paxos

Experiments with 1 client:

cd $DOMINO/azure/exp

./azure-cluster-setup.sh settings-Azure-NA-fp-4dc-3r-1c.sh

IMPORTANT: Update azure-exp-vm-ip.config according to the README under $DOMINO/azure/exp

cp azure-exp-vm-ip.config $DOMINO/exp-azure/

cd $DOMINO/exp-azure

./exp-commit-lat-fp-na-4dc-3r-1c.sh

Experiments with 2 clients in different datacenters:

cd $DOMINO/azure/exp

./azure-cluster-setup.sh settings-Azure-NA-fp-4dc-3r-2c.sh

IMPORTANT: Update azure-exp-vm-ip.config according to the README under $DOMINO/azure/exp

cp azure-exp-vm-ip.config $DOMINO/exp-azure/

cd $DOMINO/exp-azure

./exp-commit-lat-fp-na-4dc-3r-2c.sh

cd $DOMINO/exp-azure/fig

./gen-commit-lat-fastpaxos.py

A pdf file azure-commit-lat-fastpaxos.pdf will be generated.


# Script Usage Samples

Under $DOMINO/exp-azure/sbin, there are scripts that can be used to run small
set of experiments.
Here is a list of samples about how to use these scripts.

## Build binaries
build.sh settings.sh

## Analyze Azure VM's IP into location configuration files
# Uses the closest replica for each client
gen-location.py -f azure-vm.ip -l westus2 -d azure-delay.json
# Uses the replica that achieves the lowest latency for each client
gen-location.py -f azure-vm.ip -l westus2 -d azure-delay.json -e

## Deploy files to servers and clients
./deploy.sh dynamic epaxos client default.config replica-location.config key.dat

## Start servers
# Dynamic
server.sh settings.sh start dynamic
# EPaxos without thrifty
server.sh settings.sh start epaxos e false
# EPaxos with thrifty
server.sh settings.sh start epaxos e true
# Mencius
server.sh settings.sh start epaxos m false
# Multi-Paxos
server.sh settings.sh start epaxos p false

## Stop servers
# Dynamic
server.sh settings.sh stop dynamic
# EPaxos
server.sh settings.sh stop epaxos

## Check server processes
# Dynamic
server.sh settings.sh list dynamic
# EPaxos
server.sh settings.sh list epaxos

## Start clients
# Dynamic
client.sh settings.sh start d
# EPaxos
client.sh settings.sh start e
# Mencius
client.sh settings.sh start m
# Multi-Paxos
client.sh settings.sh start p

## Stop clients
client.sh settings.sh stop

## Check client processes
client.sh settings.sh list

## Logs
# Collect both server and client logs
log.sh settings.sh collect

# Delete both server and client logs
log.sh settings.sh delete 

# List both server and client logs
log.sh settings.sh list 
