package passwords

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var saltedHashRegex = regexp.MustCompile(`^\$2[ayb]\$.{56}$`)

type PasswordGenerator struct {
	cost int
}

func NewGenerator(cost int) *PasswordGenerator {
	return &PasswordGenerator{cost: cost}
}

func (pg *PasswordGenerator) New(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), pg.cost)
	if err != nil {
		return "", err
	}

	return string(p), nil
}

func (pg *PasswordGenerator) IsValid(password string) bool {
	return !saltedHashRegex.MatchString(password)
}

func (pg *PasswordGenerator) Check(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		return fmt.Errorf("mismatched hash and password")
	} else if err != nil {
		return fmt.Errorf("encountered error checking password")
	}

	return nil
}
