Uses azure-cluster-setup.sh to create cluster on Azure.

For example,

./azure-cluster-setup.sh settings-Azure-Globe-3r.sh 

azure-cluster-setup.sh will run ../sbin/vm-ip.sh to generate file
azure-exp-vm-ip.config, which consists of VMs' public and private IPs.

The public IPs are used for uploading files to Azure and controlling VMs from
outside of Azure, while the private IPs will be used for communications between
VMs.

In the azure-exp-vm-ip.config file, VMs that has the prefix of vm1- will run
clients, and VMs that has the prefix of vm2- will run replica servers.

For experiments to compare the latency of Fast Paxos and Multi-Paxos, we
currently need to manually config the generated azure-exp-vm-ip.config file,
which is as follows:

(1) After using settings-Azure-NA-fp-4dc-3r-1c.sh to create a cluster, replace vm1-
with vm2- for vm1-eastus2, vm1-westus2, and vm1-canadaeast in the generated
azure-exp-vm-ip.config file.

(2) Similarly, after using settings-Azure-NA-fp-4dc-3r-2c.sh to create a cluster,
replace vm1- with vm2- for vm1-eastus2 and vm1-canadaeast in the generated
azure-exp-vm-ip.config file.
