// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package syslog generates syslog messages.
package syslog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// The Priority is a combination of the syslog facility and
// severity. For example, LOG_USER | LOG_NOTICE.
type Priority int

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to LOG_FTP on Linux, BSD, and OS X.
	LOG_KERN Priority = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

const version = 1 // defined in RFC 5424.

// Writer generates syslog messages as defined in RFC 5424.
type writer struct {
	out io.Writer
	pri Priority
}

// NewWriter wrappes another io.Writer and returns a new
// io.Writer that generates syslog messages as defined
// in RFC 5424 and writes them to the given io.Writer.
func NewWriter(out io.Writer, pri Priority) io.Writer {
	if pri < 0 || pri > LOG_LOCAL7|LOG_DEBUG {
		panic("syslog: invalid priority: " + strconv.Itoa(int(pri)))
	}

	return &writer{
		out: out,
		pri: pri,
	}
}

// Write generates and writes a syslog message to the
// underlying io.Writer.
func (w *writer) Write(d []byte) (n int, err error) {
	if len(d) == 0 {
		return 0, nil
	}

	if d[0] != '<' {
		return w.out.Write(w.format(d))
	}

	// don't format a syslog message
	return w.out.Write(d)
}

const rfc3339Milli = "2006-01-02T15:04:05.999-07:00"

func (w *writer) format(d []byte) []byte {
	timestamp := time.Now().Format(rfc3339Milli)
	hostname, _ := os.Hostname()
	appName := os.Args[0]
	procid := os.Getpid()

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "<%d>%d %s %s %s %d - - ",
		w.pri,
		version,
		timestamp,
		hostname,
		appName,
		procid,
	)
	buf.Write(d)

	if d[len(d)-1] != '\n' {
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
