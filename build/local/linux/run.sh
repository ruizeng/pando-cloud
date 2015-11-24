#!/bin/sh

# init db
echo 'CREATE DATABASE PandoCloud' | mysql -uroot 

# start services
$GOPATH/bin/httpaccess -etcd http://localhost:2379 -rpchost localhost:20035 -httphost :8000 -loglevel debug &
$GOPATH/bin/registry -etcd http://localhost:2379 -rpchost localhost:20034 -aeskey ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP -dbhost localhost -dbname PandoCloud -dbport 3306 -dbuser root -loglevel debug &
$GOPATH/bin/apiprovider -etcd http://localhost:2379 -httphost localhost:8888 &
$GOPATH/bin/devicemanager -etcd http://localhost:2379 -rpchost localhost:20033 &
$GOPATH/bin/controller -etcd http://localhost:2379 -rpchost localhost:20032 &

