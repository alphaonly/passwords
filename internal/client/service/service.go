package service

//package describes a client service that work with grpc  operating with  model entities
import (
	"context"
	"errors"
	accountDomain "passwords/internal/domain/account"
	userDomain "passwords/internal/domain/user"
)

var (
	ErrWrongUserOrPassword   = errors.New("wrong user  or password")
	ErrNoUserExists          = errors.New("user does not exists ")
	ErrNoUserOrPassword      = errors.New("no user or password given")
	ErrUserAlreadyAuthorized = errors.New("user already authorized")
	ErrUserWasNotAuthorized  = errors.New("user was not authorized")

	Errors = []error{
		ErrWrongUserOrPassword,
		ErrNoUserExists,
		ErrNoUserOrPassword,
		ErrUserAlreadyAuthorized,
		ErrUserWasNotAuthorized,
	}
)

type clientService interface {
	Run(ctx context.Context)
	Stop()
}

type dataService interface {
	AuthorizeUser(ctx context.Context, user string, password string) (*userDomain.User, error)
	DisAuthorizeUser(ctx context.Context, user string, password string) error
	CheckUserAuthorization(ctx context.Context, user string) error
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
