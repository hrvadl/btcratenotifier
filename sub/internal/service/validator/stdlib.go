package validator

import (
	"net/mail"
	"strings"
)

func NewStdlib() *Stdlib {
	return &Stdlib{}
}

type Stdlib struct{}

func (r Stdlib) Validate(email string) bool {
	m, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	parts := strings.Split(m.Address, "@")
	if len(parts) != 2 {
		return false
	}

	return len(strings.Split(parts[1], ".")) >= 2
}
