#!/bin/bash

# set -x

hosts_list="vm-list.config"

# use UTC throughout the experiment
export TZ=UTC

# results will be stored here
base_dir=rtt

# ssh_option="-i ~/.ssh/id_rsa -y"
ssh_option="-i ~/.labssh/.ssh/id_rsa -y"

# regions=(ohio virginia california oregon mumbai seoul singapore sydney tokyo canada frankfurt ireland london saopaulo)
regions=($(cut -d " " -f 1 "$hosts_list"))

user=koya

server_command="ping -c 60"
now() {
    echo -n "$(date +%s) $(date +'%F %H:%M:%S %:z')"
}

mkdir -p "$base_dir"
output_dir="$base_dir/$(date +'%Y%m%d-%H%M')"
mkdir -p "$output_dir"

# start time
echo "> Experiments-Start: $(now)"

for i in "${regions[@]}"; do
    server_ip=$(grep "$i" "$hosts_list" | cut -d ' ' -f 2)
    mkdir -p "$output_dir/$i"

    echo ">> $i-Region-Start: $(now)"

    for j in "${regions[@]}"; do
        client_ip=$(grep "$j" "$hosts_list" | cut -d ' ' -f 2)
        output="$output_dir/$i/$j.txt"

        echo ">>> $i-$j: $(now)"

        # start client
        ssh -oStrictHostKeyChecking=no $ssh_option "$server_ip" -l $user $server_command "$client_ip" > "$output" &
    done

    echo ">> $i-Region-End: $(now)"
done

# Wait for all background processes to finish
wait

# finish time
echo "> Experiments-Finish: $(now)"
