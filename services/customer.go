// Implements the 'customer' web service
package services

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

import (
	"encoding/json"
	"../Thrift/gen-go/messages"
	"../common"
	"io"
	"net/http"
)

// The string base of the URI for which this service is registered
var root string
var addr string
var port int

// Set the root path ... optional
func InitService(r string, a string, p int) {
	root = r
	addr = a
	port = p
}

// customer GET processing
func CustomerGET(c *common.ThriftConnection, w http.ResponseWriter, r *http.Request) {
	common.Logger.Printf("GET: %s", r.URL)
	result, err := c.Conn().GetCustomer(r.URL.String()[len(root):])
	if err != nil {
		common.Logger.Printf("Failed to create new customer: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		common.Logger.Printf("Failed to convert to JSON: %s", result)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Customer POST processing
func CustomerPOST(c *common.ThriftConnection, w http.ResponseWriter, r *http.Request) {
	common.Logger.Printf("POST: %s", r.URL)
	if r.ContentLength <= 0 {
		common.Logger.Printf("Failed to detect body content in POST")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body = make([]byte, r.ContentLength)
	n, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		common.Logger.Printf("Failed to read body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	common.Logger.Printf("Body[%d]: %s", n, body)
	var thriftCust *messages.Customer = messages.NewCustomer()
	err = json.Unmarshal(body, &thriftCust)
	common.Logger.Printf("Customer: %s", thriftCust)
	result, err := c.Conn().CreateCustomer(thriftCust)
	if err != nil {
		common.Logger.Printf("Failed to create new customer: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	common.Logger.Printf("Customer Created: %s", result)
	bytes, err := json.Marshal(result)
	if err != nil {
		common.Logger.Printf("Failed to convert result to JSON: %s", result)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Handle registered hierarchy for HTTP methods (stored with 'SetRoot')
func Customer(w http.ResponseWriter, r *http.Request) {

	c, err := common.NewThriftConnection(addr, port)
	if err == nil {
		defer c.Close()
		switch r.Method {
		case "GET":
			CustomerGET(c, w, r)
		case "POST":
			CustomerPOST(c, w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
			common.Logger.Printf("Unsupported method: %s %s", r.Method, r.URL)
		}
	} else {
		common.Logger.Printf("Unable to make Thrift connection")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
