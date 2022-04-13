package handler

import (
	"github.com/artnoi43/stubborn/domain/usecase"
)

type dnsServer interface {
	ListenAndServe() error
	Addr() string
	Shutdown()
}

type dnsClient interface {
	Query(v interface{}) (*usecase.Answer, error)
	Shutdown()
}
