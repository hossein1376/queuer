package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/hossein1376/querier/config"
	"github.com/hossein1376/querier/pkg/order/worker"
	"github.com/hossein1376/querier/pkg/streamer"
)

func main() {
	var path string
	flag.StringVar(
		&path, "config", "assets/cfg.yaml", "path to the config file",
	)
	flag.Parse()

	cfg, err := config.New(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read configs: %s", err)
		return
	}

	// listen for the interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go worker.Run(cfg.Worker)
	go streamer.Run(cfg.Streamer)

	// wait until the interrupt signal is received
	<-interrupt
}
