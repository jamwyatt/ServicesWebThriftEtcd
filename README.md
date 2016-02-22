# ServicesWebThriftEtcd
Basic demo of a web service Frontend with Thrift messaging to a Backend that uses Etcd for simple storage.
This only uses HTTP as there is an expectation of SSL termination by a load balancer.
This project is being written in golang with a specific purpose of being deployable in kubernetes.

Now has initial support for the basic 'POST' verb to the REST collection of /service/customer/.
This takes HTTP POST with JSON body and converts it to a Thrift client request. The backend currently only records the receipt and returns. POST/GET for the customer collection now working.

Next work items, link in the EtcClient library and start building the new client on POST/Create. Then start supporting end-to-end query.

Both the frontend and backend take two sets of IP:port details. The frontend listens for HTTP connections and connects to the backend. The backend listens for Thrift connections and connects to Etcd.

```
$ ./webService -h
Usage of ./webService:
  -bi string
    	Backend Service addr (default "127.0.0.1")
  -bp int
    	Backend Service port (default 8081)
  -debug
    	enable debug
  -li string
    	Listening interface (default "0.0.0.0")
  -log string
    	Logging destination file, '-' for STDOUT (default "-")
  -lp int
    	Listening port (default 8080)
  -root string
    	web service root (default "service")
```

```
$ ./backEndProcessor -h
Usage of ./backEndProcessor:
  -bi string
    	Etcd Service addr (default "127.0.0.1")
  -bp int
    	Etcd Service port (default 4000)
  -debug
    	enable debug
  -li string
    	Listening interface (default "0.0.0.0")
  -log string
    	Logging destination file, '-' for STDOUT (default "-")
  -lp int
    	Thrift Listening port (default 8081)
```  
  
#Notes on dependancies:
  
Thrift go libraries
```
  go get git.apache.org/thrift.git/lib/go/thrift/...
```
UUID code from google
```
  go get github.com/pborman/uuid
```
Etcd library
```
  go get github.com/jamwyatt/etcdClientAPI
```

#Example Build

Building with 'make all' generates Thrift messages and builds two static golang binaries (backEndProcessor and webFrontEnd). Each binary is placed into a Docker container. The base containers are build off of the 'scratch' base so that only the binary files are present.

```
$ make all
make -C Thrift all
make[1]: Entering directory '/home/jamwyatt/gowork/src/ServicesWebThriftEtcd/Thrift'
thrift -gen go messages.thrift
make[1]: Leaving directory '/home/jamwyatt/gowork/src/ServicesWebThriftEtcd/Thrift'
make -C backEndProcessor all
make[1]: Entering directory '/home/jamwyatt/gowork/src/ServicesWebThriftEtcd/backEndProcessor'
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo backEndProcessor.go
docker build --rm=true -t jamwyatt/backendprocessor:latest .
Sending build context to Docker daemon 4.808 MB
Sending build context to Docker daemon 
Step 0 : FROM scratch
 ---> 
Step 1 : MAINTAINER J. Robert Wyatt
 ---> Running in 0a61b05382e4
 ---> 57b5a91f0ba0
Removing intermediate container 0a61b05382e4
Step 2 : WORKDIR /
 ---> Running in ff693bb05d8a
 ---> 6786f56f80c3
Removing intermediate container ff693bb05d8a
Step 3 : ADD backEndProcessor /
 ---> a3a91ceaca15
Removing intermediate container b980fc20029b
Step 4 : ENTRYPOINT /backEndProcessor
 ---> Running in 459073fd4d5b
 ---> 9bd0a1b0158c
Removing intermediate container 459073fd4d5b
Step 5 : EXPOSE 8081
 ---> Running in 0af16e2e5cfe
 ---> be37a0ffb1a1
Removing intermediate container 0af16e2e5cfe
Successfully built be37a0ffb1a1
touch docker
make[1]: Leaving directory '/home/jamwyatt/gowork/src/ServicesWebThriftEtcd/backEndProcessor'
make -C webService all
make[1]: Entering directory '/home/jamwyatt/gowork/src/ServicesWebThriftEtcd/webService'
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo webFrontEnd.go
docker build --rm=true -t jamwyatt/webfrontend:latest .
Sending build context to Docker daemon 4.645 MB
Sending build context to Docker daemon 
Step 0 : FROM scratch
 ---> 
Step 1 : MAINTAINER J. Robert Wyatt
 ---> Using cache
 ---> 57b5a91f0ba0
Step 2 : WORKDIR /
 ---> Using cache
 ---> 6786f56f80c3
Step 3 : ADD webFrontEnd /
 ---> 4ce07453882d
Removing intermediate container c6aea6c8b1ed
Step 4 : ENTRYPOINT /webFrontEnd
 ---> Running in 7ee2eb1fd9a5
 ---> d1ae5cf360b1
Removing intermediate container 7ee2eb1fd9a5
Step 5 : EXPOSE 8080
 ---> Running in 9ca0c4fc6118
 ---> 089d87b089f4
Removing intermediate container 9ca0c4fc6118
Successfully built 089d87b089f4
touch docker
make[1]: Leaving directory '/home/jamwyatt/gowork/src/ServicesWebThriftEtcd/webService'
$ 
```
#Example run

