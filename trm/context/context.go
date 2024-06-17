package context

import (
	"context"

	"github.com/nlnaa11/transaction-manager/trm"
	"github.com/nlnaa11/transaction-manager/trm/config"
)

var DefaultCtxManager = NewCtxManager(config.DefaultCtxKey)

type ctxManager struct {
	defaultKey trm.CtxKey
}

func NewCtxManager(key trm.CtxKey) *ctxManager {
	return &ctxManager{
		defaultKey: key,
	}
}

func (m *ctxManager) Default(ctx context.Context) trm.Transaction {
	return m.ByKey(ctx, m.defaultKey)
}

func (m *ctxManager) SetDefault(ctx context.Context, tr trm.Transaction) context.Context {
	return m.SetByKey(ctx, m.defaultKey, tr)
}

func (m *ctxManager) ByKey(ctx context.Context, key trm.CtxKey) trm.Transaction {
	if tr, ok := ctx.Value(key).(trm.Transaction); ok {
		return tr
	}

	return nil
}

func (m *ctxManager) SetByKey(ctx context.Context, key trm.CtxKey, tr trm.Transaction) context.Context {
	return context.WithValue(ctx, key, tr)
}
