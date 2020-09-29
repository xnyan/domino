#!/usr/bin/env bash

sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

usage()
{
  echo -e "$@"
}

load_setting() {
  __usage__info__="Usage: <config file>
  \t Please refer to ${sbin}/settings.sh to write a configuration file"
  
  __setting__=$1
  if [ -z $__setting__ ]; then
    usage "${__usage__info__}"; exit 1
  fi

  source ${__setting__}

  if [ -z $vm_public_key ]; then
    echo "Missing vm_public_key in ${__setting__}"; exit 1
  fi

  if [ -z $vm_username ]; then
    echo "Missing vm_username in ${__setting__}"; exit 1
  fi
}

__PARSE__RET__=()
parse_list_file() {
  __PARSE__RET__=()
  # Returns parsing results through a global shared variable
  while read line;
  do
    config=`echo $line | sed "s/#.*$//;/^$/d"`
    if [ -z "$config" ]; then
      continue
    fi
    __PARSE__RET__=("${__PARSE__RET__[@]}" "$config")
  done<$@
}

run_cmd() {
  cmd="$@"
  #echo "Executing: $cmd"
  eval $cmd
}

run_cmd_in_background() {
  cmd="$@"
  #echo "Executing (in background): $cmd"
  eval $cmd &
}

log() {
  echo "`date`: $@"
}

gen_vm_name() {
  if [ -z $1 ] || [ -z $2 ]; then
    echo "Missing vm name parameter(s)."; exit 1
  fi
  echo "vm$1-$2"
}

gen_dns_name() {
  if [ -z $1 ]; then
    echo "Missing vm name as a parameter."; exit 1
  fi
  echo "$1-dns"
}

gen_vnet_name() {
  if [ -z $1 ]; then
    echo "Missing vnet name parameter."; exit 1
  fi
  echo "vnet-$1"
}

gen_vnet_peer_name() {
  if [ -z $1 ]; then
    echo "Missing vnet peer name parameter(s). Expected two parameters"; exit 1
  fi
  echo "peer-$1-$2"
}

