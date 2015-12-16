docker#!/bin/sh

# restore dependency packages. 
cd $GOPATH/src/github.com/PandoCloud/pando-cloud
cp -r Godeps/_workspace/src/* $GOPATH/src

# install binaries
go install -v github.com/PandoCloud/pando-cloud/services/...
