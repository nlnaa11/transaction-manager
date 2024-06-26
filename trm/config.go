package trm

//go:generate mockgen -source=$GOFILE -destination=drivers/mock/$GOFILE -package mock

import "time"

type Config interface {
	CtxKey() CtxKey
	CtxKeyWithFlag() (CtxKey, bool)
	SetCtxKey(CtxKey) Config

	Timeout() time.Duration
	TimeoutWithFlag() (time.Duration, bool)
	SetTimeout(time.Duration) Config

	Cancellable() bool
	CancellableWithFlag() (bool, bool)
	SetCancellable(bool) Config

	TransactionType() TrType
	SetTransactionType(typ TrType) Config

	InheritFrom(Config) Config
}
