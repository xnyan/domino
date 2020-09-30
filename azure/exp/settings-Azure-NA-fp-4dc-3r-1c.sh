## Resource group settings
resource_group_name="na-9dc-3r-exp"
## Location options: az account list-locations --output table
## uses the value under name / Name
resource_group_location="eastus2"

## Virtual network settings
vnet_mask="16"
vnet_subnet_mask="24"
vnet_subnet_name="Subnet1"

## Virtual machine settings
## VM public key file for login, which must corresponds to the private key that
## is configured on Azure for VM logins.
vm_public_key="$HOME/.ssh/id_rsa.pub"
## VM username
vm_username="$USER"
## VM image options: az vm image list --output table
## Uses the value under Urn or UrnAlias
vm_image="UbuntuLTS"
## VM size options are different across locations: az vm list-sizes --location <location> --output table
## Uses the value under name / Name
vm_size="Standard_D4_v3"
# Accelerated Network Options: true for Standard_D4_v3
vm_acc_network="true"
## Disk options:
## Standard_LRS (for HDD)
## Premium_LRS 
## StandardSSD_LRS (for standard SSD)
## UltraSSD_LRS
vm_disk_type="StandardSSD_LRS"

## Optional configurations
## vm_dns=true to set a DNS for each VM
vm_dns="false"
## DNS format: ${vm_name}-dns.${location}.${vm_dns_suffix}
## Suffix cannot be changed
vm_dns_suffix="cloudapp.azure.com"
## vm_no_wait=true to create VMs in the background. When this is enabled, if
## Azure fails to create a VM, the error information will not be displyed by
## Azure CLI.
vm_no_wait="false"

## azure-location customized-tag vnet_ip number-of-vms
cluster_config=(
## Fast Paxos Experiment
## 4 datacenters, 3 replicas (eastus2, westus2, canadaeast), 1 client (centralus)
"eastus2                eastus2            10.1.0.0    1"
"westus2                westus2            10.3.0.0    1"
"centralus              centralus          10.6.0.0    1"
"canadaeast             canadaeast         10.5.0.0    1"
)
