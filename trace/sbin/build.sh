setting=$1

if [ -z "$setting" ]; then
  echo "Usage: <settings.sh>"; exit 1
fi

source $setting

cur_dir=`pwd`

server_path="$GOPATH/src/${server_app_path}/${server_app}"
cd $server_path
go build

client_path="$GOPATH/src/${client_app_path}/${client_app}"
cd ${client_path}
go build

cd ${cur_dir}
mv ${server_path}/${server_app} ./
mv ${client_path}/${client_app} ./
