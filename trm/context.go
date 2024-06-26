package trm

//go:generate mockgen -source=$GOFILE -destination=drivers/mock/$GOFILE -package mock

import (
	"context"
)

type CtxKey interface{}

type TrProvider func(ctx context.Context) Transaction

type CtxManager interface {
	Default(ctx context.Context) Transaction
	SetDefault(ctx context.Context, tr Transaction) context.Context

	ByKey(ctx context.Context, key CtxKey) Transaction
	SetByKey(ctx context.Context, key CtxKey, tr Transaction) context.Context
}
