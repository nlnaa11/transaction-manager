package manager

import "github.com/nlnaa11/transaction-manager/trm"

func WithCtxManager(ctxManager trm.CtxManager) Opt {
	return func(m *trManager) error {
		m.ctxManager = ctxManager

		return nil
	}
}

func WithTrConfig(config trm.Config) Opt {
	return func(m *trManager) error {
		m.trConfig = config

		return nil
	}
}
