package manager

import (
	"context"
	"errors"

	"github.com/nlnaa11/transaction-manager/trm"
)

type trManager struct {
	// создатель транзакций определенного типа
	trFactory trm.TrFactory
	// провайдер транзакций по ключу
	ctxManager trm.CtxManager
	// настройки транзакций
	trConfig trm.Config
}

type Opt func(*trManager) error

func NewTrManager(trFactory trm.TrFactory) *trManager {
	return &trManager{
		trFactory: trFactory,
	}
}

func (m *trManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DoWithConfig(ctx, m.trConfig, fn)
}

func (m *trManager) DoWithConfig(ctx context.Context, config trm.Config, fn func(ctx context.Context) error) (err error) {
	ctx, closer, err := m.initTransaction(ctx, config.InheritFrom(m.trConfig))
	if err != nil {
		return err
	}

	defer func() { err = closer(ctx, recover(), &err) }()

	return fn(ctx)
}

func (m *trManager) initTransaction(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	if initTr, ok := initByTrType[config.TransactionType()]; ok {
		return initTr(m, ctx, config)
	}

	return nil, nil, errors.New("invalid transaction type")
}
