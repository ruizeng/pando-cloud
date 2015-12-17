#!/bin/sh

SERVICES='registry devicemanager mqttaccess controller httpaccess apiprovider'

echo "stopping all services..."
docker stop `echo $SERVICES`
docker rm `echo $SERVICES`

echo "starting registry..."
docker run -d --name registry --link etcd --link mysql pandocloud/pando-cloud registry -etcd http://etcd:2379 -rpchost internal:20034 -aeskey ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP -dbhost mysql -dbname PandoCloud -dbport 3306 -dbuser root -loglevel debug

echo "starting devicemanager..."
docker run -d --name devicemanager --link etcd --link redis pandocloud/pando-cloud devicemanager -etcd http://etcd:2379 -rpchost internal:20033 -redishost redis:6379 -loglevel debug

echo "starting controller..."
docker run -d --name controller --link etcd --link mongo pandocloud/pando-cloud controller -etcd http://etcd:2379 -mongohost mongo -rpchost internal:20032 -loglevel debug

echo "starting mqttaccess..."
docker run -d --name mqttaccess -p 1883:1883 --link etcd pandocloud/pando-cloud mqttaccess -etcd http://etcd:2379 -tcphost :1883 -usetls -keyfile /go/src/github.com/PandoCloud/pando-cloud/pkg/server/testdata/key.pem -cafile /go/src/github.com/PandoCloud/pando-cloud/pkg/server/testdata/cert.pem -rpchost internal:20030 -loglevel debug

echo "starting httpaccess..."
docker run -d --name httpaccess -p 443:443 --link etcd --link redis pandocloud/pando-cloud httpaccess -etcd http://etcd:2379 -httphost :443 -redishost redis:6379 -usehttps -keyfile /go/src/github.com/PandoCloud/pando-cloud/pkg/server/testdata/key.pem -cafile /go/src/github.com/PandoCloud/pando-cloud/pkg/server/testdata/cert.pem -loglevel debug

echo "starting apiprovider..."
docker run -d --name apiprovider -p 8888:8888 --link etcd pandocloud/pando-cloud apiprovider -etcd http://etcd:2379 -httphost :8888 -loglevel debug