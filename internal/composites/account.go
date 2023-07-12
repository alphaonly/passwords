package composites

import (
	accountd "passwords/internal/adapters/db/account"
	"passwords/internal/domain/account"
	"passwords/internal/pkg/dbclient"
)

type OrderComposite struct {
	Storage account.Storage
	Service account.Service
}

func NewAccountComposite(dbClient dbclient.DBClient) *OrderComposite {
	storage := accountd.NewStorage(dbClient)
	service := account.NewService(storage)
	return &OrderComposite{
		Storage: storage,
		Service: service,
	}
}
