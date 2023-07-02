package accrual

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/alphaonly/multipass/internal/configuration"
	"github.com/alphaonly/multipass/internal/domain/order"
	"github.com/alphaonly/multipass/internal/domain/user"
)

//Periodically checking orders' accrual from remote service

type Accrual interface {
	Run(ctx context.Context)
}

type accrual struct {
	serviceAddress string
	requestTime    time.Duration //200 * time.Millisecond
	OrderStorage   order.Storage
	UserStorage    user.Storage
}

func NewAccrual(configuration *configuration.ServerConfiguration, orderStorage order.Storage, userStorage user.Storage) Accrual {
	return &accrual{
		serviceAddress: configuration.AccrualSystemAddress,
		requestTime:    time.Duration(configuration.AccrualTime) * time.Millisecond,
		OrderStorage:   orderStorage,
		UserStorage:    userStorage,
	}
}

func (acr accrual) Run(ctx context.Context) {

	ticker := time.NewTicker(acr.requestTime)
	baseURL, err := url.Parse(acr.serviceAddress)
	if err != nil {
		log.Fatal("unable to parse URL for accrual system")
	}

	httpc := resty.New().
		SetBaseURL(baseURL.String())

doItAGain:
	for {
		select {
		case <-ticker.C:
			//Getting New unprocessed orders to make a request to accrual system
			oList, err := acr.OrderStorage.GetNewOrdersList(ctx)
			if err != nil {
				log.Fatal("can not get new orders list")
			}

			for orderNumber, orderData := range oList {

				orderNumberStr := strconv.Itoa(int(orderNumber))
				req := httpc.R().
					SetHeader("Accept", "application/json")

				response := order.OrderAccrualResponse{}
				resp, err := req.
					SetResult(&response).
					Get("api/orders/" + orderNumberStr)
				if err != nil {
					log.Printf("account %v response error: %v", orderNumber, err)
					continue
				}
				log.Printf("account %v response from accrual: %v", orderNumber, resp)
				if response.Status != order.ProcessedOrder.Text {
					log.Printf("account %v response status type %v, continue", orderNumber, resp.Status())
					continue
				}

				orderData.Accrual = response.Accrual
				orderData.Status = order.ProcessedOrder.Text
				log.Printf("Saving processed account:%v", orderData)

				err = acr.OrderStorage.SaveOrder(ctx, orderData)
				if err != nil {
					log.Fatal("unable to save account")
				}
				log.Printf("Processed account saved from accrual:%v", orderData)
				log.Printf("Update user balance with processed account:%v", orderNumber)
				//Update balance in case of account accrual greater than zero
				if orderData.Accrual > 0 {

					u, err := acr.UserStorage.GetUser(ctx, orderData.User)
					if err != nil {
						log.Fatalf("Error in getting user %v data: %v", orderData.User, err.Error())
					}
					if u == nil {
						log.Fatalf("Data inconsistency with there is no user %v, but there is account %v with the user", orderData.User, orderNumber)
					}
					u.Accrual += orderData.Accrual
					err = acr.UserStorage.SaveUser(ctx, u)
					if err != nil {
						log.Fatalf("Unable to save user %v with updated accrual %v: %v", u.User, u.Accrual, err.Error())
					}
					log.Printf("Updated user:%v", u)

				}
			}

		case <-ctx.Done():
			break doItAGain
		}
	}

}

func (acr accrual) sendRequest(ctx context.Context) {}

func (acr accrual) GetResponse(ctx context.Context) {}
