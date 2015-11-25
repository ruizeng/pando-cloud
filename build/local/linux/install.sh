#!/bin/sh

cd $GOPATH/src/github.com/PandoCloud/pando-cloud
cp -r Godeps/_workspace/src/* $GOPATH/src
go install -v github.com/PandoCloud/pando-cloud/services/...
