package composites

import (
	"github.com/alphaonly/multipass/internal/adapters/api"
	"github.com/alphaonly/multipass/internal/adapters/api/order"
	"github.com/alphaonly/multipass/internal/adapters/api/user"
	"github.com/alphaonly/multipass/internal/adapters/api/withdrawal"
)

type HandlerComposite struct {
	Common     api.Handler
	User       user.Handler
	Order      order.Handler
	Withdrawal withdrawal.Handler
}

func NewHandlerComposite(
	common api.Handler,
	user user.Handler,
	order order.Handler,
	withdrawal withdrawal.Handler) *HandlerComposite {

	return &HandlerComposite{
		Common:     common,
		User:       user,
		Order:      order,
		Withdrawal: withdrawal,
	}
}
