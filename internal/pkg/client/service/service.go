package service

import (
	"context"
	"errors"
	accountDomain "passwords/internal/domain/account"
	userDomain "passwords/internal/domain/user"
)

var (
	ErrWrongUserOrPassword = errors.New("User %v does not exists or wrong password")
	ErrUserNotCreated      = errors.New("User %v was not created ")
	ErrUserExists          = errors.New("User %v  already exists ")

	ErrAccountNotCreated = errors.New("Account %v was not created")
	ErrAccountExists     = errors.New("Account %v  already exists")
	ErrGetAccounts       = errors.New("Getting accounts by user %v is impossible")
)

type clientService interface {
	Run(ctx context.Context)
	Stop()
}

type dataService interface {
	AuthorizeUser(ctx context.Context, user string, password string) (*userDomain.User, error)
	AddNewUser(ctx context.Context, user userDomain.User) error
	AddNewAccount(ctx context.Context, account accountDomain.Account) error
	GetAllAccounts(ctx context.Context, user string) (accountDomain.Accounts, error)
	Run(ctx context.Context)
	Stop()
}

type Service interface {
	clientService
	dataService
}
