package domain

type PasswordTool interface {
	PasswordGenerator
	PasswordValidator
}

type PasswordGenerator interface {
	New(password string) (string, error)
}

type PasswordValidator interface {
	IsValid(password string) bool
}
