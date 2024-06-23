package config

import (
	"time"

	"github.com/nlnaa11/transaction-manager/trm"
)

type ctxKey struct{}

var DefaultCtxKey = ctxKey{}

type SetUpCancel int8

const (
	IncorrectSetUpCancel SetUpCancel = iota
	Cancellable
	NonCancellable
)

var (
	defaultCancellability = IncorrectSetUpCancel
)

type config struct {
	ctxKey trm.CtxKey
	trType trm.TrType

	timeout        time.Duration
	cancellability SetUpCancel
}

type Opt func(c *config) error

func NewTrConfig(trType trm.TrType, oo ...Opt) (config, error) {
	c := &config{
		trType:         trType,
		ctxKey:         DefaultCtxKey,
		cancellability: defaultCancellability,
	}

	for _, o := range oo {
		if err := o(c); err != nil {
			return config{}, err
		}
	}

	return *c, nil
}

func MustTrConfig(trType trm.TrType, oo ...Opt) config {
	c, err := NewTrConfig(trType, oo...)
	if err != nil {
		panic(err)
	}

	return c
}

func (c config) CtxKey() trm.CtxKey {
	return c.ctxKey
}

func (c config) CtxKeyWithFlag() (trm.CtxKey, bool) {
	return c.CtxKey(), c.ctxKey != DefaultCtxKey
}

func (c config) SetCtxKey(key trm.CtxKey) trm.Config {
	return c.setCtxKey(key)
}

func (c config) setCtxKey(key trm.CtxKey) config {
	c.ctxKey = key

	return c
}

func (c config) Timeout() time.Duration {
	return c.timeout
}

func (c config) TimeoutWithFlag() (time.Duration, bool) {
	return c.Timeout(), c.timeout != 0
}

func (c config) SetTimeout(timeout time.Duration) trm.Config {
	return c.setTimeout(timeout)
}

func (c config) setTimeout(timeout time.Duration) config {
	c.timeout = timeout

	return c
}

func (c config) Cancellable() bool {
	return c.cancellability == Cancellable
}

func (c config) CancellableWithFlag() (bool, bool) {
	return c.Cancellable(), c.cancellability != IncorrectSetUpCancel
}

func (c config) SetCancellable(cancellable bool) trm.Config {
	return c.setCancellable(cancellable)
}

func (c config) setCancellable(cancellabe bool) config {
	if cancellabe {
		c.cancellability = Cancellable
	} else {
		c.cancellability = NonCancellable
	}

	return c
}

func (c config) TransactionType() trm.TrType {
	return c.trType
}

func (c config) SetTransactionType(typ trm.TrType) trm.Config {
	c.trType = typ

	return c
}

func (c config) InheritFrom(external trm.Config) (res trm.Config) {
	res = c

	if _, ok := c.CtxKeyWithFlag(); !ok {
		res = res.SetCtxKey(external.CtxKey())
	}
	if t, ok := c.TimeoutWithFlag(); !ok || t > external.Timeout() {
		res = res.SetTimeout(external.Timeout())
	}
	if _, ok := c.CancellableWithFlag(); !ok {
		res = res.SetCancellable(external.Cancellable())
	}

	return res
}
