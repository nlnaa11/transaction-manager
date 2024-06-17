package config

import (
	"time"

	"github.com/nlnaa11/transaction-manager/trm"
)

func WithCtxKey(key trm.CtxKey) Opt {
	return func(c *config) error {
		*c = c.setCtxKey(key)

		return nil
	}
}

func WithTimeout(timeout time.Duration) Opt {
	return func(c *config) error {
		*c = c.setTimeout(timeout)

		return nil
	}
}
