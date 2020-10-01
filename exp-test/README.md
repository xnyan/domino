# Quick Start

## Build

./sbin/build.sh settings.sh

This will generates 4 binary files, client, fastpaxos, dynamic, epaxos.

dynamic is an alias name for domino.
epaxos has the protocols of EPaxos, Mencius, and Multi-Paxos.

If the building process fails due to missing dependencies, use govendor to
fetch the dependency libs.

## Test Domino

Start Domino replica servers:

./server.sh start domino

Start a Domino client that targets sending 2 requests per second for 10s:

./client.sh domino

Stop Domino replica servers:

./server.sh stop domino

## Test Fast Paxos

./server.sh start fastpaxos

./client.sh fastpaxos

./server.sh stop fastpaxos

## Test EPaxos

./server.sh start epaxos

./client.sh epaxos

./server.sh stop epaxos

## Test Mencius

./server.sh start mencius

./client.sh mencius

./server.sh stop mencius

## Test Multi-Paxos

./server.sh start paxos

./client.sh paxos

./server.sh stop paxos

