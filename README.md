## Domino

Domino is a low-latency state machine replication protocol in WANs. 

This repo includes the implemenation of Domino and the network measruement data traces on Microsoft Azure for our [CoNEXT'20](https://conferences2.sigcomm.org/co-next/2020) paper, ["Domino: Using Network Measurements to Reduce State Machine Replication Latency in WANs."](https://dl.acm.org/doi/10.1145/3386367.3431291)

## Prerequisites

Operating System: Linux Ubuntu

Install GO 1.13+

Set up $GOROOT and $GOPATH following the GO installation

Set up $PATH to include $GOROOT/bin and $GOPATH/bin as follows:

export PATH="$PATH:$GOROOT/bin:$GOPATH/bin" 

## Downlowd

Clone this repo to $GOPATH/src/ as follows:

git clone https://github.com/xnyan/domino.git $GOPATH/src/domino

The root directory of this repo will be $GOPATH/src/domino

## Source Code
This repo has a prototype of Domino and an implementation of Fast Paxos without fault tolerance. The repo also contains the implementations of EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

/dynamic is the source code of Domino. Dynamic is currently an alias name of Domino in the repo.

/fastpaxos is the source code of Fast Paxos.

/epaxos contains the source code of EPaxos, Mencius, and Multi-Paxos.

/benchmark is the source code of benchmark clients.

/exp-test has the scripts for testing the intallation of the prototype on a local machine.

/azure has the scripts for using Azure CLI to create a clueter across different datacenters on Azure.

/exp-azure has the scripts for repeating the experiments on Azure in the Domino paper.

/trace has the scripts (/trace/azure/fig) for using the collected data traces to generate the figures about the inter-region latency on Azure in the Domino paper. It also has the source code and the scripts (/trace/azure) for collecting the inter-region latency on Azure.

## Quick Start

cd $GOPATH/src/domino/exp-test

Build:

./sbin/build.sh settings.sh

NOTE: The source code of dependency libs are already under /vendor, which are fetched by using [govendor](https://github.com/kardianos/govendor). Although [GO modules](https://blog.golang.org/migrating-to-go-modules) are widely adopted since GO 1.15, there should be no need to use GO modules for Domino as long as this repo is cloned at $GOPATH/src/domino. For users that prefer using GO modules, run "go mod init" under $GOPATH/src/domino before building executables.

Start Domino replica servers:

./server.sh start domino

Start a Domino client:

./client.sh domino

After the client completes running, stop the Domino replica servers:

./server.sh stop domino

Follow [the README file under /exp-test](https://github.com/xnyan/domino/tree/master/exp-test) to run differet protocols locally to test the installation.

## Data Traces about the Inter-Region Latency on Azure

Use the following two commands to download the data traces that are collected from Azure and used in the Domino paper. The data traces are collected under the Globe setting (6 datacenters that are globally distributed) and the NA setting (9 datacenters that are located in North America), respectively.

Data trace under the Globe setting:

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz

Data trace under the NA setting:

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-na-9dc-24h-202005071450-202005081450.tar.gz

NOTE: Extracting both of the two .tar.gz files would require about 50 GB disk spaces.

The data traces are plain text files. Each data file is named as "{host_region}-{target-region}.log.txt". The first line in each data file is a comment. After that each line consists of three timestamps: a client's sending time of a probing request, the time when the client receives the probing response from the target server, and the time when the server receives the probing request. All of the three timestamps are in nanoseconds.

To re-generate the figures about the inter-region latency on Azure in the Domino paper, please follow [the README file in /trace/azure/fig](https://github.com/xnyan/domino/tree/master/trace/azure/fig).


## Experiments on Azure

Follow [the README file under /exp-azure](https://github.com/xnyan/domino/tree/master/exp-azure) to repeat the experiments on Azure in the Domino paper.
