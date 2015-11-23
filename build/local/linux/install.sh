#!/bin/sh

cd $GOPATH/src/github.com/PandoCloud/pando-cloud
go get github.com/tools/godep
$GOPATH/bin/godep restore
go get ./...
