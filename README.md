## Requirements

Operating System: Linux Ubuntu

Install GO 1.14+

Set up $GOROOT and $GOPATH following the GO installation

Set up $PATH to include $GORROT/bin and $GOPATH/bin as follows:

export PATH="$PATH:$GOROOT/bin:$GOPATH/bin" 

## Downlowd

Clone this repo to $GOPATH/src/ as follows:

git clone https://github.com/xnyan/domino.git $GOPATH/src/domino

The root directory of this repo will be $GOPATH/src/domino

## Source Code and Data
This repo has a prototype of Domino and an implementation of Fast Paxos without fault tolerance. The repo also contains the implementations of EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

/dynamic is the source code of Domino. Dynamic is currently an alias name of Domino in the repo.

/fastpaxos is the source code of Fast Paxos.

/epaxos contains the source code of EPaxos, Mencius, and Multi-Paxos.

/benchmark is the source code of benchmark clients.

/exp-test has the scripts for testing the intallation of the prototype on a local machine.

/azure has the scripts for using Azure CLI to create a clueter across different datacenters on Azure.

/exp-azure has the scripts for repeating the experiments in the Domino paper on Azure.

/trace has the scripts (/trace/azure/fig) for using the collected data traces to generate the figures about the inter-region latency on Azure in the Domino paper. It also has the source code (/trace) and the scripts (/trace/azure) for collecting the inter-region latency on Azure.

## Quick Start

cd $GOPATH/src/domino/exp-test

Build:

./sbin/build.sh settings.sh

If the building fails due to missing dependencies, use [govendor](https://github.com/kardianos/govendor) to fetch the dependency libs. An alterative is to update to GO 1.15+ to use [GO modules](https://blog.golang.org/migrating-to-go-modules).

Start Domino replica servers:

./server.sh start domino

Start a Domino client:

./client.sh domino

After the client completes running, stop the Domino replica servers:

./server.sh stop domino

Follow [the README file under /exp-test](https://github.com/xnyan/domino/tree/master/exp-test) to run differet protocols locally to test the installation.

## Data Traces of the Inter-Region Latency on Azure

Use the following two commands to download the data traces that are collected from Azure and used in the Domino paper. The data traces are collected under the Globe setting (6 datacenters that are globally distributed) and the NA setting (9 datacenters that are located in North America), respectively.

Data trace under the Globe setting:

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz

Data trace under the NA setting:

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-na-9dc-24h-202005071450-202005081450.tar.gz

The data traces are byte files. Follow [the README file in /trace/azure/fig](https://github.com/xnyan/domino/tree/master/trace/azure/fig) to parse the files and generate the figures about the inter-region latency measurments in the Domino paper.


## Experiments on Azure

Follow [the README file under /exp-azure](https://github.com/xnyan/domino/tree/master/exp-azure) to repeat the experiments in the Domino paper on Azure. The scripts can also be used to run the prototype in a distributed setting in a local cluster by modifying configuration files to specify machine IP addresses.
