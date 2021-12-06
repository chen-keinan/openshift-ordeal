package logger

import (
	"log"
)

//OsOrdealLogger Object
type OsOrdealLogger struct {
}

//GetLog return native logger
func GetLog() *OsOrdealLogger {
	return &OsOrdealLogger{}
}

//Console print to console
func (BLogger *OsOrdealLogger) Console(str string) {
	log.SetFlags(0)
	log.Print(str)
}

//Table print to console
func (BLogger *OsOrdealLogger) Table(v interface{}) {
	log.SetFlags(0)
	log.Print(v)
}
