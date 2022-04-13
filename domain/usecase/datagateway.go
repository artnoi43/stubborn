package usecase

// package "github.com/patrickmn/go-cache"
// does not return error when calling Set
type AnswerDataGateway interface {
	Get(CacheKey) (CachedAnswer, bool)
	Set(CacheKey, CachedAnswer)
	SetMap(map[CacheKey]CachedAnswer)
	Shutdown()
}
