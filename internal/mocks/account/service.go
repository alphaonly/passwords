package mocks

import (
	"context"
	"fmt"
	"passwords/internal/domain/account"
	"passwords/internal/schema"
)

type service struct {
}

func NewaccountStorage() account.Storage { return nil }

func NewService() (sr account.Service) {
	return service{}
}

var testUser = "testuser"

var testaccount1122 = account.Account{
	Account: "1122",
	User:    testUser,
	Created: schema.CreatedTime{},
}

func (sr service) GetUsersAccounts(ctx context.Context, userName string) (accounts account.Accounts, err error) {
	// data validation
	if userName == "" {
		account.ErrBadUserAccount = fmt.Errorf("400 user is empty %v (%w)", userName, account.ErrBadUserAccount)
		return nil, account.ErrBadUserAccount
	}
	//getaccounts
	if userName == testUser {
		return account.Accounts{"1122": testaccount1122}, nil
	}
	return nil, account.ErrNoAccounts

}

func (sr service) ValidateaccountNumber(ctx context.Context, accountNumberStr string, user string) (accountNum int64, err error) {

	return 0, err
}
