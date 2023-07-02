package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alphaonly/multipass/internal/adapters/api"
	"github.com/alphaonly/multipass/internal/configuration"
	"github.com/alphaonly/multipass/internal/domain/order"
	"github.com/alphaonly/multipass/internal/domain/user"
	"github.com/alphaonly/multipass/internal/schema"
)

type Handler interface {
	PostOrders() http.HandlerFunc
	GetOrders() http.HandlerFunc
	GetBalance() http.HandlerFunc
}
type handler struct {
	Storage       order.Storage
	Service       order.Service
	UserService   user.Service
	Configuration *configuration.ServerConfiguration
}

func NewHandler(storage order.Storage, service order.Service, userService user.Service, configuration *configuration.ServerConfiguration) Handler {
	return &handler{
		Storage:       storage,
		Service:       service,
		UserService:   userService,
		Configuration: configuration,
	}
}

func (h *handler) PostOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostUserOrders invoked")
		//Get parameters from previous handler
		usr, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HTTPError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		OrderNumberByte, err := io.ReadAll(r.Body)
		if err != nil {
			api.HTTPError(w, fmt.Errorf("unrecognized body body %w", err), http.StatusBadRequest)
			return
		}

		orderNumber, err := h.Service.ValidateOrderNumber(r.Context(), string(OrderNumberByte), string(usr))
		if err != nil {
			if errors.Is(err, order.ErrBadUserOrOrder) {
				api.HTTPErrorW(w, fmt.Sprintf("account number  %v insufficient format", orderNumber), err, http.StatusBadRequest)
				return
			}

			if errors.Is(err, order.ErrNoLuhnNumber) {
				api.HTTPErrorW(w, fmt.Sprintf("account %v insufficient format", orderNumber), err, http.StatusUnprocessableEntity)
				return
			}
			if strings.Contains(err.Error(), "409") {
				if errors.Is(err, order.ErrAnotherUsersOrder) {
				}
				api.HTTPErrorW(w, fmt.Sprintf("account %v exists", orderNumber), err, http.StatusConflict)
				return
			}
			if errors.Is(err, order.ErrOrderNumberExists) {
			}
			if strings.Contains(err.Error(), "200") {
				log.Printf("account %v exists: %v", orderNumber, err.Error())
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		//Create object for a new account
		o := order.Order{
			Order:   string(OrderNumberByte),
			User:    string(usr),
			Status:  order.NewOrder.Text,
			Created: schema.CreatedTime(time.Now()),
		}
		err = h.Storage.SaveOrder(r.Context(), o)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("account's number %v not saved", orderNumber), err, http.StatusInternalServerError)
			return
		}
		//Response
		w.WriteHeader(http.StatusAccepted)
	}
}
func (h *handler) GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetUserOrders invoked")

		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HTTPError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		orderList, err := h.Service.GetUsersOrders(r.Context(), string(userName))
		if err != nil {
			if strings.Contains(err.Error(), "204") {
				api.HTTPErrorW(w, fmt.Sprintf("No orders for user %v", userName), err, http.StatusNoContent)
				return
			}
		}
		//Response
		bytes, err := json.Marshal(orderList)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("user %v account list json marshal error", userName), err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("user %v HandleGetUserOrders write response error", userName), err, http.StatusInternalServerError)
			return
		}
	}
}
func (h *handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetUserBalance invoked")
		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HTTPError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		balance, err := h.UserService.GetUserBalance(r.Context(), string(userName))
		if err != nil {
			api.HTTPError(w, fmt.Errorf("cannot get user data by userName %v from context %w", userName, err), http.StatusInternalServerError)
			return
		}
		log.Printf("Got balance %v for user %v ", balance, userName)
		//Response
		bytes, err := json.Marshal(balance)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("user %v balance json marshal error", userName), err, http.StatusInternalServerError)
			return
		}
		log.Printf("Write response balance json:%v ", string(bytes))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("user %v balance write response error", userName), err, http.StatusInternalServerError)
			return
		}
	}
}
