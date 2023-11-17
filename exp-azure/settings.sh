username="$USER"
identity=""

is_debug_mode="true"

# Node configuration file (same file name in both remote and local)
config_file="default.config"
# Location configuration files (same file name in both remote and local)
remote_server_location_file="server-location.config"
replica_location_file="replica-location.config"
remote_client_location_file="client-location.config"

# Settings for remote machines to run experiments
remote_setup="source $HOME/.profile"
# The directory to store binary and configuration files
remote_exec_path="$HOME/"
# The directory to store server logs
remote_server_log_dir="$HOME/"
# The directory to store client logs
remote_client_log_dir="$HOME/"

# Settings for the local machine to (temporarily) collect experimental results
local_log_dir="./log"

## Domino server
# Source code path, which is a relative path to $GOPATH/src
dynamic_go_app_path="domino/dynamic/server"
# The name of the binary after "go build"
dynamic_go_app="server"
# The name of the binary to use for experiments
dynamic_app="dynamic"

## EPaxos, Menciu, and Multi-Paxos server
# Source code path, which is a relative path to $GOPATH/src
epaxos_go_app_path="domino/epaxos/server"
# The name of the binary after "go build"
epaxos_go_app="server"
# The name of the binary to use for experiments
epaxos_app="epaxos"

## Fast Paxos server
# Source code path, which is a relative path to $GOPATH/src
fastpaxos_go_app_path="domino/fastpaxos/server"
# The name of the binary after "go build"
fastpaxos_go_app="server"
# The name of the binary to use for experiments
fastpaxos_app="fastpaxos"

## Client
# Source code path, which is a relative path to $GOPATH/src
client_go_app_path="domino/benchmark/benchmark-client"
# The name of the binary after "go build"
client_go_app="benchmark-client"
# The name of the binary to use for experiments
client_app="client"