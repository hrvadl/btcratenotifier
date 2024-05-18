package validator

import "net/mail"

func NewStdlib() *Stdlib {
	return &Stdlib{}
}

type Stdlib struct{}

func (r Stdlib) Validate(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
