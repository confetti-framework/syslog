// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog_test

import (
	"bytes"
	"github.com/szxp/syslog"
	"log"
	"strings"
	"testing"
)

func TestWriter(t *testing.T) {
	const msg = "this is the message details"

	buf := &bytes.Buffer{}
	wrappedBuf := syslog.NewWriter(buf, syslog.USER|syslog.NOTICE)
	logger := log.New(wrappedBuf, "", 0)
	logger.Println(msg)

	if !strings.HasSuffix(buf.String(), msg+"\n") {
		t.Fatalf("non expected msg suffix: %q", buf.String())
	}
}

func TestStructuredData(t *testing.T) {
	sd := syslog.SDElem("id1").
		Data("par1", "\"val1\"").
		Data("par2", "val2").
		SDElem("id2").
		Data("par1", "val1").
		Data("par2", "val2")

	expected := `[id1 par1="\"val1\"" par2="val2"][id2 par1="val1" par2="val2"]`
	if sd.String() != expected {
		t.Fatalf("got %v, but expected %v", sd.String(), expected)
	}
}
