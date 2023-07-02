package router

import (
	"github.com/alphaonly/multipass/internal/composites"
	"github.com/go-chi/chi"
)

func NewRouter(h *composites.HandlerComposite) chi.Router {

	var (
		post      = h.Common.Post
		get       = h.Common.Get
		basicAuth = h.User.BasicAuth
		health    = get(h.Common.Health())

		register       = post(h.User.Register())
		login          = post(h.User.Login())
		sendOrders     = post(basicAuth(h.Order.PostOrders()))
		withdraw       = post(basicAuth(h.Withdrawal.PostWithdraw()))
		getOrders      = get(basicAuth(h.Order.GetOrders()))
		balance        = get(basicAuth(h.Order.GetBalance()))
		getWithdrawals = get(basicAuth(h.Withdrawal.GetWithdrawals()))
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/check", health)
		r.Post("/api/user/register", register)
		r.Post("/api/user/login", login)
		r.Post("/api/user/orders", sendOrders)
		r.Post("/api/user/balance/withdraw", withdraw)

		r.Get("/api/user/orders", getOrders)
		r.Get("/api/user/balance", balance)
		r.Get("/api/user/withdrawals", getWithdrawals)

		//Mock for accrual system (in case similar addresses) returns +5.3 score
		r.Get("/api/orders/{number}", h.Common.AccrualScore(nil))

	})

	return r
}
