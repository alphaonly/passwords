package composites

import (
	orderh "github.com/alphaonly/multipass/internal/adapters/api/account"
	orderd "github.com/alphaonly/multipass/internal/adapters/db/account"
	"github.com/alphaonly/multipass/internal/configuration"
	"github.com/alphaonly/multipass/internal/domain/order"
	"github.com/alphaonly/multipass/internal/domain/user"
	"github.com/alphaonly/multipass/internal/pkg/dbclient"
)

type OrderComposite struct {
	Storage order.Storage
	Service order.Service
	Handler orderh.Handler
}

func NewOrderComposite(dbClient dbclient.DBClient, userService user.Service, configuration *configuration.ServerConfiguration) *OrderComposite {
	storage := orderd.NewStorage(dbClient)
	service := order.NewService(storage)
	handler := orderh.NewHandler(storage, service, userService, configuration)
	return &OrderComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
