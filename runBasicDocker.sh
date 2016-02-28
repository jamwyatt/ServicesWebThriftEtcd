#!/bin/bash


ETCDNAME=jamwyatt_etcd
BACKENDNAME=jamwyatt_backend
FRONTENDNAME=jamwyatt_frontend

docker stop $ETCDNAME $BACKENDNAME $FRONTENDNAME >/dev/null 2>&1
docker rm $ETCDNAME  $BACKENDNAME $FRONTENDNAME >/dev/null 2>&1

if [ $# -eq 1 -a "X$1" = "Xstop" ]
then
	echo "Docker processes cleaned up ... done"
	exit 0
fi

# Etcd ... /data mapped to /tmp/etcd, exposes 4001 and 7001 (TCP)
# echo "Removing the old etcd data ... you might be asked for your root password as etcd builds these out as root"
# sudo rm -rf /tmp/etcd
# docker run -d --name=$ETCDNAME -v /tmp/etcd:/data  microbox/etcd --name defaultEtcdName --data-dir /data
docker run -d --name=$ETCDNAME microbox/etcd --name defaultEtcdName
etcdIP=`docker inspect --format='{{.NetworkSettings.IPAddress}}' $ETCDNAME`
echo "etcd listening on $etcdIP:4001"

# Points at the dynamic IP of etcd
docker run -d --name=$BACKENDNAME jamwyatt/backendprocessor -bi $etcdIP
backendIP=`docker inspect --format='{{.NetworkSettings.IPAddress}}' $BACKENDNAME`
echo "backend processor listening on $backendIP:8081"

# Points at the dynamic IP of the backend processor
docker run -d --name=$FRONTENDNAME -p 8080:8080  jamwyatt/webfrontend -bi $backendIP
frontendIP=`docker inspect --format='{{.NetworkSettings.IPAddress}}' $FRONTENDNAME`
echo "frontend processor listening on $frontendIP:8080"

docker ps -n 3 --no-trunc=true
