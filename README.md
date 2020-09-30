## Requirements

Install Go 1.14

Set up $GOPATH

## Downlowd

clone this repo to $GOPATH/src/

For example, the root directory of this repo is $GOPATH/src/domino

## Source Code
This repo consists a prototype of Domino and an implementation of SMR using Fast Paxos. It also contains implementations of EPaxos, Mencius, and Multi-Paxos, which are imported from https://github.com/efficient/epaxos

/dynamic has the source code of Domino. Dynamic is an alias name of Domino in the repo because of historic reasons. TODO: rename dynamic to domino.

/fastpaxos is the source code of Fast Paxos.

/epaxos contains the source code of EPaxos, Mencius, and Multi-Paxos.

/benchmark is the source code of benchmark clients.

## Experiments on Azure

Please refer to the README file under $GOPATH/src/domino/exp-azure to run experiments on Azure.

