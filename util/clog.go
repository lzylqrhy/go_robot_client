package util

import "log"

const (
	Debug = 1 << iota
	Warning
	Error
)

func DebugLog(format string, arg ...interface{}) {
	log.Printf(format, arg)
}

func WarningLog(format string, arg ...interface{}) {
	log.Printf(format, arg)
}

func ErrorLog(format string, arg ...interface{}) {
	log.Panicf(format, arg)
}