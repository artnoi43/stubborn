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
)

func init() {
	// f receives extra config from the command-line,
	// and its values are sent to InitConfig()
	// **to overwrite** config from config.yaml.
	var f cmd.Flags
	f.Parse()
	log.Println("using config file", f.ConfigFile)
	conf, err := config.InitConfig(f)
	if err != nil {
		log.Fatalln("bad config:", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	stubborn = controller.New(ctx, conf)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	go func() {
		sig := <-sigChan
		log.Printf("received signal: %s\n", sig.String())
		// Calling cancel will trigger application-wide shutdown
		cancel()
	}()
}

func main() {
	stubborn.Run()
}
