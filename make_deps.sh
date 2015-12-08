#!/bin/sh

rm -rf ./Godeps 
$GOPATH/bin/godep save ./...
