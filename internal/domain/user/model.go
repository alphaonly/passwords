package user

type User struct {
	User     string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname,omitempty"`
	Phone    string `json:"phone,omitempty"`
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
