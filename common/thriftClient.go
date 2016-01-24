// Implements a thrift client connection

package common

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/jamwyatt/ServicesWebThriftEtcd/Thrift/gen-go/messages"
	"time"
)

/*
Copyright (C) 2015 J. Robert Wyatt

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

type ThriftConnection struct {
	addr       string
	port       int
	connection *messages.DataStoreClient
}

// Get the stored connection ... it must be initialized first!
func (c *ThriftConnection) Conn() *messages.DataStoreClient {
	if c.connection == nil {
		Logger.Printf("Thrift connection uninitialized")
	}
	return c.connection
}

func (c *ThriftConnection) Close() error {
	err := c.connection.Transport.Close()
	if err != nil {
		Logger.Printf("Failed to close Thrift connection: %s", err)
	}
	return err
}

func (c *ThriftConnection) Open() error {
	if c.connection != nil {
		Logger.Printf("WARNING re-opening thrigt client connection")
	}
	var transport thrift.TTransport
	var err error
	socket, err := thrift.NewTSocket(fmt.Sprintf("%s:%d", c.addr, c.port))
	if err != nil {
		Logger.Printf("Error opening socket: %s", err)
		return err
	}
	transportFactory := thrift.NewTTransportFactory()
	transport = transportFactory.GetTransport(socket)
	if err := transport.Open(); err != nil {
		Logger.Printf("Thrift client failed to connect to: %s:%d: %s", c.addr, c.port, err)
		return err
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	c.connection = messages.NewDataStoreClientFactory(transport, protocolFactory)
	Logger.Printf("Thrift connection: %s:%d connected", c.addr, c.port)
	return nil
}

func NewThriftConnection(addr string, port int) (*ThriftConnection, error) {
	c := ThriftConnection{addr, port, nil}
	var err error
	for i := 0; i < 5; i++ {
		err = c.Open() // Might log an error
		if err != nil {
			Logger.Printf("Thrift connection(%d): %s:%d failed: %s", i, c.addr, c.port, err)
			time.Sleep(time.Second * 1)
		} else {
			// Good connection
			break
		}
	}
	return &c, err
}
