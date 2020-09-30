#!/usr/bin/python
import azure_exec
import common as c

## Exec latency zipf 0.75 varying added delays and percentiles
exp_dir_map = {
    'dt-95th-0ms'  : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add0ms/', \
    'dt-95th-1ms'  : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add1ms/', \
    'dt-95th-2ms'  : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add2ms/', \
    'dt-95th-4ms'  : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add4ms/', \
    'dt-95th-8ms'  : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic/', \
    'dt-95th-12ms' : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add12ms/', \
    'dt-95th-16ms' : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add16ms/', \
    'dt-95th-24ms' : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add24ms/', \
    'dt-95th-36ms' : '../exp-data/azure-exec-lat-globe-6dc-3r-zipf0.75/dynamic-add36ms/', \
  }

## Exec latency zipf 0.75 varying added delays
p_list = [
    'dt-95th-0ms'  , \
    'dt-95th-1ms'  , \
    'dt-95th-2ms'  , \
    'dt-95th-4ms'  , \
    'dt-95th-8ms'  , \
    'dt-95th-12ms' , \
    'dt-95th-16ms' , \
    'dt-95th-24ms' , \
    'dt-95th-36ms' , \
    ]
## Box figures
p_tick_label = [
    '0'  , \
    '1'  , \
    '2'  , \
    '4'  , \
    '8'  , \
    '12' , \
    '16' , \
    '24' , \
    '36' , \
    ]

# Vertical box figures
#x_label = ""
x_label = "Additional Delay (ms)"
output_file="azure-exec-lat-globe-6dc-3r-zipf0.75-adddelay.pdf"
azure_exec.lat_box(output_file, exp_dir_map, p_list, p_tick_label, x_label, whisker = [5, 95], ymin=0, exp_n=10)
