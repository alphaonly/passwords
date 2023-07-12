package account

//the package describes the model of account entity
import (
	"encoding/json"
	"fmt"
	"passwords/internal/schema"
	"time"
)

// Account entity model
type Account struct {
	User        string             `json:"user"`
	Account     string             `json:"account"`
	Login       string             `json:"login"`
	Password    string             `json:"password,omitempty"`
	Description string             `json:"description,omitempty"`
	Created     schema.CreatedTime `json:"uploaded_at"`
}

// String - Implements a function to describe how to print an account instance
func (a Account) String() string {

	return fmt.Sprintf("User: %v\n", a.User) +
		fmt.Sprintf("Account: %v\n", a.Account) +
		fmt.Sprintf("Login: %v\n", a.Login) +
		fmt.Sprintf("Password: %v\n", a.Password) +
		fmt.Sprintf("Description: %v\n", a.Description) +
		fmt.Sprintf("Created timestamp: %v\n", a.Created)

}

type AType struct {
	Code int64
	Text string
}

// Accounts - describes the list of account (particularly a map for fast element access)
type Accounts map[string]Account

// String - Implements a function to describe how to appear the account list while fmt.Print*
func (a Accounts) String() string {
	var result string
	result = "\n Account IDs:\n"
	//print account list as a variety of account IDs
	//header
	result += "ID      \t|" + "Login       \t|" + "Password           \t|" + "Description                             \t|" + "Created time\n"
	for _, v := range a {
		result +=
			v.Account + "\t|" +
				v.Login + "\t|" +
				v.Password + "\t|" +
				v.Description + "\t|" +
				time.Time(v.Created).Format(time.RFC3339) + "\n"

	}
	result = result + "\n"
	return result
}

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
