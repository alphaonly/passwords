package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"passwords/internal/pkg/common/logging"
	"passwords/internal/schema"
	"time"

	pb "passwords/internal/adapters/grpc/proto"
	conf "passwords/internal/configuration/client"
	accountDomain "passwords/internal/domain/account"
	userDomain "passwords/internal/domain/user"
	grpcCl "passwords/internal/pkg/client/grpc/client"
	"passwords/internal/pkg/client/service"
)

type myClient struct {
	Configuration *conf.ClientConfiguration
	GRPCClient    *grpcCl.GRPCClient
}

func NewClient(c *conf.ClientConfiguration, grpcClient *grpcCl.GRPCClient) service.Service {

	return &myClient{
		Configuration: c,
		GRPCClient:    grpcClient,
	}
}

func (c myClient) AuthorizeUser(ctx context.Context, user string, password string) (*userDomain.User, error) {
	response, err := c.GRPCClient.Client.AuthorizeUser(ctx, &pb.AuthUserRequest{User: user, Password: password})
	if err != nil {
		return nil, fmt.Errorf(service.ErrWrongUserOrPassword.Error()+" %w", user, err)
	}
	log.Printf("User %v authorized ", response.User.Name)

	return &userDomain.User{
		User:     response.User.User,
		Password: response.User.Password,
		Name:     response.User.Name,
		Surname:  response.User.Surname,
		Phone:    response.User.Phone,
	}, nil
}
func (c myClient) AddNewUser(ctx context.Context, user userDomain.User) error {
	var request pb.AddUserRequest

	request.User = &pb.User{
		User:     user.User,
		Password: user.Password,
		Name:     user.Name,
		Surname:  user.Surname,
		Phone:    user.Phone,
	}

	response, err := c.GRPCClient.Client.AddUser(ctx, &request)
	if err != nil {
		log.Println(err)
		return err
	}
	if response.Error != "" {
		err = errors.New(response.Error)
	}
	return err
}
func (c myClient) AddNewAccount(ctx context.Context, account accountDomain.Account) error {
	var request pb.AddAccountRequest

	request.Account = &pb.Account{
		User:         account.User,
		Account:      account.Account,
		AccountLogin: account.Login,
		Password:     account.Password,
		Description:  account.Description,
	}

	response, err := c.GRPCClient.Client.AddAccount(ctx, &request)
	if err != nil {
		log.Println(err)
		return err
	}
	if response.Error != "" {
		err = errors.New(response.Error)
	}
	return err
}
func (c myClient) GetAllAccounts(ctx context.Context, user string) (accountDomain.Accounts, error) {
	var request pb.GetAllAccountsRequest

	request.UserLogin = user

	response, err := c.GRPCClient.Client.GetAllUserAccounts(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if response.Error != "" {
		err = errors.New(response.Error)
		return nil, err
	}

	accounts := make(accountDomain.Accounts)
	for _, v := range response.Accounts {
		created, err := time.Parse(time.RFC3339, v.Created)
		logging.LogFatal(err)
		accounts[v.Account] = accountDomain.Account{
			User:        v.User,
			Account:     v.Account,
			Login:       v.AccountLogin,
			Password:    v.Password,
			Description: v.Description,
			Created:     schema.CreatedTime(created),
		}
	}

	return nil, err
}

func (c myClient) Run(ctx context.Context) {

}

func (c myClient) Stop() {

}
