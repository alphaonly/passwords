package service

import (
	"context"
	"log"
	"sync"

	pb "passwords/internal/adapters/grpc/proto"
	accountStorage "passwords/internal/domain/account"
	userStorage "passwords/internal/domain/user"
	"passwords/internal/pkg/common/logging"
)

type GRPCService struct {
	pb.UnimplementedServiceServer

	Users          sync.Map
	UserStorage    userStorage.Storage    //a storage to set/get user data
	AccountStorage accountStorage.Storage //a storage to set/get account data
}

// NewGRPCService - a factory to User gRPC server service, receives used storage implementation
func NewGRPCService(storage userStorage.Storage) pb.ServiceServer {
	return &GRPCService{UserStorage: storage}
}

// AddUser - adds inbound User data to storage
func (s *GRPCService) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	var response pb.AddUserResponse

	//add user data
	user := userStorage.User{
		User:     in.User.User,
		Password: in.User.Password,
		Name:     in.User.Name,
		Surname:  in.User.Surname,
		Phone:    in.User.Phone,
	}
	//save User data
	s.UserStorage.SaveUser(ctx, &user)
	log.Printf("User %v saved through gRPC", user)
	return &response, nil
}

// GetAccount - gets user data by inbound user name from storage
func (s *GRPCService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var response pb.GetUserResponse
	//get User data
	u, err := s.UserStorage.GetUser(ctx, in.Login)
	logging.LogPrintln(err)

	log.Printf("User %v gotten through gRPC", u)

	response.User.Password = u.Password
	response.User.Name = u.Name
	response.User.Surname = u.Surname
	response.User.Phone = u.Phone

	return &response, err

}

// AddAccount - adds inbound Account data to storage
func (s *GRPCService) AddAccount(ctx context.Context, in *pb.AddAccountRequest) (*pb.AddAccountResponse, error) {
	var response pb.AddAccountResponse

	//add Account data
	Account := accountStorage.Account{
		Account:     in.Account.Account,
		User:        in.User,
		Password:    in.Account.Password,
		Description: in.Account.Description,
	}
	//save Account data
	err := s.AccountStorage.SaveAccount(ctx, Account)
	logging.LogPrintln(err)

	log.Printf("Account %v saved through gRPC", Account)
	return &response, err
}

// GetAccount - gets account data by inbound user name from storage
func (s *GRPCService) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	var response pb.GetAccountResponse
	//get Account data
	a, err := s.AccountStorage.GetAccount(ctx, in.User)
	logging.LogPrintln(err)

	log.Printf("Account %v gotten through gRPC", a)

	response.Account.User = a.User
	response.Account.AccountLogin = a.Account
	response.Account.Password = a.Password
	response.Account.Description = a.Description

	return &response, err

}

func (s *GRPCService) GetAllUserAccounts(ctx context.Context, in *pb.GetAllAccountsRequest) (*pb.GetAllAccountsResponse, error){
	var response pb.GetAllAccountsResponse

	//Get all accounts by given user name
	accounts,err:= s.AccountStorage.GetAccountsList(ctx,in.UserLogin)
	logging.LogPrintln(err)

	
	for k,v := range accounts{
		pbAccount := pb.Account{
			User:v.User,
			Account:v.Account,
			AccountLogin: v.Account,


		} 
		response.Accounts=append (response.Accounts,&pbAccount)
	}
	 
	return &response,err
}