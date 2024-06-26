package pgxv5

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/nlnaa11/transaction-manager/trm"
	"github.com/nlnaa11/transaction-manager/trm/drivers"
)

type Transaction struct {
	tx     pgx.Tx
	closer drivers.Closer
}

func NewTransaction(
	ctx context.Context,
	opts pgx.TxOptions,
	db Transactor,
) (context.Context, *Transaction, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{
		tx:     tx,
		closer: drivers.NewCloser(),
	}

	go tr.waitDone(ctx)

	return ctx, tr, nil
}

func (t *Transaction) waitDone(ctx context.Context) {
	// н, у WithoutCancel ctx doneCh == nil
	if ctx.Done() == nil {
		return
	}

	select {
	case <-ctx.Done():
		t.closer.Close()
	case <-t.closer.Closed():
		return
	}
}

func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin начинает псевдо-вложенную транзакцию
func (t *Transaction) Begin(
	ctx context.Context,
	_ trm.Config,
) (context.Context, trm.Transaction, error) {
	tx, err := t.tx.Begin(ctx)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{
		tx:     tx,
		closer: drivers.NewCloser(),
	}

	return ctx, tr, nil
}

func (t *Transaction) Commit(ctx context.Context) error {
	defer t.closer.Close()

	return t.tx.Commit(ctx)
}

func (t *Transaction) Rollback(ctx context.Context) error {
	defer t.closer.Close()

	return t.tx.Rollback(ctx)
}

func (t *Transaction) IsActive() bool {
	select {
	case <-t.closer.Closed():
		return false
	default:
		return true
	}
}
