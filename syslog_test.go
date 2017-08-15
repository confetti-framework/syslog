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
	wrappedBuf := syslog.NewWriter(buf, syslog.LOG_USER|syslog.LOG_NOTICE)
	logger := log.New(wrappedBuf, "", 0)
	logger.Println(msg)

	if !strings.HasSuffix(buf.String(), msg+"\n") {
		t.Fatalf("non expected msg suffix: %q", buf.String())
	}
}
