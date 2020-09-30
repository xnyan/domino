#!/usr/bin/env bash

sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

if [ -z $GOPATH ]; then
  echo "ERROR: GOPATH is not set."; exit 1
fi

usage()
{
  echo -e "$@"
}

load_settings() {
  usage_info="Usage: <settings.sh>
  \t Please refer to $sbin/settings-default.sh to write a configuration file"
  
  setting=$1
  if [ -z $setting ]; then
    usage "${usage_info}"; exit 1
  fi

  source $setting

  if [ ! -z $identity ]; then
    SSH_OPTIONS="-i $identity"
  fi

  if [ ! -z $username ]; then
    USER_AT="${username}@"
  fi
}

install_app() {
  app__path=$1
  echo "Uses $sbin/install.sh to build and install $app__path."; 
  $sbin/install.sh -a ${app__path}
  if [ $? != 0 ]; then
    exit 1
  fi
}

COMMON__CONFIG_LIST=()
parse_config_file() {
  #Returns via using a shared variable
  while read line;
  do
    config=`echo $line | sed "s/#.*$//;/^$/d"`
    if [ -z "$config" ]; then
      continue
    fi
    COMMON__CONFIG_LIST=("${COMMON__CONFIG_LIST[@]}" "$config")
  done<$@
}

run_cmd() {
  cmd="$@"
  eval $cmd
}

run_cmd_in_background() {
  cmd="$@"
  eval $cmd &
}
