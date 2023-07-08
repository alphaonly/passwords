package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	pb "passwords/internal/adapters/grpc/proto"
	accountStorage "passwords/internal/domain/account"
	userStorage "passwords/internal/domain/user"
)

var (
	ErrUserNotAuthorized = errors.New("user %v not authorized")
	ErrUserAuthorization = errors.New("user %v authorization error %w")
)

type GRPCService struct {
	pb.UnimplementedServiceServer

	AuthorizedUsers sync.Map               // keep information about users that are currently authorized
	userStorage     userStorage.Storage    //a storage to set/get user data
	accountStorage  accountStorage.Storage //a storage to set/get account data
}

// NewGRPCService - a factory to User gRPC server service, receives used storage implementation
func NewGRPCService(userStor userStorage.Storage, accountStor accountStorage.Storage) pb.ServiceServer {
	return &GRPCService{
		userStorage:    userStor,
		accountStorage: accountStor,
	}
}

// AuthorizeUser - checks user and password and authorize  inbound User data in server
func (s *GRPCService) AuthorizeUser(ctx context.Context, in *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
	var response pb.AuthUserResponse

	//get User data
	user, err := s.userStorage.GetUser(ctx, in.User)
	if err != nil {
		response.Error = ErrUserNotAuthorized.Error()
		return &response, fmt.Errorf(ErrUserAuthorization.Error(), in.User, err)
	}
	if user == nil || user.Password != in.Password {
		response.Error = ErrUserNotAuthorized.Error()
		return &response, ErrUserNotAuthorized
	}
	response.User = &pb.User{
		User:     user.User,
		Password: user.Password,
		Name:     user.Name,
		Surname:  user.Surname,
		Phone:    user.Phone,
	}

	log.Printf("User %v authorized on client through gRPC", user)
	return &response, nil
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
	err := s.userStorage.SaveUser(ctx, &user)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}

	log.Printf("User %v saved through gRPC", user)
	return &response, nil
}

// GetAccount - gets user data by inbound user name from storage
func (s *GRPCService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var response pb.GetUserResponse
	//get User data
	u, err := s.userStorage.GetUser(ctx, in.Login)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, err
	}
	response.User = &pb.User{
		User:     u.User,
		Password: u.Password,
		Name:     u.Name,
		Surname:  u.Surname,
		Phone:    u.Phone,
		Created:  time.Time(u.Created).Format(time.RFC3339),
	}
	log.Printf("User %v gotten through gRPC", u)
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
	err := s.accountStorage.SaveAccount(ctx, Account)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, err
	}

	log.Printf("Account %v saved through gRPC", Account)
	return &response, err
}

// GetAccount - gets account data by inbound user name from storage
func (s *GRPCService) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	var response pb.GetAccountResponse
	//get Account data
	a, err := s.accountStorage.GetAccount(ctx, in.User, in.Account)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, err
	}
	response.Account = &pb.Account{
		User:         a.User,
		Account:      a.Account,
		AccountLogin: a.Login,
		Password:     a.Password,
		Description:  a.Description,
		Created:      time.Time(a.Created).Format(time.RFC3339),
	}

	log.Printf("Account %v has been gotten through gRPC", a)
	return &response, err

}

// GetAllUserAccounts - gets all accounts by given user name
func (s *GRPCService) GetAllUserAccounts(ctx context.Context, in *pb.GetAllAccountsRequest) (*pb.GetAllAccountsResponse, error) {
	var response pb.GetAllAccountsResponse

	//Get all accounts by given username for storage
	accounts, err := s.accountStorage.GetAccountsList(ctx, in.UserLogin)
	if err != nil {
		log.Println(err)
		response.Error = err.Error()
		return &response, err
	}
	for _, v := range accounts {
		pbAccount := pb.Account{
			User:         v.User,
			Account:      v.Account,
			AccountLogin: v.Login,
			Password:     v.Password,
			Description:  v.Description,
			Created:      time.Time(v.Created).Format(time.RFC3339),
		}
		response.Accounts = append(response.Accounts, &pbAccount)
	}
	return &response, err
}
