package account

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"passwords/internal/pkg/common/logging"
	"passwords/internal/pkg/dbclient"
	"passwords/internal/pkg/dbclient/postgres"
	"time"

	"passwords/internal/domain/account"
	accountDomain "passwords/internal/domain/account"
	"passwords/internal/schema"
)

type accountStorage struct {
	client dbclient.DBClient
}

func NewStorage(client dbclient.DBClient) account.Storage {
	return &accountStorage{client: client}
}

func (s accountStorage) GetAccount(ctx context.Context, user string, account string) (acc *account.Account, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	d := DBAccountsDTO{
		UserID:    sql.NullString{String: user, Valid: true},
		AccountID: sql.NullString{String: account, Valid: true},
	}
	row := conn.QueryRow(ctx, selectLineAccountsTable, &d.UserID, &d.AccountID)
	err = row.Scan(&d.UserID, &d.AccountID, &d.Login, &d.Password, &d.Descr, &d.CreatedAt)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, err
	}
	created, err := time.Parse(time.RFC3339, d.CreatedAt.String)

	return &accountDomain.Account{
		User:        d.UserID.String,
		Account:     d.AccountID.String,
		Login:       d.Login.String,
		Password:    d.Password.String,
		Description: d.Descr.String,
		Created:     schema.CreatedTime(created),
	}, nil
}
func (s accountStorage) SaveAccount(ctx context.Context, a accountDomain.Account) (err error) {
	if !s.client.Connect(ctx) {
		return errors.New(postgres.Message[0])
	}

	d := &DBAccountsDTO{
		UserID:    sql.NullString{String: a.User, Valid: true},
		AccountID: sql.NullString{String: a.Account, Valid: true},
		Login:     sql.NullString{String: a.Login, Valid: true},
		Password:  sql.NullString{String: a.Password, Valid: true},
		Descr:     sql.NullString{String: a.Description, Valid: true},
		CreatedAt: sql.NullString{String: time.Time(a.Created).Format(time.RFC3339), Valid: true},
	}

	conn, err := s.client.GetConn()
	if err != nil {
		return err
	}
	tag, err := conn.Exec(ctx, createOrUpdateIfExistsAccountsTable, d.UserID, d.AccountID, d.Login, d.Password, d.Descr, d.CreatedAt)
	logging.LogFatalf(postgres.Message[3], err)
	log.Println(tag)
	return err
}

func (s accountStorage) GetAccountsList(ctx context.Context, user string) (al accountDomain.Accounts, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	al = make(accountDomain.Accounts)

	d := DBAccountsDTO{
		UserID: sql.NullString{String: user, Valid: true},
	}

	rows, err := conn.Query(ctx, selectAllAccountsTableByUser, &d.UserID)
	if err != nil {
		log.Printf(postgres.Message[4], err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.UserID, &d.AccountID, &d.Login, &d.Password, &d.Descr, &d.CreatedAt)
		logging.LogFatalf(postgres.Message[5], err)
		created, err := time.Parse(time.RFC3339, d.CreatedAt.String)
		logging.LogFatalf(postgres.Message[6], err)
		al[d.AccountID.String] = accountDomain.Account{
			User:        d.UserID.String,
			Account:     d.AccountID.String,
			Login:       d.Login.String,
			Password:    d.Password.String,
			Description: d.Descr.String,
			Created:     schema.CreatedTime(created),
		}
	}

	return al, nil
}
