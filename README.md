[![Build Status](https://travis-ci.org/szxp/syslog.svg?branch=master)](https://travis-ci.org/szxp/syslog)
[![Build Status](https://ci.appveyor.com/api/projects/status/github/szxp/syslog?branch=master&svg=true)](https://ci.appveyor.com/project/szxp/syslog)
[![GoDoc](https://godoc.org/github.com/szxp/syslog?status.svg)](https://godoc.org/github.com/szxp/syslog)
[![Go Report Card](https://goreportcard.com/badge/github.com/szxp/syslog)](https://goreportcard.com/report/github.com/szxp/syslog)

# syslog
Syslog package provides logger that generates syslog 
messages as defined in RFC 5424.


## Example Writer
```go
package main

import (
	"github.com/szxp/syslog"
	"log"
	"os"
)

func main() {
	const msg = "Start HTTP server (addr=:8080)"

	hostname := "laptop"
	appName := "testapp"
	procid := "123"
	wrappedWriter := syslog.NewWriter(os.Stdout, syslog.USER|syslog.NOTICE, hostname, appName, procid)
	logger := log.New(wrappedWriter, "", 0)
	logger.Println(msg)

	// Output is similar to this:
	// <13>1 2017-08-15T23:13:15.33+02:00 laptop testapp 123 - - Start HTTP server (addr=:8080)
}
```


## Example Logger
```go
package main

import (
	"fmt"
	"github.com/szxp/syslog"
	"os"
)

func main() {
	buf := &bytes.Buffer{}
	l := syslog.NewLogger(buf, "hostname", "appName", "procid")

	sd := syslog.StructuredData{}
	sd.Element("id1").Set("par1", "val1")
	l.Log(syslog.USER|syslog.ERR, "LoginFailed", sd, "login failed: %s", "username")

	fmt.Print(buf.String())

	// Output is similar to this:
	// <11>1 2017-08-15T23:13:15.335+02:00 hostname appName procid LoginFailed [id1 par1="val1"] login failed: username
}	
```

