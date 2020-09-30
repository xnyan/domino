
#data_dir="../trace/trace-azure-globe-6dc-24h-202005170045-202005180045"

data_dir=$1
if [ -z $1 ]; then
  echo "Missing the log file direcotry. Usage: <log_dir>"
  exit 1
fi

log_list=(
eastus2-westus2.log
eastus2-francecentral.log
eastus2-australiaeast.log
#eastus2-eastasia.log
#eastus2-southeastasia.log
)

for l in ${log_list[@]}
do
  lf="$data_dir/$l"
  echo "Parsing $lf"
  ls -lh $lf
  parser -f $lf
  echo "Parsed $lf to ${lf}.txt"
done
