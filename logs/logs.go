package logs

import (
	"log"
	"os"
)

var Logger *log.Logger
var file *os.File

func Logs() {
	path, pathError := os.Getwd();
	if pathError != nil {
		log.Fatal(pathError)
	}

	logPath := path + "\\logs.log"

	var err error
	if file, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		log.Fatal(err)
	}

	Logger = log.New(file, "", log.Ltime)
}

func CloseLogs() {
	if file != nil {
		file.Close()
	}
}
