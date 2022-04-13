package adapter

import (
	"log"

	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/data/repository"
	"github.com/artnoi43/stubborn/domain/usecase"
)

// answerCacherAdapter implements datagateway.AnswerDataGateway
// using repository.Cacher as the implementation.
type answerCacherAdapter struct {
	cacher repository.Cacher
}

const (
	notFound bool = false
	found    bool = true
)

// New returns answerCacherAdapter,
// which implements datagateway.AnswerDataGateway
// using repository.Cacher as the implementation.
func New(conf *repository.Config) *answerCacherAdapter {
	c := repository.New(conf)
	return &answerCacherAdapter{
		cacher: c,
	}
}

func (a *answerCacherAdapter) Get(k usecase.CacheKey) (usecase.CachedAnswer, bool) {
	answers, err := a.cacher.GetRR(k)
	if err != nil {
		if errors.Is(err, repository.ErrCacheMissed) {
			return nil, notFound
		}
	}
	return answers, found
}

func (a *answerCacherAdapter) Set(k usecase.CacheKey, v usecase.CachedAnswer) {
	if len(v) > 0 {
		a.cacher.SetRR(k, v)
	}
}

func (a *answerCacherAdapter) SetMap(m map[usecase.CacheKey]usecase.CachedAnswer) {
	if m == nil {
		log.Println("SetMap failed: nil map")
		return
	}
	a.cacher.SetMapRR(m)
}

func (a *answerCacherAdapter) Shutdown() {
	log.Println("shutting down answerCacherAdapter")
	a.cacher.Flush()
	log.Println("answerCacherAdapter shutdown gracefully")
}
