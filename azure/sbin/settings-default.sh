## Resource group settings
resource_group_name="testing-cluster"
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
"eastus2      eastus2       10.1.0.0    2"
"westus       westus        10.2.0.0    1"
#"westus2      westus2       10.3.0.0    1"
)

#### A list of Azure locations from location.sh ####
### Asia-Pacific
#eastasia
#southeastasia
#japanwest
#japaneast
#koreacentral
#koreasouth
#australiaeast
#australiasoutheast
#australiacentral
#australiacentral2
#southindia
#centralindia
#westindia
### North America
#eastus
#eastus2
#westus
#westus2
#centralus
#westcentralus
#northcentralus
#southcentralus
#canadacentral
#canadaeast
### Europe
#northeurope
#westeurope
#uksouth
#ukwest
#francecentral
#francesouth
#switzerlandnorth
#switzerlandwest
#germanynorth
#germanywestcentral
#norwaywest
#norwayeast
### Middle East
#uaecentral
#uaenorth
### North America
#brazilsouth
### Africa
#southafricanorth
#southafricawest
