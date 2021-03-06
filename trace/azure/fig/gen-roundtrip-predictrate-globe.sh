dc_l=(
eastus2
westus2
francecentral
australiaeast
eastasia
southeastasia
)

data_dir="trace-azure-globe-6dc-24h-202005170045-202005180045"
output_dir="roundtrip-predictrate-globe/"

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
