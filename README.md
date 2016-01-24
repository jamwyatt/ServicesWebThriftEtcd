# ServicesWebThriftEtcd
Basic demo of a web service Frontend with Thrift messaging to a Backend that uses Etcd for simple storage. This only uses HTTP as there is an expectation of SSL termination by a load balancer. This project is being written in golang with a specific purpose of being deployable in kubernetes.

Now has initial support for the basic 'POST' verb to the REST collection of /service/customer/. This takes HTTP POST with JSON body and converts it to a Thrift client request. The backend currently only records the receipt and returns.

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
  
