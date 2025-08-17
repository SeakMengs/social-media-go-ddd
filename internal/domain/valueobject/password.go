package valueobject

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	PwMaxLength = 255
	PwMinLength = 8
)

type Password struct {
	hash string
}

func NewPassword(plain string) (Password, error) {
	if strings.TrimSpace(plain) == "" {
		return Password{}, fmt.Errorf(ErrPwEmpty)
	}

	if len(plain) < PwMinLength {
		return Password{}, fmt.Errorf(ErrPwMinLength, PwMinLength)
	}

	if len(plain) > PwMaxLength {
		return Password{}, fmt.Errorf(ErrPwMaxLength, PwMaxLength)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(bytes)}, nil
}

func PasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p Password) GetHash() string {
	return p.hash
}

func (p Password) Match(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain)) == nil
}
