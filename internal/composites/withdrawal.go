package composites

import (
	withdrawalh "github.com/alphaonly/multipass/internal/adapters/api/withdrawal"
	withdrawald "github.com/alphaonly/multipass/internal/adapters/db/withdrawal"
	"github.com/alphaonly/multipass/internal/configuration"
	"github.com/alphaonly/multipass/internal/domain/order"
	"github.com/alphaonly/multipass/internal/domain/user"
	"github.com/alphaonly/multipass/internal/domain/withdrawal"
	"github.com/alphaonly/multipass/internal/pkg/dbclient"
)

type WithdrawalComposite struct {
	Storage withdrawal.Storage
	Service withdrawal.Service
	Handler withdrawalh.Handler
}

func NewWithdrawalComposite(
	dbClient dbclient.DBClient,
	configuration *configuration.ServerConfiguration,
	userStorage user.Storage,
	orderService order.Service,
) *WithdrawalComposite {
	storage := withdrawald.NewStorage(dbClient)
	service := withdrawal.NewService(storage, userStorage, orderService)
	handler := withdrawalh.NewHandler(storage, service, orderService, configuration)
	return &WithdrawalComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
