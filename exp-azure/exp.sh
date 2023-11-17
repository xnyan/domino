#!/usr/bin/env bash
vm_list=`cat vm-list.config`
start_time=`date +%s`
cat exp-list.config | while read LINE
do
    exp_id=`echo $LINE | cut -d " " -f 1`
    clients=`echo $LINE | cut -d " " -f 2`
    leader=`echo $LINE | cut -d " " -f 3 | cut -d "," -f 1`
    replicas=`echo $LINE | cut -d " " -f 3`
    replica_arr=(${replicas//,/ })
    client_arr=(${clients//,/ })
    server_id=1
    for r in "${replica_arr[@]}";
    do
        echo $server_id $r `echo "$vm_list" | grep $r | cut -d " " -f 2`
        ((server_id=server_id+1))
    done > server-location.config
    server_id=1
    port_num=10001
    for r in "${replica_arr[@]}";
    do
        if [ $r == $leader ]
        then
            echo $server_id $r `echo "$vm_list" | grep $r | cut -d " " -f 3` $port_num L L
        else
            echo $server_id $r `echo "$vm_list" | grep $r | cut -d " " -f 3` $port_num L F
        fi
        ((server_id=server_id+1))
        ((port_num=port_num+1))
    done  > replica-location.config
    client_id=1
    for c in "${client_arr[@]}";
    do
        echo $client_id $c `echo "$vm_list" | grep $c | cut -d " " -f 2` 1 $leader
        ((client_id=client_id+1))
    done > client-location.config
    ./do-all-paxos.sh $exp_id "Leader"

    client_id=1
    for c in "${client_arr[@]}";
    do
        leader=`python3 leader_opt.py $c $replicas`
        echo $client_id $c `echo "$vm_list" | grep $c | cut -d " " -f 2` 1 $leader
        ((client_id=client_id+1))
    done > client-location.config
    ./do-all-paxos.sh $exp_id "noLeader"

    end_time=`date +%s`
    run_time=$((end_time - start_time))
    echo "exp id $exp_id finish time --------> $run_time"

done
end_time=`date +%s`
run_time=$((end_time - start_time))
echo "all finish time --------> $run_time"