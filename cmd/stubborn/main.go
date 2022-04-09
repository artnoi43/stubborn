package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/artnoi43/stubborn/cmd"
	"github.com/artnoi43/stubborn/config"
	"github.com/artnoi43/stubborn/controller"
)

var (
	stubborn controller.Stubborn
	ctx      context.Context
	cancel   context.CancelFunc
	sigChan  chan os.Signal
)

func init() {
	var f cmd.Flags
	f.Parse()
	conf, err := config.InitConfig(f.ConfigFile)
	if err != nil {
		log.Fatalln("bad config:", err.Error())
	}
	ctx, cancel = context.WithCancel(context.Background())
	stubborn = controller.New(ctx, conf)
	sigChan = make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
}

func main() {
	go func() {
		sig := <-sigChan
		log.Printf("received signal: %s\n", sig.String())
		// Calling cancel will trigger application-wide shutdown
		cancel()
	}()
	stubborn.Run()
}
