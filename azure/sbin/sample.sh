## Location options:
az account list-locations --output table

## VM image options: 
az vm image list --output table

## VM size options are different across locations:
az vm list-sizes --location <location> --output table

## List the details of all of the VMs
az vm list -g cluster-us

## List private and public IPs for all VMs or a specific VM
az vm list-ip-addresses -g cluster-us --output table
az vm list-ip-addresses -g cluster-us -n vm1-eastus --output table

## Get the details of the given VM
az vm get-instance-view -g cluster-us -n vm1-eastus

## Get a list of your subscriptions
az account list --output table

## Use az account set with the subscription ID or name you want to switch to
az account set --subscription "My Demos"
