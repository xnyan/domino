# Build binarys

./build.sh

Binary files, client and server, will be created.

# Network measurements for collecting inter-region latency on Azure

## Configuration

Creates a clueter on Azure by using the readme and the scripts in $GOPATH/src/domino/azure

Update location.config with datacenter IDs and VM private IPs

Update remote.config with datacenter IDs and VM public IPs

Update exp.config to specify probing intervals and durations, replace the $HOME
with the target log directory on Azure VMs.

Configures settings.sh to specify username and experiment directories on Azure
VMs by replacing the $USER and $HOME fields

## Start to run the network measurements

./exp.sh

## Parsing collected data

The collected data are binary files. Use the parser in $GOPATH/src/domino/trace/parser to parse convert a binary file to a plain-text file.
