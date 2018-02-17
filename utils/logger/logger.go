package logger

import (
	"log"
	"bytes"
	"os"
)

var (
	buf		bytes.Buffer
	logger	*log.Logger
)

func buildPrefix(prefix string) string {
	return prefix + ": "
}

func Setup(prefix string) {
	logger = log.New(os.Stdout, buildPrefix(prefix), log.LstdFlags)
}

func Log(level string, message string, args ...interface{}) {
	var logs []interface{}
	logs = append(logs, message)

	if len(args) > 0 {
		logs = append(logs, args)
	}

	switch level {
	case "info":
		logger.Println(logs...)

	case "error":
		logger.Println(logs...)

	default:
		logger.Println("Not supported log level")

	}
}