package account

import (
	"encoding/json"
	"passwords/internal/schema"
)

type Account struct {
	Account     string             `json:"account"`
	User        string             `json:"user"`
	Password    string             `json:"password,omitempty"`
	Description string             `json:"description,omitempty"`
	Created     schema.CreatedTime `json:"uploaded_at"`
}

type AType struct {
	Code int64
	Text string
}

type Accounts map[string]Account

func (a Accounts) MarshalJSON() ([]byte, error) {

	oArray := make([]Account, len(a))
	i := 0
	for _, v := range a {
		oArray[i] = Account{
			Account:     v.Account,
			User:        v.User,
			Password:    v.Password,
			Description: v.Description,
			Created:     v.Created,
		}
		i++
	}
	bytes, err := json.Marshal(&oArray)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func (a Accounts) UnmarshalJSON(b []byte) error {
	var oArray []Account
	if err := json.Unmarshal(b, &oArray); err != nil {
		return err
	}
	for _, v := range oArray {
		a[v.Account] = v
	}
	return nil
}
