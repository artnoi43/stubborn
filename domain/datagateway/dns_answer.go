package datagateway

import (
	"github.com/artnoi43/stubborn/domain/usecase"
)

// package "github.com/patrickmn/go-cache"
// does not return error when calling Set
type AnswerDataGateway interface {
	Get(usecase.CacheKey) (usecase.CachedAnswer, bool)
	Set(usecase.CacheKey, usecase.CachedAnswer)
	SetMap(map[usecase.CacheKey]usecase.CachedAnswer)
	Shutdown()
}
