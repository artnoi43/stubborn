package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/config"
	"github.com/artnoi43/stubborn/lib/cacher"
	"github.com/artnoi43/stubborn/lib/dohclient"
	"github.com/artnoi43/stubborn/lib/rediswrapper"
	"github.com/artnoi43/stubborn/lib/server"
)

var (
	confLocation      = "config/config.yaml"
	confDir, confFile = path.Split(confLocation)
	confFileExt       = strings.Split(confFile, ".")
)

func init() {
	if len(confFileExt) < 2 {
		log.Fatalln("bad config file location:", confLocation)
	}
}

func main() {
	conf, err := config.InitConfig(confDir, confFileExt[0], confFileExt[1])
	if err != nil {
		log.Fatalln("bad config:", err.Error())
	}

	ctx := context.Background()
	redisCli := rediswrapper.New(ctx, &conf.RedisConfig)
	cacher := cacher.New(&conf.CacherConfig, redisCli)
	dnsServer := server.NewDNSServer(&conf.ServerConfig)
	dohClient := dohclient.New()
	server := server.New(ctx, &conf.ServerConfig, dnsServer, dohClient, cacher)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-sigChan
		server.DnsServer.Shutdown()
		server.DohClient.Close()
		server.Cacher.Redis.Cli.FlushDB(ctx)
		server.Cacher.Redis.Cli.Close()
		log.Println("Shutting down workers")
	}()

	dns.HandleFunc(".", server.HandleDnsReq)
	if err := server.Start(); err != nil {
		log.Fatalln("DNS Server error", err.Error())
	}
	wg.Wait()
	os.Exit(0)
}
