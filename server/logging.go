package main

import (
	"io"
	"log"
)

var (
	// Trace logs
	Trace *log.Logger

	// Info logs
	Info *log.Logger

	//Warning logs
	Warning *log.Logger

	// Error logs
	Error *log.Logger
)

// InitLogging initializes logging
func InitLogging(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
