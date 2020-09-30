#Remote settings
username="$USER"
identity=""
remote_setup="source $HOME/.profile"
remote_exec_path="./"
##NOTE: This log dir must be the same as the dir (in ${config_file}) for probing log
remote_log_dir="$remote_exec_path"
##Maps from datancenter IDs to server ips. 
##A server's public ip may be different from the one used in the program
remote_location_file="location.config"

#Program parameters
is_debug_mode="false"
config_file="exp.config"
location_file="location.config"

#Server binary
server_app="server"
#Relative path to $GOPATH/src
server_app_path="latency"

#Client binary
client_app="client"
#Relative path to $GOPATH/src
client_app_path="latency"

#Log collection local directory
local_log_dir="./log"
