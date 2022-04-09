package dataadapter

import (
	"log"

	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/data/repository"
	"github.com/artnoi43/stubborn/domain/datagateway"
	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

// dataAdapter implements datagateway.AnswerDataGateway
type dataAdapter struct {
	cacher repository.Cacher
}

const (
	notFound bool = false
	found    bool = true
)

func New(repoConf *repository.Config) datagateway.AnswerDataGateway {
	c := repository.New(repoConf)
	return &dataAdapter{
		cacher: c,
	}
}

func (a *dataAdapter) Get(k dnsutils.Key) (entity.CachedAnswer, bool) {
	answers, err := a.cacher.GetRR(k)
	if err != nil {
		if errors.Is(err, repository.ErrCacheMissed) {
			return nil, notFound
		}
	}
	return answers, found
}

func (a *dataAdapter) Set(k dnsutils.Key, v entity.CachedAnswer) {
	if len(v) > 0 {
		a.cacher.SetRR(k, v)
	}
}

func (a *dataAdapter) SetMap(m map[dnsutils.Key]entity.CachedAnswer) {
	a.cacher.SetMapRR(m)
}

func (a *dataAdapter) Shutdown() {
	log.Println("shutting down dataAdapter")
	a.cacher.Flush()
	log.Println("dataAdapter shutdown gracefully")
}
