package main

import (
	"context"
	"time"

	"github.com/athanbase/log"
	"github.com/segmentio/kafka-go"
)

type BrokerWriter struct {
	writer *kafka.Writer
}

func (w BrokerWriter) Write(b []byte) (int, error) {
	b1 := make([]byte, len(b))
	copy(b1, b) // b is reused
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	err := w.writer.WriteMessages(
		ctx,
		kafka.Message{Value: b1},
	)
	cancel()
	if err != nil {
		panic(err)
	}
	return len(b1), err
}

func (w BrokerWriter) Sync() error {
	return w.writer.Close()
}

func main() {
	w := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "logs",
		Balancer:     &kafka.RoundRobin{},
		Async:        true,
		BatchTimeout: time.Millisecond * 10,
		BatchSize:    2,
	}

	bw := &BrokerWriter{writer: w}

	logger := log.New(bw, log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))
	log.ResetDefault(logger)

	log.Infof("%s", "test write to kafka info")
	log.Debugf("%s", "test debug")

	logger.SetLevel(log.DebugLevel)

	log.Debugf("%s", "test level changed to debug")
	defer log.Sync()
}
