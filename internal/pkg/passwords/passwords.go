package passwords

import (
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
