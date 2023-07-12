package user

import (
	"fmt"
	"passwords/internal/schema"
	"time"
)

type User struct {
	User     string             `json:"login"`
	Password string             `json:"password"`
	Name     string             `json:"name"`
	Surname  string             `json:"surname,omitempty"`
	Phone    string             `json:"phone,omitempty"`
	Created  schema.CreatedTime `json:"created,omitempty"`
}

func (u User) String() string {

	return fmt.Sprintf("User: %v\n", u.User) +
		fmt.Sprintf("Name: %v\n", checkEmpty(u.Name)) +
		fmt.Sprintf("Surname: %v\n", checkEmpty(u.Surname)) +
		fmt.Sprintf("Phone: %v\n", checkEmpty(u.Phone)) +
		fmt.Sprintf("Created timestamp: %v\n", checkEmpty(time.Time(u.Created).Format(time.RFC3339)))

}

func checkEmpty(s string) string {
	if s == "" {
		return "<empty>"
	}
	return s
}

func (u User) Equals(u2 *User) (ok bool) {
	if u2 == nil {
		return false
	}
	if u.User == u2.User && u.Password == u2.Password {
		return true
	}
	return false
}
