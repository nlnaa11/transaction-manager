package trm

//go:generate mockgen -source=$GOFILE -destination=drivers/mock/$GOFILE -package mock

import (
	"context"
)

type TransactionManager interface {
	Do(context.Context, func(context.Context) error) error
	DoWithConfig(context.Context, Config, func(context.Context) error) error
}
