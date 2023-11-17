#!/usr/bin/env bash
protocols=(dynamic epaxos fastpaxos mencius paxos)

for protocol in "${protocols[@]}"
do
    touch latency/$protocol.txt
    cat exp-list.config | while read LINE
    do
        exp_id=`echo $LINE | cut -d " " -f 1`
        data=`cat latency/$exp_id/$protocol/server* | grep -a --binary-files=text -E '[[:digit:]]+ ns'`
        num=`echo "$data" | wc -l`
        echo -n $LINE 
        echo -n " "
        if [ "$num" == "0" ] || [ "$num" == "1" ]; then
            echo "error"
        else 
            # 平均値の計算
            average=$(echo "$data" | rev | cut -d " " -f 2 | rev | awk '{sum+=$1} END {printf "%.0f", sum/NR}')

            # 標準偏差の計算
            variance=$(echo "$data" | rev | cut -d " " -f 2 | rev | awk -v avg="$average" '{sum+=($1-avg)^2} END {printf "%.0f", sum/NR}')
            stddev=$(echo "$variance" | awk '{printf "%.0f", sqrt($0)}')

            echo -n "$average "
            echo -n "$stddev "
            echo "$num"
        fi
    done > latency/$protocol.txt
done
