package api

import (
	"encoding/json"
	"fmt"
	"github.com/alphaonly/multipass/internal/pkg/common/logging"
	"log"
	"net/http"
	"strconv"

	"github.com/alphaonly/multipass/internal/configuration"
	"github.com/alphaonly/multipass/internal/domain/order"
	"github.com/go-chi/chi"
)

type Handler interface {
	Health() http.HandlerFunc
	Post(next http.Handler) http.HandlerFunc
	Get(next http.Handler) http.HandlerFunc
	AccrualScore(next http.Handler) http.HandlerFunc
}

func NewHandler(
	configuration *configuration.ServerConfiguration) Handler {

	return &handler{
		Configuration: configuration}
}

type handler struct {
	Configuration *configuration.ServerConfiguration
}

func (h *handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("200 OK"))
		logging.LogFatal(err)
	}
}

func (h *handler) Get(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetValidation invoked")
		//Validation
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
			return
		}
		if next != nil {
			//call further handler with context parameters
			next.ServeHTTP(w, r)
			return
		}
	}
}
func (h *handler) Post(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostValidation invoked")
		//Validation
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
			return
		}
		if next != nil {
			//call further handler with context parameters
			next.ServeHTTP(w, r)
			return
		}
	}
}

func (h *handler) AccrualScore(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetOrderAccrual invoked")
		//Handling
		orderNumberStr := chi.URLParam(r, "number")
		if orderNumberStr == "" {
			HTTPError(w, fmt.Errorf("Account number  %v is empty", orderNumberStr), http.StatusBadRequest)
			return
		}

		_, err := strconv.ParseInt(orderNumberStr, 10, 64)
		if err != nil {
			HTTPError(w, fmt.Errorf("Account number  %v is bad format", orderNumberStr), http.StatusBadRequest)
			return
		}

		accrual := 5.3

		OrderAccrualResponse := order.OrderAccrualResponse{
			Order:   orderNumberStr,
			Status:  "PROCESSED",
			Accrual: accrual,
		}

		//Response
		bytes, err := json.Marshal(&OrderAccrualResponse)
		if err != nil {
			HTTPErrorW(w, "Account accrual response json marshal error", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "apilication/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			HTTPErrorW(w, "Account accrual response write response error", err, http.StatusInternalServerError)
			return
		}

	}
}
