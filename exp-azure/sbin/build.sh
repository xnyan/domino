#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

source $sbin/common.sh

setting_file=$1
load_settings ${setting_file}

cur_dir=`pwd`

#Dynamic
local_dynamic_path="$GOPATH/src/${dynamic_go_app_path}"
cd $local_dynamic_path
log "Building Dynamic in $local_dynamic_path"
CGO_ENABLED=0 go build

#EPaxos
local_epaxos_path="$GOPATH/src/${epaxos_go_app_path}"
cd $local_epaxos_path
log "Building EPaxos in $local_epaxos_path"
CGO_ENABLED=0 go build

#Fast Paxos
local_fastpaxos_path="$GOPATH/src/${fastpaxos_go_app_path}"
cd $local_fastpaxos_path
log "Building FastPaxos in $local_fastpaxos_path"
CGO_ENABLED=0 go build

#Client
local_client_path="$GOPATH/src/${client_go_app_path}"
cd $local_client_path
log "Building Client in $local_client_path"
CGO_ENABLED=0 go build

cd ${cur_dir}
mv $local_dynamic_path/${dynamic_go_app} ./${dynamic_app}
mv $local_epaxos_path/${epaxos_go_app} ./${epaxos_app}
mv $local_fastpaxos_path/${fastpaxos_go_app} ./${fastpaxos_app}
mv $local_client_path/${client_go_app} ./${client_app}
