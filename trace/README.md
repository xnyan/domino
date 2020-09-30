This is the source code for collecting the inter-regaion latency on Azure.

# Build

cd $GOPATH/src/domino/trace/azure 

./build.sh

Binary files, client and server, will be generated.


# Collect new data traces

Folder $GOPATH/src/domino/trace/azure has the readme and scripts for making
network measurements on Azure.


# Generate figures based on the data traces

Folder $GOPATH/src/domino/trace/azure/fig has the readme and scripts showing
how to generate the figures about the inter-region latency in the Domino paper.

