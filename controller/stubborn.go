package controller

import (
	"context"
	"log"
	"sync"

	"github.com/artnoi43/stubborn/config"
	"github.com/artnoi43/stubborn/domain/dataadapter"
	"github.com/artnoi43/stubborn/domain/usecase/dnsserver"
	"github.com/artnoi43/stubborn/domain/usecase/handler"
)

type Stubborn interface {
	Run()
}

type stubborn struct {
	ctx     context.Context
	handler handler.Handler
}

func New(ctx context.Context, conf *config.Config) Stubborn {
	// Init dnsServer, answerDataGateway, and handler
	s := dnsserver.NewDNSServer(&conf.ServerConfig)
	d := dataadapter.New(&conf.CacherConfig)
	h := handler.New(ctx, &conf.HandlerConfig, s, d)
	return &stubborn{
		ctx:     ctx,
		handler: h,
	}
}

func (app *stubborn) Run() {
	// Calling start will start app.dnsServer
	// and handlers will be registered with dns.HandleFunc
	go app.handler.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-app.ctx.Done()
		app.handler.Shutdown()
		log.Println("stubborn shutdown gracefully")
	}()
	wg.Wait()
}
