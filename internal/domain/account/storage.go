package account

import "context"

type Storage interface {
	GetAccount(ctx context.Context, name string) (account *Account, err error)
	SaveAccount(ctx context.Context, account Account) (err error)
	GetAccountsList(ctx context.Context, userName string) (accounts Accounts, err error)
}
