#!/bin/sh

SERVICES='etcd mysql redis mongo registry'

echo "stopping all services..."
docker stop `echo $SERVICES`
docker rm `echo $SERVICES`

echo "starting etcd"
docker run -d --name etcd elcolio/etcd

echo "starting mysql"
docker run -d --name mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true mysql
docker run -it --link mysql --rm mysql sh -c 'exec mysql -hmysql -e"CREATE DATABASE PandoCloud"'

echo "starting redis"
docker run -d --name redis redis

echo "starting mongo"
docker run -d --name mongo mongo

echo "starting registry"
docker run -d --name registry --link etcd --link mysql pandocloud/pando-cloud registry -etcd http://etcd:2379 -rpchost internal:20034 -aeskey ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP -dbhost mysql -dbname PandoCloud -dbport 3306 -dbuser root -loglevel debug
