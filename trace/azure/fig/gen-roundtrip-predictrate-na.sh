dc_l=(
eastus2
canadacentral
canadaeast
centralus
northcentralus
southcentralus
westcentralus
westus2
westus
)

data_dir="trace-azure-na-9dc-24h-202005071450-202005081450"
output_dir="roundtrip-predictrate-na/"

mkdir -p ${output_dir}

for src_dc in ${dc_l[@]}
do
  for dst_dc in ${dc_l[@]}
  do
    if [ "${src_dc}" == "${dst_dc}" ]; then
      continue
    fi
    echo ${src_dc} ${dst_dc} ${data_dir} ${output_dir}
    ./gen-roundtrip-predictrate-by-host.py ${src_dc} ${dst_dc} ${data_dir} ${output_dir}
  done
done
