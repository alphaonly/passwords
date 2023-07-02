package account

import "context"

type Storage interface {
	GetOrder(ctx context.Context, orderNumber int64) (o *Account, err error)
	SaveOrder(ctx context.Context, o Account) (err error)
	GetOrdersList(ctx context.Context, userName string) (ol Orders, err error)
	GetNewOrdersList(ctx context.Context) (ol Orders, err error)
}
