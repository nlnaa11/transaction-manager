package manager

import (
	"context"
	"errors"

	"github.com/nlnaa11/transaction-manager/trm"
	"go.uber.org/multierr"
)

func newTrClose(tr trm.Transaction, cancel context.CancelFunc) TrCloser {
	return func(ctx context.Context, p interface{}, processTrErr *error) error {
		defer cancel()

		// recovering from panic
		if p != nil {
			if tr.IsActive() {
				if err := tr.Rollback(ctx); err != nil {
					// log
				}
			}

			panic(p)
		}

		hasErr := *processTrErr != nil
		isCtxCanceled := errors.Is(*processTrErr, context.Canceled)
		isCtxDeadlineExceeded := errors.Is(*processTrErr, context.DeadlineExceeded)
		isCtxErr := isCtxCanceled || isCtxDeadlineExceeded

		ctxErr := ctx.Err()

		if ctxErr != nil {
			if !hasErr {
				*processTrErr = ctxErr
			} else if !isCtxCanceled && errors.Is(ctxErr, context.Canceled) ||
				!isCtxDeadlineExceeded && errors.Is(ctxErr, context.DeadlineExceeded) {
				*processTrErr = multierr.Combine(*processTrErr, ctxErr)
			}

			isCtxErr = true
			hasErr = true
		}

		if !tr.IsActive() {
			if hasErr {
				if isCtxErr || errors.Is(*processTrErr, trm.ErrAlreadyClosed) {
					return *processTrErr
				}

				return multierr.Combine(*processTrErr, trm.ErrAlreadyClosed)
			}

			return trm.ErrAlreadyClosed
		}

		if hasErr {
			if rollbackErr := tr.Rollback(ctx); rollbackErr != nil {
				return multierr.Combine(*processTrErr, trm.ErrRollback, rollbackErr)
			}

			return *processTrErr
		}

		if commitErr := tr.Commit(ctx); commitErr != nil {
			return multierr.Combine(trm.ErrCommit, commitErr)
		}

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
