#!/bin/bash -e
# WARNING: FOR PIPELINE USE UNLESS YOU WANT TO INSTALL GO
echo "For pipeline use!"

export PATH=$PATH:$PWD/goinstall/go/bin
export GOPATH=$PWD/goinstall/workspace
export PATH=$PATH:$GOPATH/bin

PLUGINDIR=$GOPATH/src/git.ng.bluemix.net/csc/pipeline-tools/deploySpace
mkdir -p $PLUGINDIR

mv deploySpace/* $PLUGINDIR

cd $PLUGINDIR
go get
go test
go install
