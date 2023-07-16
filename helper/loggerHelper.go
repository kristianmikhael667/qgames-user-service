package helper

import (
	"time"

	"github.com/ccpaging/log4go"
)

func Logger(logType string, message string, status string) {

	time := time.Now().Format("2006-01-02")
	log := log4go.NewLogger()

	flw := log4go.NewFileLogWriter("logs/go-"+"Mode_Dev "+time+".log", false)
	flw.SetFormat("[%D %T] [%L] %M")
	flw.SetRotate(false)
	flw.SetRotateSize(0)
	flw.SetRotateLines(0)
	flw.SetRotateDaily(false)
	log.AddFilter("file", log4go.FINE, flw)

	switch logType {
	case "info":
		log.Info(message)
	case "error":
		log.Error(message)
	}

	flw.Close()
	log.Close()
}
