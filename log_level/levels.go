package log_level

import "github.com/lanvard/syslog"

type Level = syslog.Priority // severity

const (
	// Level/severity
	// These are the same on Linux, BSD, and OS X.
	EMERGENCY syslog.Priority = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)
