package uuidv4

import "github.com/google/uuid"

type UUIDv4Generator struct{}

func NewGenerator() *UUIDv4Generator {
	return &UUIDv4Generator{}
}

func (gen *UUIDv4Generator) New() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (gen *UUIDv4Generator) IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
