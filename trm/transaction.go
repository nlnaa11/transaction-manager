package trm

import "context"

type TrType int32

const (
	IndependentTransaction TrType = 1 << iota
	NestedTransaction
	NoTransaction
)

type TrFactory func(context.Context, Config) (context.Context, Transaction, error)

type NestedTrFactory interface {
	Begin(context.Context, Config) (context.Context, Transaction, error)
}

type Transaction interface {
	Transaction() interface{}

	Commit() error
	Rollback() error

	Close()
}
