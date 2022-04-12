package controller

import (
	"context"
	"log"
	"sync"

	"github.com/artnoi43/stubborn/config"
	"github.com/artnoi43/stubborn/domain/adapter"
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

// New inits dnsServer, answerDataGateway,
// and embeds them into stubborn.handler
func New(ctx context.Context, conf *config.Config) Stubborn {
	return &stubborn{
		ctx: ctx,
		handler: handler.New(
			ctx,
			&conf.HandlerConfig,
			dnsserver.New(&conf.ServerConfig),
			adapter.New(&conf.CacherConfig),
		),
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
