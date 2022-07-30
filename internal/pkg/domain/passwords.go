package domain

type PasswordTool interface {
	PasswordGenerator
	PasswordValidator
	PasswordChecker
}

type PasswordGenerator interface {
	New(password string) (string, error)
}

type PasswordValidator interface {
	IsValid(password string) bool
}

type PasswordChecker interface {
	Check(hash, password string) error
}
