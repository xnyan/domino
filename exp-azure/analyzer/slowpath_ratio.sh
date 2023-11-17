#!/usr/bin/env bash
protocols=(epaxos fastpaxos)

touch latency/$protocol.txt
cat exp-list.config | while read LINE
do
    for protocol in "${protocols[@]}"
    do
        #grep -r -w "[0-1] , [0-1] , [0-1] , [0-9]\+ , [0-9]\+"
        exp_id=`echo $LINE | cut -d " " -f 1`
        if [ "$protocol" == "epaxos" ]; then
            slowpath=`cat latency/$exp_id/$protocol/client*.log | grep -w "0 , 1 , 0 , [0-9]\+ , [0-9]\+"`
            fastpath=`cat latency/$exp_id/$protocol/client*.log | grep -w "0 , 1 , 1 , [0-9]\+ , [0-9]\+"`
        else
            slowpath=`cat latency/$exp_id/$protocol/client*.log | grep -w "1 , 1 , 0 , [0-9]\+ , [0-9]\+"`
            fastpath=`cat latency/$exp_id/$protocol/client*.log | grep -w "1 , 1 , 1 , [0-9]\+ , [0-9]\+"`
        fi
        slow_num=$(echo "$slowpath" | wc -l)
        fast_num=$(echo "$fastpath" | wc -l)
        total=$((slow_num + fast_num))
        slow_ratio=$(awk "BEGIN { ratio = $slow_num / $total; printf \"%.2f\", ratio }")

        #echo "Total: $total"
        #echo "Slow Num: $slow_num"
        #echo "Fast Num: $fast_num"
        echo -n "$protocol,$slow_ratio,"
    done
    echo
done
