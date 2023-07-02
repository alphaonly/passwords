package postgres

import (
	"context"
	"fmt"
	"github.com/alphaonly/multipass/internal/pkg/common/logging"
	"github.com/alphaonly/multipass/internal/pkg/dbclient"
	"log"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresClient struct {
	dataBaseURL string
	pool        *pgxpool.Pool
	conn        *pgxpool.Conn
}

func (pc postgresClient) GetConn() (*pgxpool.Conn, error) {
	// func (pc postgresClient) GetConn() (*pgxpool.Conn, error) {
	if reflect.ValueOf(pc.conn).IsNil() {
		// if pc.conn == nil {
		return nil, fmt.Errorf(Message[8])
	}
	return pc.conn, nil
}

func (pc postgresClient) GetPull() (*pgxpool.Pool, error) {
	if reflect.ValueOf(pc.pool).IsNil() {
		// if pc.pool == nil {
		return nil, fmt.Errorf(Message[9])
	}
	return pc.pool, nil
}

func NewPostgresClient(ctx context.Context, dataBaseURL string) dbclient.DBClient {
	//get params
	pc := postgresClient{dataBaseURL: dataBaseURL}
	//connect db
	var err error

	pc.pool, err = pgxpool.New(ctx, pc.dataBaseURL)
	if err != nil {
		logging.LogFatalf(Message[0], err)
		return nil
	}

	err = pc.checkTables(ctx)
	logging.LogFatalf(Message[10], err)

	return &pc
}

func (pc *postgresClient) Connect(ctx context.Context) (ok bool) {
	ok = false
	var err error

	if reflect.TypeOf(pc.pool) == nil {

		pc.pool, err = pgxpool.New(ctx, pc.dataBaseURL)
		logging.LogFatalf(Message[0], err)
	}
	for i := 0; i < 10; i++ {
		pc.conn, err = pc.pool.Acquire(ctx)

		if err != nil {
			log.Println(Message[12] + " " + err.Error())
			time.Sleep(time.Millisecond * 200)
			continue
		}
		break
	}

	err = pc.conn.Ping(ctx)
	if err != nil {
		logging.LogFatalf(Message[0], err)
	}

	ok = true
	return ok
}

func (pc postgresClient) checkTables(ctx context.Context) error {
	if reflect.TypeOf(pc.conn) == nil {
		return fmt.Errorf(Message[9])
	}
	var err error
	// check users table exists
	err = CreateTable(ctx, pc, checkIfUsersTableExists, createUsersTable)
	logging.LogFatalf("error:", err)
	// check orders table exists
	err = CreateTable(ctx, pc, checkIfOrdersTableExists, createOrdersTable)
	logging.LogFatalf("error:", err)
	// check withdrawals table exists
	err = CreateTable(ctx, pc, checkIfWithdrawalsTableExists, createWithdrawalsTable)
	logging.LogFatalf("error:", err)

	return nil
}
