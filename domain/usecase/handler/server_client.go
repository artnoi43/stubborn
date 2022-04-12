package handler

import (
	"github.com/artnoi43/stubborn/domain/entity"
)

type dnsServer interface {
	ListenAndServe() error
	Addr() string
	Shutdown()
}

type dnsClient interface {
	Query(v interface{}) (*entity.Answer, error)
	Shutdown()
}
