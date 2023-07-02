package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alphaonly/multipass/internal/pkg/common/logging"
	"github.com/alphaonly/multipass/internal/pkg/dbclient"
	"github.com/alphaonly/multipass/internal/pkg/dbclient/postgres"
	"log"
	"strconv"
	"time"

	"github.com/alphaonly/multipass/internal/domain/order"
	"github.com/alphaonly/multipass/internal/schema"
)

type orderStorage struct {
	client dbclient.DBClient
}

func NewStorage(client dbclient.DBClient) order.Storage {
	return &orderStorage{client: client}

}

func (s orderStorage) GetOrder(ctx context.Context, orderNumber int64) (o *order.Order, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	d := DBOrdersDTO{orderID: sql.NullInt64{Int64: orderNumber, Valid: true}}
	row := conn.QueryRow(ctx, selectLineOrdersTable, &d.orderID)
	err = row.Scan(&d.orderID, &d.userID, &d.status, &d.accrual, &d.createdAt)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, err
	}
	created, err := time.Parse(time.RFC3339, d.createdAt.String)

	return &order.Order{
		Order:   strconv.FormatInt(d.orderID.Int64, 10),
		User:    d.userID.String,
		Status:  order.OrderTypesByCode[d.status.Int64].Text,
		Accrual: d.accrual.Float64,
		Created: schema.CreatedTime(created),
	}, nil
}
func (s orderStorage) SaveOrder(ctx context.Context, o order.Order) (err error) {
	if !s.client.Connect(ctx) {
		return errors.New(postgres.Message[0])
	}

	orderInt, err := strconv.ParseInt(o.Order, 10, 64)
	if err != nil {
		log.Fatal(fmt.Errorf("error in converting account number %v to string:%w", o.Order, err))
	}

	d := &DBOrdersDTO{
		orderID:   sql.NullInt64{Int64: orderInt, Valid: true},
		userID:    sql.NullString{String: o.User, Valid: true},
		status:    sql.NullInt64{Int64: order.OrderTypesByText[o.Status].Code, Valid: true},
		accrual:   sql.NullFloat64{Float64: o.Accrual, Valid: true},
		createdAt: sql.NullString{String: time.Time(o.Created).Format(time.RFC3339), Valid: true},
	}

	conn, err := s.client.GetConn()
	if err != nil {
		return err
	}
	tag, err := conn.Exec(ctx, createOrUpdateIfExistsOrdersTable, d.orderID, d.userID, d.status, d.accrual, d.createdAt)
	logging.LogFatalf(postgres.Message[3], err)
	log.Println(tag)
	return err
}

func (s orderStorage) GetOrdersList(ctx context.Context, userName string) (ol order.Orders, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	ol = make(order.Orders)

	d := DBOrdersDTO{userID: sql.NullString{String: userName, Valid: true}}

	rows, err := conn.Query(ctx, selectAllOrdersTableByUser, &d.userID)
	if err != nil {
		log.Printf(postgres.Message[4], err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.orderID, &d.userID, &d.status, &d.accrual, &d.createdAt)
		logging.LogFatalf(postgres.Message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logging.LogFatalf(postgres.Message[6], err)
		ol[d.orderID.Int64] = order.Order{
			Order:   strconv.FormatInt(d.orderID.Int64, 10),
			User:    d.userID.String,
			Status:  order.OrderTypesByCode[d.status.Int64].Text,
			Accrual: d.accrual.Float64,
			Created: schema.CreatedTime(created),
		}
	}

	return ol, nil
}

func (s orderStorage) GetNewOrdersList(ctx context.Context) (ol order.Orders, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	logging.LogFatalf("", err)

	defer conn.Release()

	ol = make(order.Orders)

	d := DBOrdersDTO{status: sql.NullInt64{Int64: order.NewOrder.Code, Valid: true}}

	rows, err := conn.Query(ctx, selectAllOrdersTableByStatus, &d.status)
	if err != nil {
		log.Printf(postgres.Message[4], err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.orderID, &d.userID, &d.status, &d.accrual, &d.createdAt)
		logging.LogFatalf(postgres.Message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logging.LogFatalf(postgres.Message[6], err)
		ol[d.orderID.Int64] = order.Order{
			Order:   strconv.FormatInt(d.orderID.Int64, 10),
			User:    d.userID.String,
			Status:  order.OrderTypesByCode[d.status.Int64].Text,
			Accrual: d.accrual.Float64,
			Created: schema.CreatedTime(created),
		}
	}

	return ol, nil
}
