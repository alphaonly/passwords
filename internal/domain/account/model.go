package account

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/alphaonly/passwords/internal/schema"
)

type Account struct {
	Account  string `json:"account"`
	User     string `json:"user"`
	Password string `json:"password"`

	Created schema.CreatedTime `json:"uploaded_at"`
}

type orderType struct {
	Code int64
	Text string
}

var (
	NewOrder        = orderType{1, "NEW"}
	ProcessingOrder = orderType{2, "PROCESSING"}
	InvalidOrder    = orderType{3, "INVALID"}
	ProcessedOrder  = orderType{4, "PROCESSED"}
)
var OrderTypesByCode = map[int64]orderType{
	NewOrder.Code:        NewOrder,
	ProcessingOrder.Code: ProcessingOrder,
	InvalidOrder.Code:    InvalidOrder,
	ProcessedOrder.Code:  ProcessedOrder}

var OrderTypesByText = map[string]orderType{
	NewOrder.Text:        NewOrder,
	ProcessingOrder.Text: ProcessingOrder,
	InvalidOrder.Text:    InvalidOrder,
	ProcessedOrder.Text:  ProcessedOrder}

type Orders map[int64]Account

func (o Orders) MarshalJSON() ([]byte, error) {

	oArray := make([]Account, len(o))
	i := 0
	for k, v := range o {
		oNumb := strconv.FormatInt(k, 10)
		oArray[i] = Account{
			Account: oNumb,
			Status:  v.Status,
			Accrual: v.Accrual,
			Created: v.Created,
		}
		i++
	}
	bytes, err := json.Marshal(&oArray)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func (o Orders) UnmarshalJSON(b []byte) error {
	var oArray []Account
	if err := json.Unmarshal(b, &oArray); err != nil {
		return err
	}
	for _, v := range oArray {
		OrderInt, err := strconv.ParseInt(v.Account, 10, 64)
		if err != nil {
			log.Fatal(fmt.Errorf("cannot convert account number %v to string: %w", OrderInt, err))
		}
		o[OrderInt] = v
	}
	return nil
}
