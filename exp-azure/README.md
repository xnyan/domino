# Domino Experiments on Azure

## Build binaries

./build.sh settings.sh

If the building fails due to missing dependencies, use govendor to fetch the dependency libs.

After building, there will be 4 binaries as follows:

dynamic: the executbale server for Domino. dynamic is a alias name for Domino in current repo. TODO rename dynamic to domino

epaxos: the executable server for EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

fastpaxos: the executable server for Fast Paxos.

client: the benchmark client.

## Configurations
Generates data key file:

./sbin/key-gen.py -k 1000000 -l 8 > key.dat

## Commit latency on Azure

## Execution latency on Azure

# Script Usage Samples
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
# Mencius with early commit
server.sh settings.sh start epaxos m true
# Multi-Paxos without thrifty
server.sh settings.sh start epaxos p false
# Multi-Paxos with thrifty
server.sh settings.sh start epaxos p true

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
