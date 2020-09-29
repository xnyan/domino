#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"
source $sbin/common.sh

usage_info="Usage: [-l|t|a] [-l for name list (default) ] | [-t for table format] | [-a for all descriptions]"

mode=$1

##List Azure locations
if [ -z $mode ]; then
  usage $usage_info
  az account list-locations | grep name
elif [ "$mode" == "-l" ]; then
  az account list-locations | grep name
elif [ "$mode" == "-a" ]; then
  az account list-locations
elif [ "$mode" == "-t" ]; then
  az account list-locations --output table
else
  usage $usage_info
fi
