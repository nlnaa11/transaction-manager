package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/nlnaa11/transaction-manager/trm"
)

type TrCloser func(context.Context, error) error

type InitTr func(*trManager, context.Context, trm.Config) (context.Context, TrCloser, error)

var initByTrType = map[trm.TrType]InitTr{
	trm.NestedIndependentTr: (*trManager).nestedIndependentTr,
	trm.PseudoNestedTr:      (*trManager).pseudoNestedTr,
	trm.NestedSubTr:         (*trManager).nestedSubTr,
}

func (m *trManager) nestedIndependentTr(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	ctx, cancel := m.withCancel(ctx, config)

	ctx, tr, err := m.trFactory(ctx)
	if err != nil {
		return nil, nil, errors.New("can't create new nested independent transaction")
	}

	return m.ctxManager.SetByKey(ctx, config.CtxKey, tr),
		nil,
		nil
}

func (m *trManager) nestedSubTr(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	tr := m.ctxManager.ByKey(ctx, config.CtxKey())
	ctx, cancel := m.withCancel(ctx, config)

	if tr == nil {
		return nil, nil, errors.New("can't get main pseudo nested transaction for subtransaction")
	}

	return ctx, nil, nil
}

func (m *trManager) pseudoNestedTr(ctx context.Context, config trm.Config) (context.Context, TrCloser, error) {
	tr := m.ctxManager.ByKey(ctx, config.CtxKey())
	ctx, cancel := m.withCancel(ctx, config)

	if tr == nil {
		ctx, tr, err := m.trFactory(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("can't create new pseudo nested transaction: %w", err)
		}

		return m.ctxManager.SetByKey(ctx, config.CtxKey, tr),
			nil,
			nil
	}

	if nestedFactory, ok := tr.(trm.NestedTrFactory); ok {
		ctx, tr, err := nestedFactory.Begin(ctx, config)
		if err != nil {
			return nil, nil, fmt.Errorf("can't begin pseudo nested transaction: %w", err)
		}

		return m.ctxManager.SetByKey(ctx, config.CtxKey(), tr),
			nil,
			nil
	}

	return ctx, nil, nil
}

func (m *trManager) withCancel(ctx context.Context, config trm.Config) (context.Context, context.CancelFunc) {
	if t, ok := config.TimeoutWithFlag(); ok {
		return context.WithTimeout(ctx, t)
	}

	return ctx, nilCancelFunc
}

func nilCancelFunc()
