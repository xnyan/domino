#!/usr/bin/env bash
protocols=(epaxos fastpaxos)

for protocol in "${protocols[@]}"
do
    touch latency/$protocol.txt
    cat exp-list.config | while read LINE
    do
        exp_id=`echo $LINE | cut -d " " -f 1`
        clients=`echo $LINE | cut -d " " -f 2` 
        replicas=`echo $LINE | cut -d " " -f 3` 
        IFS=',' read -ra client <<< "$clients"
        total=0
        average=0
        declare -a stddevs
        for c in "${client[@]}"; do
            if [ "$protocol" == "epaxos" ]; then
                data=`cat latency/$exp_id/$protocol/client*-$c-*.log | grep -w "0 , 1 , 1 , [0-9]\+ , [0-9]\+"`
            else
                data=`cat latency/$exp_id/$protocol/client*-$c-*.log | grep -w "1 , 1 , 1 , [0-9]\+ , [0-9]\+"`
            fi
            num=`echo "$data" | wc -l`
            if [ "$num" == "0" ] || [ "$num" == "1" ]; then
                err=true
                break
            fi
            err=false
            av=`echo "$data" | sed -e "s/ , /,/g" | cut -d "," -f 4 | sort -n | awk 'BEGIN{OFMT="%d"} { a[i++]=$1; } END { x=int((i+1)/2); if (x < (i+1)/2) print (a[x-1]+a[x])/2; else print a[x-1]; }'`
            total=$((total + av))
            stddev=$(echo "$data" | sed -e "s/ , /,/g" | cut -d "," -f 4 | awk -v mean="$av" '{ sum += ($1 - mean)^2 } END { printf "%.2f", sqrt(sum/NR) }')
            stddevs+=("$stddev")
        done
        average=$(awk "BEGIN {printf \"%.0f\", $total/${#client[@]}}")  # total を client の数で割る
        echo -n $LINE
        echo -n " "
        if [ "$err" == "false" ]; then
            # Calculate the average of the standard deviations
            stddev_avg=$(echo "${stddevs[*]}" | awk -v RS=" " '{sum+=$1} END {printf "%.2f", sum/NR}')
            echo "$average $stddev_avg"
        else
            echo "error"
        fi
    done > latency/$protocol.txt
done