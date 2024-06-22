package manager

import (
	"context"

	"github.com/nlnaa11/transaction-manager/trm"
)

func newTrClose(tr trm.Transaction, cancel context.CancelFunc) TrCloser {
	return func(ctx context.Context, p interface{}, err *error) error {

		return nil
	}
}

func newNilClose(cancel context.CancelFunc) TrCloser {
	return func(_ context.Context, p interface{}, err *error) error {
		defer cancel()

		if p != nil {
			panic(p)
		}

		if *err != nil {
			return *err
		}

		return nil
	}
}
