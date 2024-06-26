package pgxv5

import (
	"context"

	"github.com/nlnaa11/transaction-manager/trm"
)

type trFactory struct {
	db Transactor
}

func NewTrFactory(db Transactor) *trFactory {
	return &trFactory{db}
}

func (f *trFactory) BeginTx(
	ctx context.Context,
	cfg trm.Config,
) (context.Context, trm.Transaction, error) {
	c, _ := cfg.(Config)

	return NewTransaction(ctx, c.TxOptions(), f.db)
}
