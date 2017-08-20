// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog_test

import (
	"bytes"
	"github.com/szxp/syslog"
	"log"
	"reflect"
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