First this is first, use the script to start an 'etcd' container and the two local containers. The script ensures that only the front end container publishes its exposed port. In terms of docker, this is a very basic setup and it records and uses the individual docker IP addresses for literal connections between the layers. This would be much better handled as services in kubernetes using something like skydns to propagate the dns service names (that's TBD for next time).

```
$ ./runBasicDocker.sh 
Removing the old etcd data ... you might be asked for your root password as etcd builds these out as root
[sudo] password for jamwyatt: 
b4e6ee9e93b8bb0b2460751acd42ae16beda5d6c1cb3a19ea5ddda1cfb52529b
etcd listening on 172.17.0.83:4001
b728181d93149cde4ff0bb33d1cb3790eb7b554b601e70254c61e3efbc8e72ca
backend processor listening on 172.17.0.84:8081
c023c5d2ad9c1dd45e713a62ffe7c27e12440cf956c8455af0267c1914fe6cf0
frontend processor listening on 172.17.0.85:8080
CONTAINER ID                                                       IMAGE                              COMMAND                                               CREATED                  STATUS                  PORTS                    NAMES
c023c5d2ad9c1dd45e713a62ffe7c27e12440cf956c8455af0267c1914fe6cf0   jamwyatt/webfrontend:latest        "/webFrontEnd -bi 172.17.0.84"                        Less than a second ago   Up Less than a second   0.0.0.0:8080->8080/tcp   jamwyatt_frontend   
b728181d93149cde4ff0bb33d1cb3790eb7b554b601e70254c61e3efbc8e72ca   jamwyatt/backendprocessor:latest   "/backEndProcessor -bi 172.17.0.83"                   1 seconds ago            Up Less than a second   8081/tcp                 jamwyatt_backend    
b4e6ee9e93b8bb0b2460751acd42ae16beda5d6c1cb3a19ea5ddda1cfb52529b   microbox/etcd:latest               "/bin/etcd --name defaultEtcdName --data-dir /data"   1 seconds ago            Up Less than a second   4001/tcp, 7001/tcp       jamwyatt_etcd       
$
```

Once this is run, the front end is listening on all interfaces of the host, on port 8080.

#Example usage

Here's an example using 'curl' to add and query a customer

```
$ cat newCust.json 
{
    "FirstName":"Robert",
    "LastName":"Wyatt",
    "Addr1":"123 Somewhere St",
    "Addr2":"",
    "City":"Sometown",
    "StateProvince":"XY",
    "PostalZip":"123456",
    "Country":"USA"
}
Roberts-iMac:tmp jamwyatt$ curl -s -o - -H 'Content-Type: application/json' -X POST http://192.168.1.147:8080/service/customer/ -d @newCust.json
{"Ok":true,"Key":"9b0269cb-d908-11e5-8d71-0242ac11004b","Err":""}
$ curl -s -o - http://192.168.1.147:8080/service/customer/9b0269cb-d908-11e5-8d71-0242ac11004b | python -m json.tool
{
    "Addr1": "123 Somewhere St",
    "Addr2": "",
    "City": "Sometown",
    "Country": "USA",
    "FirstName": "Robert",
    "LastName": "Wyatt",
    "PostalZip": "123456",
    "StateProvince": "XY"
}

```

