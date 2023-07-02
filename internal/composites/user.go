package composites

import (
	userh "github.com/alphaonly/multipass/internal/adapters/api/user"
	userd "github.com/alphaonly/multipass/internal/adapters/db/user"
	"github.com/alphaonly/multipass/internal/configuration"
	"github.com/alphaonly/multipass/internal/domain/user"
	"github.com/alphaonly/multipass/internal/pkg/dbclient"
)

type UserComposite struct {
	Storage user.Storage
	Service user.Service
	Handler userh.Handler
}

func NewUserComposite(dbClient dbclient.DBClient, configuration *configuration.ServerConfiguration) *UserComposite {
	storage := userd.NewStorage(dbClient)
	service := user.NewService(storage)
	handler := userh.NewHandler(storage, service, configuration)
	return &UserComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
