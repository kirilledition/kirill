package cmd

import (
	"io"
	"log"
	"os"
	"strings"
)

var logger *log.Logger

func getLogger(filename string) (*log.Logger, *os.File, error) {
	logFilename := filename + ".log"
	logFile, err := os.Create(logFilename)
	if err != nil {
		return nil, nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	loggerInstance := log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime)
	loggerInstance.Printf("kirill - Version: %s, Git Commit: %s, Logging to file %s", GitTag, GitCommit[:7], logFilename)

	return loggerInstance, logFile, nil
}

func getCommandLine() string {
	return "Invoked with following arguments: \n" + strings.Join(os.Args, " \\ \n")
}
