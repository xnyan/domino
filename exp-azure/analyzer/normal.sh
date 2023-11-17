#!/usr/bin/env bash
protocols=(dynamic epaxos fastpaxos mencius paxos)

for protocol in "${protocols[@]}"
do
    touch latency/$protocol.txt
    cat exp-list.config | while read LINE
    do
        exp_id=`echo $LINE | cut -d " " -f 1`
        data=`cat latency/$exp_id/$protocol/client* | grep -w "[0-1] , 1 , [0-1] , [0-9]\+ , [0-9]\+"`
        num=`echo "$data" | wc -l`
        echo -n $LINE 
        echo -n " "
        if [ "$num" == "0" ] || [ "$num" == "1" ]; then
            echo "error"
        else 
            # Calculate mean
            mean=$(echo "$data" | sed -e "s/ , /,/g" | cut -d "," -f 4 | awk '{ sum += $1 } END { printf "%.2f", sum/NR }')

            # Calculate standard deviation
            std_dev=$(echo "$data" | sed -e "s/ , /,/g" | cut -d "," -f 4 | awk -v mean="$mean" '{ sum += ($1 - mean)^2 } END { printf "%.2f", sqrt(sum/NR) }')

            echo "$mean $std_dev"
        fi
    done > latency/$protocol.txt
done
