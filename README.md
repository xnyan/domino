## Requirements

Operating Systems: Linux Ubuntu

Install Go 1.14+

Set up $GOROOT and $GOPATH following the Go language installation

Set up $PATH to include $GORROT/bin and $GOPATH/bin as follows:

export PATH="$PATH:$GOROOT/bin:$GOPATH/bin" 

## Downlowd

Clone this repo $GOPATH/src/ as follows:

git clone https://github.com/xnyan/domino.git $GOPATH/src/domino

The root directory of this repo will be $GOPATH/src/domino

## Source Code
This repo has a prototype of Domino and an implementation of Fast Paxos without fault tolerance. The repo also contains the implementations of EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

/exp-test has the scripts for quick start and testing the intallation.

/dynamic is the source code of Domino. Dynamic is currently an alias name of Domino in the repo.

/fastpaxos is the source code of Fast Paxos.

/epaxos contains the source code of EPaxos, Mencius, and Multi-Paxos.

/benchmark is the source code of benchmark clients.

/azure has scripts of using Azure CLI to create a clueter on Azure.

/exp-azure has the scripts to repeat the experiments on Azure in the Domino paper.

/trace has the scripts (/trace/azure/fig) to generate the figures about inter-region latency on Azure in the paper based on the collected data traces. It also has the source code and the scripts (/trace/azure) for collecting the inter-region latency.

## Quick Start

Follow [the README file under /exp-test](https://github.com/xnyan/domino/tree/master/exp-test) to run differet protocols locally for testing the installation.

## Data Traces for Inter-Region Latency on Azure

Use the following two commands to download the data traces for the Globe setting (6 datacenters that are globally distributed) and the NA setting (9 datacenters that are located in North America).

Globe setting data trace:

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz

NA setting data trace:

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-na-9dc-24h-202005071450-202005081450.tar.gz

The data traces are byte files. Follow [the README file in /trace/azure/fig](https://github.com/xnyan/domino/tree/master/trace/azure/fig) to parse the files and generate the figures about the inter-region latency measurments in the Domino paper.


## Experiments on Azure

Refer to [the README file under /exp-azure](https://github.com/xnyan/domino/tree/master/exp-azure) to run the experiments on Azure in the Domino paper.
