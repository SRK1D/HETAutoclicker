package logger

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

var logger *log.Logger

/*
-	 Function Name: initialise
-	 Description: Initialises the global variable logger to a logging instance
*/
func initialise() {
	if logger != nil {
		return
	}

	logger = log.New()
	logger.SetFormatter(&log.TextFormatter{ForceColors: true})
	logger.SetOutput(os.Stdout)
}

/*
-	Function Name: QuickLog
-	Description: Quickly logs the message and the field
*/
func QuickLog(logType LogType, message string, fields ...log.Fields) {

	// Defines a variable with multiple arguments of type any
	var logFunc func(...interface{})
	var entry *log.Entry

	initialise()

	if len(fields) > 0 {
		entry = logger.WithFields(fields[0])
	} else {
		entry = logger.WithFields(log.Fields{})
	}

	// Assigns the function variable to the valid type
	switch logType {
	case TC_INFO:
		logFunc = entry.Info
		break
	case TC_WARN:
		logFunc = entry.Warn
		break
	case TC_ERROR:
		logFunc = entry.Error
		break
	case TC_FATAL:
		logFunc = entry.Fatal
		break
	case TC_PANIC:
		logFunc = entry.Panic
		break
	default:
		panic(errors.New("logging -> Log -> Invalid log type"))
	}

	logFunc(message)
}
