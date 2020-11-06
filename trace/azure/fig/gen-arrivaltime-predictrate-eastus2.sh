dc_l=(
#eastus2
westus2
#francecentral
#australiaeast
#eastasia
#southeastasia
)

data_dir="trace-azure-globe-6dc-24h-202005170045-202005180045"
output_dir="arrivaltime-globe"

mkdir -p ${output_dir}

src_dc="eastus2"
pth="95.0"
window="1000"

date
for dst_dc in ${dc_l[@]}
do
  if [ "${src_dc}" == "${dst_dc}" ]; then
    continue
  fi
  echo ${src_dc} ${dst_dc} ${data_dir} ${output_dir}
  dat_f="${src_dc}-${dst_dc}.log.txt"
  for pth in 5 10 15 20 25 30 35 40 45 50 55 60 65 70 75 80 85 90 95 99
  do
  for window in 100 200 400 600 800 1000
  do
    ./gen-arrivaltime.py ${data_dir} ${dat_f} ${output_dir} ${pth} ${window}
  done # window
  done # pth
done
date
