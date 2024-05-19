package validator

import (
	"net/mail"
	"strings"
)

func NewStdlib() *Stdlib {
	return &Stdlib{}
}

// Stdlib is an extensions of built-in mail parser/validator.
// It not only looking whether string contains at symbol (@),
// but also check validity of domain.
type Stdlib struct{}

// Validate method checks if provided email is valid.
// Performs couple of simple checks on over-all mail structrure.
// Also, checks domain (subdomain) on valid structure.
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
