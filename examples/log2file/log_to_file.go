package main

import (
	"os"

	"github.com/athanbase/log"
)

func main() {
	f, err := os.OpenFile(
		"demo1.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logger := log.New(f, log.InfoLevel, log.AddCallerSkip(1), log.WithCaller(true))
	log.ResetDefault(logger)
	defer log.Sync()

	log.Info("demo1: ", log.String("app", "start ok"), log.Int("version", 2))
}
