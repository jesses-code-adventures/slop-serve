package auth

import (
	"github.com/matthewhartstonge/argon2"
)

type Authorizer interface {
	Verify(password string, existing string) (bool, error)
	Hash(password string) (string, error)
}

type PasswordHandler struct {
	argon argon2.Config
}

func NewPasswordHandler() *PasswordHandler {
	config := argon2.DefaultConfig()
	return &PasswordHandler{config}
}

func (p *PasswordHandler) Hash(password string) (string, error) {
	encoded, err := p.argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (p *PasswordHandler) Verify(password string, existing_encoded string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(existing_encoded))
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}
