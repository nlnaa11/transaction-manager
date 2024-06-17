package trm

import (
	"context"
)

type TransactionManager interface {
	Do(context.Context, func(context.Context) error) error
	DoWithConfig(context.Context, Config, func(context.Context) error) error
}
