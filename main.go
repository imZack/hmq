package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/fhmq/hmq/broker"
	"github.com/fhmq/hmq/logger"
	"go.uber.org/zap"
)

var log = logger.Get()

func getKey(id string) ([]byte, error) {
	fmt.Println("getKey called!", id)
	return []byte("123451"), nil
}

func main() {
	config, err := broker.ConfigureConfig(os.Args[1:])
	if err != nil {
		log.Fatal("configure broker config error", zap.Error(err))
	}

	config.TlsInfo.PSK = true
	config.TlsInfo.PSKGetKey = getKey

	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("New Broker error: ", zap.Error(err))
	}
	b.Start()

	s := waitForSignal()
	log.Info("signal received, broker closed.", zap.Any("signal", s))
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
