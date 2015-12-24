FROM golang:1.5.2

MAINTAINER razr razr.china@gmail.com

ADD . /go/src/github.com/PandoCloud/pando-cloud

RUN go get -u github.com/PandoCloud/pando-cloud/...