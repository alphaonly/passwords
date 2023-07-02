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
	"passwords/internal/schema"
)

type accountStorage struct {
	client dbclient.DBClient
}

func NewStorage(client dbclient.DBClient) account.Storage {
	return &accountStorage{client: client}
}

func (s accountStorage) GetAccount(ctx context.Context, name string) (acc *account.Account, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	d := DBAccountsDTO{AccountID: sql.NullString{String: name, Valid: true}}
	row := conn.QueryRow(ctx, selectLineAccountsTable, &d.AccountID)
	err = row.Scan(&d.AccountID, &d.userID, &d.password, &d.descr, &d.createdAt)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, err
	}
	created, err := time.Parse(time.RFC3339, d.createdAt.String)

	return &account.Account{
		Account:     d.AccountID.String,
		User:        d.userID.String,
		Password:    d.password.String,
		Description: d.descr.String,
		Created:     schema.CreatedTime(created),
	}, nil
}
func (s accountStorage) SaveAccount(ctx context.Context, a account.Account) (err error) {
	if !s.client.Connect(ctx) {
		return errors.New(postgres.Message[0])
	}

	accountName := a.Account

	d := &DBAccountsDTO{
		AccountID: sql.NullString{String: accountName, Valid: true},
		userID:    sql.NullString{String: a.User, Valid: true},
		password:  sql.NullString{String: a.Password, Valid: true},
		descr:     sql.NullString{String: a.Description, Valid: true},
		createdAt: sql.NullString{String: time.Time(a.Created).Format(time.RFC3339), Valid: true},
	}

	conn, err := s.client.GetConn()
	if err != nil {
		return err
	}
	tag, err := conn.Exec(ctx, createOrUpdateIfExistsAccountsTable, d.AccountID, d.userID, d.password, d.descr, d.createdAt)
	logging.LogFatalf(postgres.Message[3], err)
	log.Println(tag)
	return err
}

func (s accountStorage) GetAccountsList(ctx context.Context, userName string) (al account.Accounts, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	al = make(account.Accounts)

	d := DBAccountsDTO{userID: sql.NullString{String: userName, Valid: true}}

	rows, err := conn.Query(ctx, selectAllAccountsTableByUser, &d.userID)
	if err != nil {
		log.Printf(postgres.Message[4], err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.AccountID, &d.userID, &d.password, &d.descr, &d.createdAt)
		logging.LogFatalf(postgres.Message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logging.LogFatalf(postgres.Message[6], err)
		al[d.AccountID.String] = account.Account{
			Account:     d.AccountID.String,
			User:        d.userID.String,
			Password:    d.password.String,
			Description: d.descr.String,
			Created:     schema.CreatedTime(created),
		}
	}

	return al, nil
}
