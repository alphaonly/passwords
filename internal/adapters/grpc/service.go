package service

import (
	"context"
	"fmt"
	"log"
	"passwords/internal/client/service"
	"passwords/internal/schema"
	"sync"
	"time"

	pb "passwords/internal/adapters/grpc/proto"
	accountStorage "passwords/internal/domain/account"
	userStorage "passwords/internal/domain/user"
)

type GRPCService struct {
	pb.UnimplementedServiceServer

	AuthorizedUsers   *sync.Map              // keep information about users that are currently authorized
	LastUserOperation *sync.Map              // keep information about user's last operation with server
	userStorage       userStorage.Storage    //a storage to set/get user data
	accountStorage    accountStorage.Storage //a storage to set/get account data
}

// NewGRPCService - a factory to User gRPC server service, receives used storage implementation
func NewGRPCService(
	userStor userStorage.Storage,
	accountStor accountStorage.Storage,
	authUsers *sync.Map,
	lastOperation *sync.Map) pb.ServiceServer {
	return &GRPCService{
		userStorage:       userStor,
		accountStorage:    accountStor,
		AuthorizedUsers:   authUsers,
		LastUserOperation: lastOperation,
	}
}

// AuthorizeUser - checks user and password and authorizes  inbound User data in server
func (s *GRPCService) AuthorizeUser(ctx context.Context, in *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
	var response pb.AuthUserResponse

	//get User data
	user, err := s.userStorage.GetUser(ctx, in.User)
	if err != nil {
		response.Error = service.ErrWrongUserOrPassword.Error()
		return &response, service.ErrWrongUserOrPassword
	}
	//No such a user in storage
	if user == nil {
		switch {
		case in.Password != "": //No user but password is given, it may be created
			{
				response.Error = service.ErrNoUserExists.Error()
				return &response, service.ErrNoUserExists
			}
		case in.Password == "": //No user and password is not given, error
			{
				response.Error = service.ErrWrongUserOrPassword.Error()
				return &response, service.ErrUserWasNotAuthorized
			}
		}
	}
	//Check password
	if user.Password != in.Password {
		response.Error = service.ErrWrongUserOrPassword.Error()
		return &response, service.ErrWrongUserOrPassword
	}
	//check user has been authorized
	_, authorized := s.AuthorizedUsers.Load(user.User)
	if authorized {
		return nil, service.ErrUserAlreadyAuthorized
	}
	//mark user as authorized
	s.AuthorizedUsers.Store(user.User, schema.CreatedTime(time.Now()))
	//remember user's operation(authorization) time
	s.LastUserOperation.Store(user.User, time.Now())
	//response
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

// DisAuthorizeUser - checks user and password and authorize  inbound User data in server
func (s *GRPCService) DisAuthorizeUser(ctx context.Context, in *pb.DisAuthUserRequest) (*pb.DisAuthUserResponse, error) {
	var response pb.DisAuthUserResponse

	//get User data
	user, err := s.userStorage.GetUser(ctx, in.User)
	if err != nil {
		response.Error = service.ErrWrongUserOrPassword.Error()
		return &response, fmt.Errorf(service.ErrWrongUserOrPassword.Error()+"%w", err)
	}
	if user == nil || user.Password != in.Password {
		response.Error = service.ErrWrongUserOrPassword.Error()
		return &response, service.ErrWrongUserOrPassword
	}
	//check user has been authorized
	_, authorized := s.AuthorizedUsers.Load(user.User)
	if !authorized {
		return nil, service.ErrUserWasNotAuthorized
	}
	//mark user as disauthorized
	s.AuthorizedUsers.Delete(user.User)
	//delete user's operation mark
	s.LastUserOperation.Delete(user.User)

	log.Printf("User %v disauthorized in server through gRPC by client demand", user)
	//response
	return &response, nil
}

// CheckAuthorization  - checks current user authorization in server
func (s *GRPCService) CheckAuthorization(ctx context.Context, in *pb.CheckAuthUserRequest) (*pb.CheckAuthUserResponse, error) {
	var response pb.CheckAuthUserResponse

	_, ok := s.AuthorizedUsers.Load(in.User)
	if !ok {
		response.Error = service.ErrUserWasNotAuthorized.Error()
	}
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
	//Authorize user as it is added during first authorization in client
	s.AuthorizedUsers.Store(in.User.User, schema.CreatedTime(time.Now()))
	//remember user's operation(add user - like first time authorized) time
	s.LastUserOperation.Store(in.User.User, time.Now())
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
	//remember user's operation(get user) time
	s.LastUserOperation.Store(u.User, time.Now())

	log.Printf("User %v gotten through gRPC", u)
	return &response, err
}

// AddAccount - adds inbound Account data to storage
func (s *GRPCService) AddAccount(ctx context.Context, in *pb.AddAccountRequest) (*pb.AddAccountResponse, error) {
	var response pb.AddAccountResponse
	//add Account data
	Account := accountStorage.Account{
		User:        in.Account.User,
		Account:     in.Account.Account,
		Login:       in.Account.AccountLogin,
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
	//remember user's operation(add account) time
	s.LastUserOperation.Store(in.User, time.Now())

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

// GetAllUserAccounts - gets all accounts by given user's name
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

	//remember user's operation(get all accounts) time
	s.LastUserOperation.Store(in.UserLogin, time.Now())

	return &response, err
}
