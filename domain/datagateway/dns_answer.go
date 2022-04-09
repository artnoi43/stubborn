package datagateway

import (
	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

// package "github.com/patrickmn/go-cache"
// does not return error when calling Set
type AnswerDataGateway interface {
	Get(dnsutils.Key) (entity.CachedAnswer, bool)
	Set(dnsutils.Key, entity.CachedAnswer)
	SetMap(map[dnsutils.Key]entity.CachedAnswer)
	Shutdown()
}
