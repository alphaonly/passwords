package composites

import (
	
	userd "passwords/internal/adapters/db/user"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/domain/user"
	"passwords/internal/pkg/dbclient"
)

type UserComposite struct {
	Storage user.Storage
	Service user.Service
}

func NewUserComposite(dbClient dbclient.DBClient, configuration *configuration.ServerConfiguration) *UserComposite {
	storage := userd.NewStorage(dbClient)
	service := user.NewService(storage)
	return &UserComposite{
		Storage: storage,
		Service: service,
	}
}
