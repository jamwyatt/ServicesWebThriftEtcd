// Front end WebService with Thrift connection to backend
package main

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
	"flag"
	"fmt"
	"../common"
	"../services"
	"log"
	"net/http"
	"time"
)

var logger *log.Logger

type service struct {
	url  string
	f    func(http.ResponseWriter, *http.Request)
	init func(root string, beAddr string, bePort int) // Service init with services root and Thrift server details
}

// List of web end point services
var myServices []service = []service{
	{"customer", services.Customer, services.InitService}, // customer.go
}

func initWebServices(root string, addr string, port int) {
	logger.Printf("Web Service root: %s\n", root)
	for _, v := range myServices {
		var path string = fmt.Sprintf("/%s/%s/", root, v.url)
		logger.Printf("Registering: %s", path)
		http.HandleFunc(path, v.f)
		if v.init != nil {
			v.init(path, addr, port)
		}
	}
}

//
// Main args: -lp <port> -li <ip> -bp <port> -bi <ip> -log <file>
//	-lp <port>	Listening port (default=8080)
//	-li <ip>	Listening interface (default=0.0.0.0)
//	-bp <port>	Backend Service port (default=8081)
//	-bi <ip>	Backend Service addr (default=127.0.0.1)
//	-log <file>	file to log to, or '-' for stdout (default='-')
//	-root <string>	uri root service (default='service')
func main() {
	var listenPort = flag.Int("lp", 8080, "Listening port")
	var listenIp = flag.String("li", "0.0.0.0", "Listening interface")
	var backendPort = flag.Int("bp", 8081, "Backend Service port")
	var backendIp = flag.String("bi", "127.0.0.1", "Backend Service addr")
	var logDest = flag.String("log", "-", "Logging destination file, '-' for STDOUT")
	var serviceRoot = flag.String("root", "service", "web service root")
	var debug = flag.Bool("debug", false, "enable debug")
	flag.Parse()

	logger = common.GetLogger(*logDest, *debug)
	logger.Print("Starting WebService")
	logger.Printf("Listening: %s:%d\t\tBackend Service: %s:%d\n", *listenIp, *listenPort, *backendIp, *backendPort)

	// Start the webservice
	initWebServices(*serviceRoot, *backendIp, *backendPort)
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", *listenIp, *listenPort),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 10 * 1024 * 1024,
	}
	logger.Fatal(s.ListenAndServe())
}
