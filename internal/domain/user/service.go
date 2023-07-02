package user

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrUserPassEmpty  = fmt.Errorf("400 user or password is empty")
	ErrInternal       = fmt.Errorf("500 internal error: ")
	ErrLoginOccupied  = fmt.Errorf("409 login is occupied")
	ErrSaveUser       = fmt.Errorf("500 cannot save user in storage")
	ErrLogPassUnknown = errors.New("401 login or password is unknown")
)

type Service interface {
	RegisterUser(ctx context.Context, u *User) (err error)
	AuthenticateUser(ctx context.Context, u *User) (err error)
	CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error)
}

type service struct {
	Storage Storage
}

func NewService(s Storage) (sr Service) {
	return &service{Storage: s}
}

func (sr service) RegisterUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return ErrUserPassEmpty
	}
	// Check if username exists
	userChk, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", u.User, ErrInternal)
		ErrInternal = fmt.Errorf("500 internal error in getting user %v: %w", u.User, ErrInternal)
		return ErrInternal
	}
	if userChk != nil {
		//login has already been occupied
		ErrLoginOccupied = fmt.Errorf("409 login %v is occupied (%w)", userChk.User, ErrLoginOccupied)
		return ErrLoginOccupied
	}
	err = sr.Storage.SaveUser(ctx, u)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", ErrInternal)
		ErrInternal = fmt.Errorf(" 500 cannot save user in storage %w", ErrInternal)
		return ErrInternal
	}
	return nil
}

func (sr service) AuthenticateUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return ErrUserPassEmpty
	}
	// Check if username exists
	userInStorage, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", u.User, ErrInternal)
		ErrInternal = fmt.Errorf("500 internal error in getting user %v: %w", u.User, ErrInternal)
		log.Println(ErrInternal)
		return ErrInternal
	}
	if !u.Equals(userInStorage) {
		return ErrLogPassUnknown
	}

	return nil
}

func (sr service) CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error) {
	// data validation
	if login == "" || password == "" {
		return false, ErrUserPassEmpty
	}
	// Check if username authorized
	u, err := sr.Storage.GetUser(ctx, login)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", ErrInternal)
		ErrInternal = fmt.Errorf("500 checking user authorization, can not get user from storage: %w", ErrInternal)
		log.Println(ErrInternal)
		return false, ErrInternal
	}
	if u == nil {
		ErrLogPassUnknown = fmt.Errorf("401 no user in storage means not authorized(%w)", ErrLogPassUnknown)
		return false, ErrLogPassUnknown
	}
	if !u.Equals(&User{User: login, Password: password}) {
		return false, nil
	}

	return true, nil
}
