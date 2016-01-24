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
	"../common"
	"net/http"
)

// The string base of the URI for which this service is registered
var root string

// Set the root path ... optional
func SetRoot(r string) {
	root = r
}

// customer GET processing
func CustomerGET(w http.ResponseWriter, r *http.Request) {
	common.Logger.Printf("GET: %s", r.URL)
}

// Customer POST processing
func CustomerPOST(w http.ResponseWriter, r *http.Request) {
	common.Logger.Printf("POST: %s", r.URL)
}

// Handle registered hierarchy for HTTP methods (stored with 'SetRoot')
func Customer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		CustomerGET(w, r)
	case "POST":
		CustomerPOST(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		common.Logger.Printf("Unsupported method: %s %s", r.Method, r.URL)
	}
}
