// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog_test

import (
	"bytes"
	"github.com/lanvard/syslog"
	"github.com/lanvard/syslog/log_level"
	"log"
	"reflect"
	"strings"
	"testing"
)

func Test_writer(t *testing.T) {
	const msg = "this is the message details"

	buf := &bytes.Buffer{}
	hostname := "laptop"
	appName := "testapp"
	procid := "123"
	wrappedBuf := syslog.NewWriter(buf, syslog.USER|log_level.NOTICE, hostname, appName, procid)
	logger := log.New(wrappedBuf, "", 0)
	logger.Println(msg)

	expectedPrefix := "<13>1"
	if !strings.HasPrefix(buf.String(), expectedPrefix) {
		t.Fatalf("non-expected prefix: %s", buf.String())
	}

	if !strings.HasSuffix(buf.String(), msg+"\n") {
		t.Fatalf("non-expected msg suffix: %s", buf.String())
	}
}

func Test_logger(t *testing.T) {
	buf := &bytes.Buffer{}
	l := syslog.NewLogger(buf, syslog.USER, "hostname", "appName", "procid")

	sd := syslog.StructuredData{}
	sd.Element("id1").Set("par1", "val1")
	l.Log(log_level.ERROR, "LoginFailed", sd, "login failed: %s", "username")

	expectedPrefix := "<11>1"
	if !strings.HasPrefix(buf.String(), expectedPrefix) {
		t.Fatalf("non-expected prefix: %s", buf.String())
	}

	expectedSuffix := "appName procid LoginFailed [id1 par1=\"val1\"] login failed: username\n"
	if !strings.HasSuffix(buf.String(), expectedSuffix) {
		t.Fatalf("non-expected suffix: %s", buf.String())
	}
}

func Test_structured_data(t *testing.T) {
	sd := syslog.StructuredData{}
	sd.Element("id1").
		Set("par1", "\"val1\"").
		Set("par2", "val2")
	sd.Element("id2").
		Set("par1", "val1").
		Set("par2", "val2")

	expectedIds := []string{"id1", "id2"}
	ids := sd.Ids()
	if !reflect.DeepEqual(ids, expectedIds) {
		t.Fatalf("got ids: %v, but expected: %v", ids, expectedIds)
	}

	expectedValue := "val2"
	value := sd.Element("id2").Get("par2")
	if value != expectedValue {
		t.Fatalf("got value: %v, but expected: %v", value, expectedValue)
	}

	expectedString := `[id1 par1="\"val1\"" par2="val2"][id2 par1="val1" par2="val2"]`
	if sd.String() != expectedString {
		t.Fatalf("got string: %v, but expected: %v", sd.String(), expectedString)
	}
}
