package account

import (
	"context"
	"fmt"
)

var (
	ErrBadUserAccount      = fmt.Errorf("400 user is empty or bad account number")
	ErrNoAccounts          = fmt.Errorf("204 no Accounts")
	ErrNoLuhnNumber        = fmt.Errorf("422 not Lihn number")
	ErrAccountNumberExists = fmt.Errorf("200 account exists with user")
	ErrAnotherUsersAccount = fmt.Errorf("409 account exists with another user")
)

type Service interface {
	GetUsersAccounts(ctx context.Context, userName string) (Accounts Accounts, err error)
}

type service struct {
	Storage Storage
}

func NewService(s Storage) (sr Service) {
	return service{Storage: s}
}

func (sr service) GetUsersAccounts(ctx context.Context, userName string) (Accounts Accounts, err error) {
	// data validation
	if userName == "" {
		ErrBadUserAccount = fmt.Errorf("400 user is empty %v (%w)", userName, ErrBadUserAccount)
		return nil, ErrBadUserAccount
	}
	//getAccounts
	Accountslist, err := sr.Storage.GetAccountsList(ctx, userName)
	if err != nil {
		ErrNoAccounts = fmt.Errorf(err.Error()+"(%w)", ErrNoAccounts)
		ErrNoAccounts = fmt.Errorf("204 no Accounts for user %v %w", userName, ErrNoAccounts)
		return nil, ErrNoAccounts
	}
	return Accountslist, nil
}
