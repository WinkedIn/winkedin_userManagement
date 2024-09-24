package logger

import (
	"log"
	"os"
	"sync"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
	once        sync.Once
)

func Init() {
	once.Do(func() {
		InfoLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
		ErrorLogger = log.New(os.Stdout, "ERROR: ", log.LstdFlags)
		WarnLogger = log.New(os.Stdout, "WARN: ", log.LstdFlags)
	})
}
