// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog_test

import (
	"github.com/szxp/syslog"
	"log"
	"os"
)

func Example() {
	const msg = "Start HTTP server (addr=:8080)"

	wrappedWriter := syslog.NewWriter(os.Stdout, syslog.LOG_USER|syslog.LOG_NOTICE)
	logger := log.New(wrappedWriter, "", 0)
	logger.Println(msg)

	// Output is similar to this:
	// <13>1 2017-08-15T23:13:15.335+02:00 laptop /path/to/myprogram 21650 - - Start HTTP server (addr=:8080)
}
