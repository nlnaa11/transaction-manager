package manager

import (
	"context"

	"github.com/nlnaa11/transaction-manager/trm"
	"go.uber.org/multierr"
)

type TrCloser func(context.Context, interface{}, *error) error

type InitTr func(*trManager, context.Context, trm.Config) (context.Context, TrCloser, error)

var initByTrType = map[trm.TrType]InitTr{
	trm.IndependentTransaction: (*trManager).independentTransaction,
	trm.NestedTransaction:      (*trManager).nestedTransaction,
	trm.NoTransaction:          (*trManager).noTransaction,
}

func (m *trManager) independentTransaction(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	ctx, cancel := m.withCancel(ctx, config)

	ctx, tr, err := m.trFactory(ctx, config)
	if err != nil {
		return nil, nil, multierr.Combine(trm.ErrBegin, err)
	}

	return m.ctxManager.SetByKey(ctx, config.CtxKey, tr),
		newTrClose(tr, cancel),
		nil
}

func (m *trManager) nestedTransaction(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	tr := m.ctxManager.ByKey(ctx, config.CtxKey())
	ctx, cancel := m.withCancel(ctx, config)

	if tr == nil {
		ctx, tr, err := m.trFactory(ctx, config)
		if err != nil {
			return nil, nil, multierr.Combine(trm.ErrBegin, err)
		}

		return m.ctxManager.SetByKey(ctx, config.CtxKey, tr),
			newTrClose(tr, cancel),
			nil
	}

	if nestedFactory, ok := tr.(trm.NestedTrFactory); ok {
		ctx, tr, err := nestedFactory.Begin(ctx, config)
		if err != nil {
			return nil, nil, multierr.Combine(trm.ErrNestedBegin, err)
		}

		return m.ctxManager.SetByKey(ctx, config.CtxKey(), tr),
			newTrClose(tr, cancel),
			nil
	}

	return ctx, newNilClose(cancel), nil
}

func (m *trManager) noTransaction(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	tr := m.ctxManager.ByKey(ctx, config.CtxKey())
	ctx, cancel := m.withCancel(ctx, config)

	if tr == nil {
		return ctx, newNilClose(cancel), nil
	}

	return m.ctxManager.SetByKey(ctx, config.CtxKey, nil),
		newNilClose(cancel),
		nil
}

func (m *trManager) withCancel(ctx context.Context, config trm.Config) (context.Context, context.CancelFunc) {
	if t, ok := config.TimeoutWithFlag(); ok {
		return context.WithTimeout(ctx, t)
	}

	if config.Cancellable() {
		return context.WithCancel(ctx)
	}

	return context.WithoutCancel(ctx), nilCancelFunc
}

func nilCancelFunc() {}
