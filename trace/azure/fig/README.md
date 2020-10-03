# Generate the figures for the inter-region latency on Azure in the Domino paper

Follow the following the steps:

## Install the trace parser to $GOPATH/bin

cd $GOPATH/src/domino/trace/parser

go install

export PATH="$PATH:$GOROOT/bin:$GOPATH/bin" 

## Download and parse trace Files

cd $GOPATH/src/domino/trace/azure/fig

curl -JLO https://rgw.cs.uwaterloo.ca/BERNARD-domino/trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz ./

tar -xvzf trace-azure-globe-6dc-24h-202005170045-202005180045.tar.gz

The trace files will be under the folder trace-azure-globe-6dc-24h-202005170045-202005180045

./parse.sh trace-azure-globe-6dc-24h-202005170045-202005180045

## Generate the heatmap figures

cd $GOPAHT/src/domino/trace/azure/fig

./gen-fig-heatmap.py

The heatmap figures will be pdf files in the folder trace-azure-globe-6dc-24h-202005170045-202005180045

## Generate the box-and-whisker figure

cd $GOPAHT/src/domino/trace/azure/fig

./gen-fig-box.py

A pdf file eastus2-westus2-dist.pdf will be generated in current folder

## Generate the figure showing the correct prediction rate

cd $GOPAHT/src/domino/trace/azure/fig

./gen-predict-rate.py
(This would take some time to complete.)

./gen-fig-predict.py

A pdf file eastus2-westus2-predictrate.pdf will be generated in current folder
