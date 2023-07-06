package account_test

import (
	"passwords/internal/domain/account"

	"testing"

	// mockAccount "passwords/internal/mocks/account"

	"github.com/golang/mock/gomock"
)

func TestGetUsersAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// s := mockAccount.NewMockStorage(ctrl)

	tests := []struct {
		name           string
		userName       string
		returnAccounts account.Accounts
		returnErr      error
		want           error
	}{
		{
			name:           "#1 Positive",
			userName:       "testuser",
			returnAccounts: account.Accounts{"1233": account.Account{Account: "1233", User: "testuser"}},
			returnErr:      nil,
			want:           nil,
		},
		{
			name:     "#2 Negative - no Accounts for user",
			userName: "testuser",
			// returnAccounts: account.Accounts{1233: account.Account{Account: "1233", User: "testuser2", Status: account.NewAccount.Text}},
			returnErr: account.ErrNoAccounts,
			want:      account.ErrNoAccounts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {

			// s.EXPECT().GetAccountsList(context.Background(), tt.userName).Return(tt.returnAccounts, tt.returnErr)
			// service := account.NewService(s)

			// _, err := service.GetUsersAccounts(context.Background(), tt.userName)
			// log.Println(err)

			// if !assert.Equal(t, true, errors.Is(err, tt.want)) {
			// 	t.Errorf("Error %v but want %v", err, tt.want)
			// }

		})

	}
}
