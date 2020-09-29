ip_file="azure-exp-vm-ip.config"
t=60

settings=$1

# Creates a cluster on Azure
../sbin/cluster.sh $settings


echo "Waiting $t seconds for VMs to start up."
sleep $t

# Output VM IPs into a file
../sbin/vm-ip.sh $settings > ${ip_file}

echo "List VM IPs"
cat ${ip_file}

echo "Output VM IPs to file ${ip_file}"
echo "NOTE: Azure may have delays in starting up a VM. Make sure that all VMs successfully start running on Azure and their IPs are in ${ip_file}"
