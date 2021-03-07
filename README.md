[![codecov](https://codecov.io/gh/confetti-framework/syslog/branch/main/graph/badge.svg?token=CZOT2GFEH4)](https://codecov.io/gh/confetti-framework/syslog)

# Syslog
Syslog package provides logger that generates syslog 
messages as defined in RFC 5424.

## Example Logger
```go
package main

import (
	"fmt"
	"github.com/confetti-framework/syslog"
	"os"
)

func main() {
	buf := &bytes.Buffer{}
	l := syslog.NewLogger(buf, syslog.USER, "hostname", "appName", "procid")

	// without structured data
	syslog.Info(l, "ImageUploaded", nil, "image uploaded by %s: %s", "username", "image.jpg")

	// with structured data
	sd := syslog.StructuredData{}
	sd.Element("id1").Set("par1", "val1")
	syslog.Error(l, "LoginFailed", sd, "login failed: %s", "username")

	fmt.Print(buf.String())

	// Output is similar to this:
	// <14>1 2017-08-15T23:13:15.335+02:00 hostname appName procid ImageUploaded - image uploaded by username: image.jpg
	// <11>1 2017-08-15T23:13:15.335+02:00 hostname appName procid LoginFailed [id1 par1="val1"] login failed: username
}	
```


## Example Writer
```go
package main

import (
	"github.com/confetti-framework/syslog"
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


