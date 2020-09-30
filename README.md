## Requirements

Operating Systems: Ubuntu 12.04+

Install Go 1.14+

Set up $GOPATH

## Downlowd

clone this repo to $GOPATH/src/

For example, the root directory of this repo is $GOPATH/src/domino

## Source Code
This repo has a prototype of Domino and an implementation of Fast Paxos without fault tolerance. The repo also contains the implementations of EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

/dynamic is the source code of Domino. Dynamic is an alias name of Domino in the repo because of historic reasons. TODO: rename dynamic to domino.

/fastpaxos is the source code of Fast Paxos.

/epaxos contains the source code of EPaxos, Mencius, and Multi-Paxos.

/benchmark is the source code of benchmark clients.

/azure has scripts of using Azure CLI to create a clueter on Azure.

/exp-azure has the scripts to repeat the experiments on Azure in the Domino paper.

## Experiments on Azure

Please refer to the README file under /exp-azure to run experiments on Azure, which also shows how to build the binaries.

