package composites

import (
	accountd "passwords/internal/adapters/db/account"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/domain/account"
	"passwords/internal/domain/user"
	"passwords/internal/pkg/dbclient"
)

type OrderComposite struct {
	Storage account.Storage
	Service account.Service
}

func NewAccountComposite(dbClient dbclient.DBClient, userService user.Service, configuration *configuration.ServerConfiguration) *OrderComposite {
	storage := accountd.NewStorage(dbClient)
	service := account.NewService(storage)
	return &OrderComposite{
		Storage: storage,
		Service: service,
	}
}
