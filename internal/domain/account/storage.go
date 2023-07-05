package account

import "context"

type Storage interface {
	GetAccount(ctx context.Context, user string, account string) (acc *Account, err error)
	SaveAccount(ctx context.Context, account Account) (err error)
	GetAccountsList(ctx context.Context, user string) (accounts Accounts, err error)
}
