package domain

import (
	"fmt"
	"net/mail"
	"time"
)

type User struct {
	ID         string         `json:"id"`
	Attributes UserAttributes `json:"attributes"`
	Password   string         `json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	ModifiedAt time.Time      `json:"modifiedAt"`
}

type UserAttributes struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

func NewUser(id, userName, email, password string) *User {
	return &User{
		ID: id,
		Attributes: UserAttributes{
			UserName: userName,
			Email:    email,
		},
		Password: password,
	}
}

func (u *User) Validate(idValidator IDValidator) error {
	if !idValidator.IsValid(u.ID) {
		return fmt.Errorf("id format is invalid")
	}

	// UserName
	switch {
	case u.Attributes.UserName == "":
		return fmt.Errorf("user name must not be empty")
	case len(u.Attributes.UserName) > 20:
		return fmt.Errorf("user name must not be longer than 20 characters")
	}

	// Email
	switch {
	case u.Attributes.Email == "":
		return fmt.Errorf("email must not be empty")
	case len(u.Attributes.Email) > 255:
		return fmt.Errorf("email must not be longer than 255 characters")
	}

	if _, err := mail.ParseAddress(u.Attributes.Email); err != nil {
		return fmt.Errorf("email format is invalid: %v", err)
	}

	// Password
	switch {
	case u.Password == "":
		return fmt.Errorf("password must not be empty")
	case len(u.Password) > 255:
		return fmt.Errorf("password must not be longer than 255 characters")
	}

	return nil
}
