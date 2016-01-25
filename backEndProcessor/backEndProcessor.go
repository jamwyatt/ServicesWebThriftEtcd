// Thrift and golang experiment
package main

/*
Thrift and golang experiment
Copyright (C) 2016 J. Robert Wyatt

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

/*

 */

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/jamwyatt/ServicesWebThriftEtcd/Thrift/gen-go/messages"
	"github.com/jamwyatt/ServicesWebThriftEtcd/common"
	"github.com/jamwyatt/etcdClientAPI/etcdMisc"
	"log"
	"net/http"
)

var logger *log.Logger

const customerRoot string = "/customers/"

type DataStoreHandler struct {
	etcdAddr   string
	etcdPort   int
	debug      bool
	connection etcdMisc.EtcdConnection
}

func NewDataStoreHandler(backendIp *string, backendPort *int, debug *bool) *DataStoreHandler {
	return &DataStoreHandler{*backendIp, *backendPort, *debug, etcdMisc.EtcdConnection{}}
}

func (p *DataStoreHandler) writeEtcdValue(path string, key string) int {
	_, err := p.connection.SetValue(path, key)
	if err != nil {
		logger.Printf("Failed to store %s in %s in etcd: %s", key, path, err)
		return 1
	}
	return 0
}

func (p *DataStoreHandler) makeConnection() {
	// Only make once
	if p.connection.Client != nil {
		return
	}
	var err error
	p.connection, err = etcdMisc.MakeEtcdConnection(&http.Client{Timeout: 0}, nil, "http", p.etcdAddr, p.etcdPort)
	if err != nil {
		logger.Printf("Failed to connect to etcd: %s", err)
	}
}

func (p *DataStoreHandler) GetCustomer(key string) (r *messages.Customer, err error) {
	logger.Printf("GetCustomer: %s", key)
	p.makeConnection() // Caching open
	path := customerRoot + key + "/" + "details"
	resp, err := p.connection.GetValue(path, false, false)
	if err != nil {
		logger.Printf("Failed to retrieve etcd[%s]: %s", path, err)
		return &messages.Customer{}, errors.New("Failures during customer creation")
	}

	var customer messages.Customer
	err = json.Unmarshal([]byte(resp.Node.Value), &customer)
	if err != nil {
		return &messages.Customer{}, errors.New("Failures during customer creation")
	}
	return &customer, nil
}

func (p *DataStoreHandler) CreateCustomer(customer *messages.Customer) (r *messages.Result_, err error) {
	uuid := uuid.NewUUID()
	logger.Printf("CreateCustomer: %s -> %s", uuid, customer)
	p.makeConnection() // Caching open
	path := customerRoot + uuid.String() + "/" + "details"
	bytes, err := json.Marshal(customer)
	if err != nil {
		return &messages.Result_{false, uuid.String(), "Failed to Create customer"}, errors.New("Failures during customer creation")
	}
	_, err = p.connection.SetValue(path, string(bytes))
	if err != nil {
		logger.Printf("Failed to store in etcd[%s]: %s", path, err)
		return &messages.Result_{false, uuid.String(), "Failed to Create customer"}, errors.New("Failures during customer creation")
	}

	return &messages.Result_{true, uuid.String(), ""}, nil
}

func (p *DataStoreHandler) GetAllCustomers(template *messages.Customer) (r map[string]*messages.Customer, err error) {
	return nil, nil
}

func (p *DataStoreHandler) UpdateCustomer(key string, customer *messages.Customer) (r *messages.Result_, err error) {
	logger.Printf("UpdateCustomer: %s", customer)
	return nil, nil
}

func runServer(listenIp *string, listenPort *int, backendIp *string, backendPort *int, debug *bool) error {
	transport, err := thrift.NewTServerSocket(fmt.Sprintf("%s:%d", *listenIp, *listenPort))
	if err != nil {
		return err
	}

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	handler := NewDataStoreHandler(backendIp, backendPort, debug)
	processor := messages.NewDataStoreProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	logger.Printf("Thrift Service configured")
	return server.Serve()
}

//
// Main args: -lp <port> -li <ip> -ep <port> -ei <ip> -log <file>
//      -lp <port>      Thrift Listening port
//      -li <ip>        Thrift Listening interface
//      -ep <port>      Etcd Service port
//      -ei <ip>        Etcd Service addr
//      -log <file>     file to log to, or '-' for stdout (default='-')
//      -debug		Enable debugging details
func main() {

	var listenPort = flag.Int("lp", 8081, "Thrift Listening port")
	var listenIp = flag.String("li", "0.0.0.0", "Listening interface")
	var backendPort = flag.Int("bp", 4001, "Etcd Service port")
	var backendIp = flag.String("bi", "127.0.0.1", "Etcd Service addr")
	var logDest = flag.String("log", "-", "Logging destination file, '-' for STDOUT")
	var debug = flag.Bool("debug", false, "enable debug")
	flag.Parse()

	logger = common.GetLogger(*logDest, *debug)
	logger.Print("Starting Backend Processor")
	logger.Printf("Thrift Listening: %s:%d\t\tEtcd Service: %s:%d\n", *listenIp, *listenPort, *backendIp, *backendPort)

	err := runServer(listenIp, listenPort, backendIp, backendPort, debug)
	if err != nil {
		logger.Printf("Failed to start Thrift service: %s\n", err)
	}
}
