// Packe of common routines for the WebServices applications
package common

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

/*

   Logging tools

*/

import (
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

// Singleton pattern to return a common logging element.
// This also sets the 'Logger' in exported space and it is
// available from there, once it is setup.
func GetLogger(logDest string, debug bool) *log.Logger {
	// Setup only once
	if Logger == nil {
		var dest = os.Stdout
		if logDest != "-" {
			var err error
			dest, err = os.OpenFile(logDest, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to open log destination: %s\n", logDest)
				os.Exit(-1)
			}
		}
		Logger = log.New(dest, "logger: ", log.Lshortfile)
		Logger.SetPrefix("WebService: ")
		var flags int = log.Ldate | log.Ltime
		if debug {
			flags |= log.Lshortfile | log.Lmicroseconds
		}
		Logger.SetFlags(flags)
	}
	return Logger
}
