package main

import (
	"github.com/juxuny/webhook/executor"
	"log"
)

type defaultLogger struct {
}

func (t *defaultLogger) Printf(format string, values ...interface{}) {
	log.Printf(format, values...)
}

func (t *defaultLogger) Println(message ...interface{}) {
	log.Println(message...)
}

func NewDefaultLogger() executor.Logger {
	return &defaultLogger{}
}
