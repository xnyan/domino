# Generate the figures for the inter-region latency on Azure in the Domino paper

Follow the following the steps:

## Download trace Files

cd $GOPATH/src/domino/trace/azure/fig

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz

tar -xvzf trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz

The trace files will be under the folder trace-azure-globe-6dc-24h-202005170045-202005180045

## Generate the heatmap figures about network roundtrip delays

cd $GOPATH/src/domino/trace/azure/fig

./gen-fig-heatmap.py

The heatmap figures will be pdf files in the folder trace-azure-globe-6dc-24h-202005170045-202005180045

## Generate the box-and-whisker figure about network roundtrip delays

cd $GOPATH/src/domino/trace/azure/fig

./gen-fig-box.py

A pdf file eastus2-westus2-dist.pdf will be generated in the current folder

## Generate the figure about the correct prediction rate for request arrival time at replicas

cd $GOPATH/src/domino/trace/azure/fig

./gen-arrivaltime-predictrate-eastus2.sh
(This would take some time to complete.)

./gen-fig-predict-arrivaltime.py

A pdf file eastus2-westus2-arrivaltime-predictrate.pdf will be generated in the current folder
