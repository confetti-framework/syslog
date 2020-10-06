// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package syslog provides logger that generates syslog
// messages as defined in RFC 5424.
package syslog

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

// The Priority is a combination of the syslog facility and
// severity. For example, USER | NOTICE.
type Priority int

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	EMERG Priority = iota
	ALERT
	CRIT
	ERR
	WARNING
	NOTICE
	INFO
	DEBUG
)

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to FTP on Linux, BSD, and OS X.
	KERN Priority = iota << 3
	USER
	MAIL
	DAEMON
	AUTH
	SYSLOG
	LPR
	NEWS
	UUCP
	CRON
	AUTHPRIV
	FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOCAL0
	LOCAL1
	LOCAL2
	LOCAL3
	LOCAL4
	LOCAL5
	LOCAL6
	LOCAL7
)

const version = 1 // defined in RFC 5424.

// NewWriter wrappes another io.Writer and returns a new
// io.Writer that generates syslog messages as defined
// in RFC 5424 and writes them to the given io.Writer.
// The returned io.Writer is NOT safe for concurrent use
// by multiple goroutines.
func NewWriter(out io.Writer, pri Priority, hostname, appName, procid string) io.Writer {
	return &writer{
		out,
		pri,
		hostname,
		appName,
		procid,
	}
}

// Writer generates syslog messages as defined in RFC 5424.
type writer struct {
	out      io.Writer
	pri      Priority
	hostname string
	appName  string
	procid   string
}

var nl = []byte{'\n'}

// Write generates and writes a syslog message to the
// underlying io.Writer.
func (w *writer) Write(d []byte) (int, error) {
	if len(d) == 0 {
		return 0, nil
	}

	// don't format a syslog message
	if d[0] != '<' {
		d = formatSyslog(
			w.pri,
			time.Now(),
			"",
			w.hostname,
			w.appName,
			w.procid,
			"",
			nil,
			d)
	}

	n, err := w.out.Write(d)
	if d[len(d)-1] != '\n' && err == nil {
		w.out.Write(nl)
	}
	return n, err
}

const rfc3339Milli = "2006-01-02T15:04:05.999-07:00"

func formatSyslog(
	pri Priority,
	timestamp time.Time,
	timeFormat string,
	hostname string,
	appName string,
	procid string,
	msgid string,
	structData StructuredData,
	msg []byte,
) []byte {
	if timeFormat == "" {
		timeFormat = rfc3339Milli
	}

	ts := timestamp.Format(timeFormat)
	hostname = defaultIfEmpty(hostname, "-")
	appName = defaultIfEmpty(appName, "-")
	procid = defaultIfEmpty(procid, "-")
	msgid = defaultIfEmpty(msgid, "-")

	sd := ""
	if structData != nil {
		sd = structData.String()
	}
	sd = defaultIfEmpty(sd, "-")

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "<%d>%d %s %s %s %s %s %s ",
		pri,
		version,
		ts,
		hostname,
		appName,
		procid,
		msgid,
		sd,
	)
	buf.Write(msg)

	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func defaultIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// Logger generates syslog messages.
type Logger interface {

	// Log generates a syslog message.
	Log(severity Priority, msgId string, sd StructuredData, msgFormat string, a ...interface{})
}

// NewLogger returns a new syslog logger that writes to
// the specified io.Writer.
// The returned Logger is safe for concurrent use by
// multiple goroutines.
func NewLogger(w io.Writer, facility Priority, hostname, appName, procid string) Logger {
	return &logger{
		sync.Mutex{},
		w,
		facility,
		hostname,
		appName,
		procid,
	}
}

type logger struct {
	mu       sync.Mutex
	w        io.Writer
	facility Priority
	hostname string
	appName  string
	procid   string
}

func (l *logger) Log(severity Priority, msgId string, sd StructuredData, msgFormat string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(msgFormat, a...)
	l.w.Write(formatSyslog(
		Priority(l.facility|severity),
		time.Now(),
		"",
		l.hostname,
		l.appName,
		l.procid,
		msgId,
		sd,
		[]byte(msg)))
}

// StructuredData provides a mechanism to express information in a well
// defined, easily parseable and interpretable data format. There are
// multiple usage scenarios.  For example, it may express meta-
// information about the syslog message or application-specific
// information such as traffic counters or IP addresses.
//
// StructuredData can contain zero, one, or multiple structured data
// elements, which are referred to as SDElement.
type StructuredData map[string]SDElement

// Element returns an SDElement associated with the given id.
// If an element with the id does not exist a new SDElement
// will be created.
func (d StructuredData) Element(id string) SDElement {
	elem, ok := d[id]
	if !ok {
		elem = make(SDElement, 1)
		d[id] = elem
	}
	return elem
}

// Ids returns the ids of the SDElements in lexicographical order.
func (d StructuredData) Ids() []string {
	ids := make([]string, 0, len(d))
	for id := range d {
		if len(d[id]) > 0 {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

// Strings returns the string representation of the structured data.
func (d StructuredData) String() string {
	r := strings.NewReplacer(`"`, `\"`, `\`, `\\`, `]`, `\]`)
	buf := &bytes.Buffer{}
	for _, id := range d.Ids() {
		elem := d[id]
		if len(elem) > 0 {
			buf.WriteByte('[')
			buf.WriteString(id)
			for _, name := range elem.Names() {
				buf.WriteByte(' ')
				fmt.Fprintf(buf, `%s="%s"`, name, r.Replace(elem[name]))
			}
			buf.WriteByte(']')
		}
	}
	return buf.String()
}

// SDElement represents a structured data element and consists
// name-value pairs.
type SDElement map[string]string

// Set sets a value associated with the specified name.
func (e SDElement) Set(name, value string) SDElement {
	e[name] = value
	return e
}

// Get returns a value associated with the specified name.
func (e SDElement) Get(name string) string {
	value, ok := e[name]
	if !ok {
		return ""
	}
	return value
}

// Names returns the parameter names in lexicographical order.
func (e SDElement) Names() []string {
	names := make([]string, 0, len(e))
	for name := range e {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func Emergency(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(EMERG, msgId, sd, format, a...)
}

func Critical(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(CRIT, msgId, sd, format, a...)
}

func Alert(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(ALERT, msgId, sd, format, a...)
}

func Error(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(ERR, msgId, sd, format, a...)
}

func Warning(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(WARNING, msgId, sd, format, a...)
}

func Notice(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(NOTICE, msgId, sd, format, a...)
}

func Info(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(INFO, msgId, sd, format, a...)
}

func Debug(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	if l == nil {
		return
	}
	l.Log(DEBUG, msgId, sd, format, a...)
}
