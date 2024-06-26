package pgxv5

import (
	"context"

	"github.com/nlnaa11/transaction-manager/trm"
	trmcontext "github.com/nlnaa11/transaction-manager/trm/context"
)

var DefaultTrProvider = MustTrProvider(trmcontext.DefaultCtxManager)

type trProvider struct {
	ctxManager trm.CtxManager
}

func NewTrProvider(ctxManager trm.CtxManager) (*trProvider, error) {
	return &trProvider{
		ctxManager: ctxManager,
	}, nil
}

func MustTrProvider(ctxManager trm.CtxManager) *trProvider {
	p, err := NewTrProvider(ctxManager)
	if err != nil {
		panic(err)
	}

	return p
}

func (p *trProvider) DefaultTrOrDB(ctx context.Context, db Tr) Tr {
	if tr := p.ctxManager.Default(ctx); tr != nil {
		return p.convert(tr)
	}

	return db
}

func (p *trProvider) TrOrDB(ctx context.Context, key trm.CtxKey, db Tr) Tr {
	if tr := p.ctxManager.ByKey(ctx, key); tr != nil {
		return p.convert(tr)
	}

	return db
}

func (p *trProvider) convert(tr trm.Transaction) Tr {
	if tx, ok := tr.Transaction().(Tr); ok {
		return tx
	}

	return nil
}
