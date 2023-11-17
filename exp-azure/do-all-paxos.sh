#!/usr/bin/env bash

wait_client_processes(){
    timer=0 
    while true
    do
        remaining_processes=`./sbin/client.sh settings.sh list | grep ./client | wc -l`
        if [[ $remaining_processes -eq 0 ]]; then
            echo "Client process is terminated"
            break
        else
            echo "Client process is running"
        fi
        sleep 10
        ((timer=timer+1))
        if [ $timer -eq 10 ]; then
            ./sbin/client.sh settings.sh stop
        fi
    done
}
wait_server_processes(){
    timer=0 
    while true
    do
        remaining_processes=`./sbin/server.sh settings.sh list $1 | grep ./server | wc -l`
        if [[ $remaining_processes -eq 0 ]]; then
            echo "Server process is terminated"
            break
        else
            echo "Server process is running"
        fi
        sleep 10
        ((timer=timer+1))
        if [ $timer -eq 10 ]; then
            ./sbin/server.sh settings.sh stop $1 
        fi
    done
}
./deploy.sh dynamic fastpaxos epaxos client default.config replica-location.config key.dat
./sbin/log.sh settings.sh delete

case $2 in
    "noLeader")
        # Epaxos exp
        ./sbin/server.sh settings.sh start epaxos e false
        ./sbin/client.sh settings.sh start e
        wait_client_processes
        ./sbin/server.sh settings.sh stop epaxos
        wait_server_processes epaxos
        ./sbin/log.sh settings.sh collect
        mkdir -p latency/$1/epaxos
        mv log/* latency/$1/epaxos/

        # Mencius exp
        ./sbin/server.sh settings.sh start epaxos m true
        ./sbin/client.sh settings.sh start m
        wait_client_processes
        ./sbin/server.sh settings.sh stop epaxos
        wait_server_processes epaxos
        ./sbin/log.sh settings.sh collect
        mkdir -p latency/$1/mencius
        mv log/* latency/$1/mencius/
        ;;
    "Leader")
        # Domino exp
        ./sbin/server.sh settings.sh start dynamic
        ./sbin/client.sh settings.sh start d
        wait_client_processes
        ./sbin/server.sh settings.sh stop dynamic
        wait_server_processes dynamic
        ./sbin/log.sh settings.sh collect
        mkdir -p latency/$1/dynamic
        mv log/* latency/$1/dynamic/

        # MultiPaxos exp
        ./sbin/server.sh settings.sh start epaxos p false
        ./sbin/client.sh settings.sh start p
        wait_client_processes
        ./sbin/server.sh settings.sh stop epaxos
        wait_server_processes epaxos
        ./sbin/log.sh settings.sh collect
        mkdir -p latency/$1/paxos
        mv log/* latency/$1/paxos/

        # FastPaxos exp
        ./sbin/server.sh settings.sh start fastpaxos
        ./sbin/client.sh settings.sh start fp
        wait_client_processes
        ./sbin/server.sh settings.sh stop fastpaxos
        wait_server_processes fastpaxos
        ./sbin/log.sh settings.sh collect
        mkdir -p latency/$1/fastpaxos
        mv log/* latency/$1/fastpaxos/
        ;;
    *)
esac
