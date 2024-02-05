package exec

import (
	"fmt"
	"log"
)

type Ilogger interface {
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})

	Errorln(...interface{})
	Warnln(...interface{})
	Infoln(...interface{})
	Debugln(...interface{})
}

type stdDebugLogger struct{}

// Errorf -
func (l stdDebugLogger) Errorf(format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("ERROR: %s\n", format), v...)
}

// Warnf -
func (l stdDebugLogger) Warnf(format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("WARN: %s\n", format), v...)
}

// Infof -
func (l stdDebugLogger) Infof(format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("INFO: %s\n", format), v...)
}

// Debugf -
func (l stdDebugLogger) Debugf(format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("DEBUG: %s\n", format), v...)
}

// Errorln -
func (l stdDebugLogger) Errorln(v ...interface{}) {
	l.Errorf("%s", v...)
}

// Warnln -
func (l stdDebugLogger) Warnln(v ...interface{}) {
	l.Warnf("%s", v...)
}

// Infoln -
func (l stdDebugLogger) Infoln(v ...interface{}) {
	l.Infof("%s", v...)
}

// Debugln -
func (l stdDebugLogger) Debugln(v ...interface{}) {
	l.Debugf("%s", v...)
}
