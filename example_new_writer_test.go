// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog_test

import (
	"github.com/lanvard/syslog"
	"log"
	"os"
)

func ExampleNewWriter() {
	const msg = "Start HTTP server (addr=:8080)"

	hostname := "laptop"
	appName := "testapp"
	procid := "123"
	wrappedWriter := syslog.NewWriter(os.Stdout, syslog.USER|syslog.NOTICE, hostname, appName, procid)
	logger := log.New(wrappedWriter, "", 0)
	logger.Println(msg)

	// Output is similar to this:
	// <13>1 2017-08-15T23:13:15.335+02:00 laptop testapp 123 - - Start HTTP server (addr=:8080)
}
