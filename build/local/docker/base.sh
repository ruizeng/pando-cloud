#!/bin/sh

SERVICES='etcd mysql redis mongo'

echo "stopping all services..."
docker stop `echo $SERVICES`
docker rm `echo $SERVICES`

echo "starting etcd..."
docker run -d --name etcd elcolio/etcd

echo "starting mysql..."
docker run -d --name mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true mysql

echo "starting redis..."
docker run -d --name redis redis

echo "starting mongo..."
docker run -d --name mongo mongo
