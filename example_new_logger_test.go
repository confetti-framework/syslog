// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog_test

import (
	"bytes"
	"fmt"
	"github.com/szxp/syslog"
)

func ExampleNewLogger() {
	buf := &bytes.Buffer{}
	l := syslog.NewLogger(buf, syslog.USER, "hostname", "appName", "procid")

	// without structured data
	l.Log(syslog.INFO, "ImageUploaded", nil, "image uploaded by %s: %s", "username", "image.jpg")

	// with structured data
	sd := syslog.StructuredData{}
	sd.Element("id1").Set("par1", "val1")
	l.Log(syslog.ERR, "LoginFailed", sd, "login failed: %s", "username")

	fmt.Print(buf.String())

	// Output is similar to this:
	// <14>1 2017-08-15T23:13:15.335+02:00 hostname appName procid ImageUploaded - image uploaded by username: image.jpg
	// <11>1 2017-08-15T23:13:15.335+02:00 hostname appName procid LoginFailed [id1 par1="val1"] login failed: username
}
