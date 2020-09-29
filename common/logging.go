package common

import (
	"github.com/op/go-logging"
	"os"
)

/*
%{id}        Sequence number for log message (uint64).
%{pid}       Process id (int)
%{time}      Time when log occurred (time.Time)
%{level}     Log level (Level)
%{module}    Module (string)
%{program}   Basename of os.Args[0] (string)
%{message}   Message (string)
%{longfile}  Full file name and line number: /a/b/c/d.go:23
%{shortfile} Final file name element and line number: d.go:23
%{callpath}  Callpath like main.a.b.c...c  "..." meaning recursive call ~. meaning truncated path
%{color}     ANSI color based on log level
%{longpkg}   Full package path, eg. github.com/go-logging
%{shortpkg}  Base package path, eg. go-logging
%{longfunc}  Full function name, eg. littleEndian.PutUint32
%{shortfunc} Base function name, eg. PutUint32
%{callpath}  Call function path, eg. main.a.b.c
*/
var loggerFormat = logging.MustStringFormatter(
	`%{color}%{level:.4s} %{time:15:04:05.000} %{shortpkg} %{shortfile} %{shortfunc} >%{color:reset} %{message}`,
	//`%{color}%{time:15:04:05.000} %{shortpkg} %{shortfile} %{module} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func ConfigLogger(isDebug bool) {
	// Logger settings
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)

	logBackendFormat := logging.NewBackendFormatter(logBackend, loggerFormat)
	logBackendLevel := logging.AddModuleLevel(logBackendFormat)

	logBackendLevel.SetLevel(logging.ERROR, "")
	logBackendLevel.SetLevel(logging.CRITICAL, "")
	logBackendLevel.SetLevel(logging.INFO, "")

	if isDebug {
		logBackendLevel.SetLevel(logging.DEBUG, "")
	}

	logging.SetBackend(logBackendLevel)
}
