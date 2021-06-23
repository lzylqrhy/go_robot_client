package util

import (
	"github/go-robot/global"
	"log"
)

func DebugLog(format string, arg ...interface{}) {
	if (global.LogLevel & global.Debug) == global.Debug {
		log.Printf("D "+format, arg...)
	}
}

func WarningLog(format string, arg ...interface{}) {
	if (global.LogLevel & global.Warning) == global.Warning {
		log.Printf("W "+format, arg...)
	}
}

func ErrorLog(format string, arg ...interface{}) {
	if (global.LogLevel & global.Error) == global.Error {
		log.Panicf("E "+format, arg...)
	}
}

func FatalLog(format string, arg ...interface{}) {
	if (global.LogLevel & global.Fatal) == global.Fatal {
		log.Fatalf("F "+format, arg...)
	}
}