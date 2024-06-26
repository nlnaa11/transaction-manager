package pgxv5

import (
	"github.com/jackc/pgx/v5"

	"github.com/nlnaa11/transaction-manager/trm"
)

type Opt func(*Config) error

func WithTxOptions(options pgx.TxOptions) Opt {
	return func(c *Config) error {
		*c = c.setTxOptions(options)

		return nil
	}
}

type Config struct {
	trm.Config
	txOpts pgx.TxOptions
}

func NewConfig(cfg trm.Config, oo ...Opt) (Config, error) {
	c := &Config{
		Config: cfg,
	}

	for _, o := range oo {
		if err := o(c); err != nil {
			return Config{}, err
		}
	}

	return *c, nil
}

func MustConfig(cfg trm.Config, oo ...Opt) Config {
	c, err := NewConfig(cfg, oo...)
	if err != nil {
		panic(err)
	}

	return c
}

func (c Config) TxOptions() pgx.TxOptions {
	return c.txOpts
}

func (c Config) setTxOptions(options pgx.TxOptions) Config {
	c.txOpts = options

	return c
}

func (c Config) InheritFrom(cfg trm.Config) trm.Config {
	external, ok := cfg.(Config)
	if ok {
		var emptyTxOpts pgx.TxOptions

		if c.txOpts == emptyTxOpts {
			c.setTxOptions(external.TxOptions())
		}

	}

	c.Config = c.Config.InheritFrom(cfg)

	return c
}
