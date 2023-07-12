package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	grpcCl "passwords/internal/client/grpc/client"
	"passwords/internal/client/service"
	"passwords/internal/pkg/common/logging"
	"passwords/internal/schema"
	"time"

	pb "passwords/internal/adapters/grpc/proto"
	conf "passwords/internal/configuration/client"
	accountDomain "passwords/internal/domain/account"
	userDomain "passwords/internal/domain/user"
)

var (
	ErrUserDisAuthorization = fmt.Errorf("user disauthorization error")
)

type clientService struct {
	Configuration *conf.ClientConfiguration
	GRPCClient    *grpcCl.GRPCClient
}

func NewClientService(c *conf.ClientConfiguration, grpcClient *grpcCl.GRPCClient) service.Service {

	return &clientService{
		Configuration: c,
		GRPCClient:    grpcClient,
	}
}
func (c clientService) CheckUserAuthorization(ctx context.Context, user string) error {
	response, err := c.GRPCClient.Client.CheckAuthorization(ctx, &pb.CheckAuthUserRequest{User: user})
	logging.LogFatal(err)
	if response.Error != "" {
		return fmt.Errorf(response.Error+" %w", service.ErrUserWasNotAuthorized)
	}

	return err
}

func (c clientService) AuthorizeUser(ctx context.Context, user string, password string) (*userDomain.User, error) {
	if user == "" || password == "" {
		return nil, service.ErrNoUserOrPassword
	}
	response, err := c.GRPCClient.Client.AuthorizeUser(ctx, &pb.AuthUserRequest{User: user, Password: password})
	if err != nil {

		return nil, CovertStatusError(err)
		//switch {
		//
		//case strings.Contains(err.Error(), service.ErrWrongUserOrPassword.Error()): return nil, service.ErrWrongUserOrPassword
		//case strings.Contains(err.Error(), service.ErrNoUserExists.Error()):return nil, service.ErrNoUserExists
		//case strings.Contains(err.Error(), service.ErrUserWasNotAuthorized.Error()):return nil, service.ErrUserWasNotAuthorized
		//}

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

// DisAuthorizeUser - ends user authorization in server
func (c clientService) DisAuthorizeUser(ctx context.Context, user string, password string) error {
	if user == "" || password == "" {
		return service.ErrNoUserOrPassword
	}
	response, err := c.GRPCClient.Client.DisAuthorizeUser(ctx, &pb.DisAuthUserRequest{User: user, Password: password})
	if err != nil {
		return fmt.Errorf(ErrUserDisAuthorization.Error()+"%w", err)
	}
	if response.Error != "" {
		err = errors.New(response.Error)
		return fmt.Errorf(ErrUserDisAuthorization.Error()+"%w", err)
	}
	return nil
}

// AddNewUser - client adds a new user
func (c clientService) AddNewUser(ctx context.Context, user userDomain.User) error {
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
func (c clientService) AddNewAccount(ctx context.Context, account accountDomain.Account) error {
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
func (c clientService) GetAllAccounts(ctx context.Context, user string) (accountDomain.Accounts, error) {
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

	return accounts, err
}

func (c clientService) Run(ctx context.Context) {
	//if c.Configuration.Login == "" || c.Configuration.Password ==""{
	//	log.Fatal(service.ErrNoUserOrPassword)
	//}
	////Authorize user by given parameters
	//authUserRequest := pb.AuthUserRequest{User: c.Configuration.Login,Password: c.Configuration.Password}
	//response,err :=c.GRPCClient.Client.AuthorizeUser(ctx,&authUserRequest)
	//if err!=nil{
	//	log.Fatal(service.ErrWrongUserOrPassword)
	//}
	//response.

}

func (c clientService) Stop() {

}
