package trm

//go:generate mockgen -source=$GOFILE -destination=drivers/mock/$GOFILE -package mock

import (
	"context"
	"errors"
	"fmt"
)

type TrType int8

const (
	IncorrectTrType TrType = iota
	IndependentTransaction
	NestedTransaction
	NoTransaction
)

type TrFactory interface {
	BeginTx(context.Context, Config) (context.Context, Transaction, error)
}

type NestedTrFactory interface {
	Begin(context.Context, Config) (context.Context, Transaction, error)
}

type Transaction interface {
	Transaction() interface{}

	Commit(context.Context) error
	Rollback(context.Context) error

	IsActive() bool
}

var (
	ErrAlreadyClosed = errors.New("transaction already closed")

	ErrBegin    = errors.New("transaction begin")
	ErrRollback = errors.New("transaction rollback")
	ErrCommit   = errors.New("transaction commit")

	ErrNestedBegin    = fmt.Errorf("nested %w", ErrBegin)
	ErrNestedRollback = fmt.Errorf("nested %w", ErrRollback)
	ErrNestedCommit   = fmt.Errorf("nested %w", ErrCommit)
)
