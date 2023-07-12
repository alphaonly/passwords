package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"passwords/internal/pkg/common/logging"
	"passwords/internal/pkg/dbclient"
	"passwords/internal/pkg/dbclient/postgres"
	"passwords/internal/schema"
	"strings"
	"time"

	"passwords/internal/domain/user"
)

type userStorage struct {
	client dbclient.DBClient
}

func NewStorage(client dbclient.DBClient) user.Storage {
	return &userStorage{client: client}

}

func (s userStorage) GetUser(ctx context.Context, name string) (u *user.User, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	d := DBUsersDTO{userID: sql.NullString{String: name, Valid: true}}
	row := conn.QueryRow(ctx, selectLineUsersTable, &d.userID)
	err = row.Scan(&d.userID, &d.password, &d.name, &d.surname, &d.phone, &d.created)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, nil
		}
		return nil, err
	}

	created, err := time.Parse(time.RFC3339, d.created.String)
	logging.LogFatal(err)

	return &user.User{
		User:     d.userID.String,
		Password: d.password.String,
		Name:     d.name.String,
		Surname:  d.surname.String,
		Phone:    d.phone.String,
		Created:  schema.CreatedTime(created),
	}, nil
}

func (s userStorage) SaveUser(ctx context.Context, u *user.User) (err error) {
	if !s.client.Connect(ctx) {
		return errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	d := DBUsersDTO{
		userID:   sql.NullString{String: u.User, Valid: true},
		password: sql.NullString{String: u.Password, Valid: true},
		name:     sql.NullString{String: u.Name, Valid: true},
		surname:  sql.NullString{String: u.Surname, Valid: true},
		phone:    sql.NullString{String: u.Phone, Valid: true},
		created:  sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
	}

	tag, err := conn.Exec(ctx, createOrUpdateIfExistsUsersTable, d.userID, d.password, d.name, d.surname, d.phone, d.created)
	logging.LogFatalf(postgres.Message[3], err)
	log.Println(tag)
	return err
}
