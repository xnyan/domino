#!/usr/bin/python
import box as b
import common as c

src_dc = 'eastus2'
dst_dc = 'westus2'
dat_file='trace-azure-globe-6dc-24h-202005170045-202005180045/'+src_dc+'-'+dst_dc+'.log.txt'
output_file=src_dc+'-'+dst_dc+'-dist.pdf'

#time in seconds
start_t=12*60*60
length_t=60
window_t=1
slide_t=0.5

b.custom_slide_window_box_range(output_file, dat_file, start_t, length_t, window_t, slide_t, 's', 'ms', 'ms')
#b.slide_window_box_range(output_file, dat_file, start_t, length_t, window_t, slide_t, 's', 'ms', 'ms')
