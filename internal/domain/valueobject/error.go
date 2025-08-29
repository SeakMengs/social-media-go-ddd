package valueobject

import (
	"errors"
	"strconv"
)

func ErrPwMaxLength(max int) error {
	return errors.New("password length must not exceed " + strconv.Itoa(max) + " character")
}

func ErrPwMinLength(min int) error {
	return errors.New("password length must at least be " + strconv.Itoa(min) + " character long")
}

var ErrPwEmpty = errors.New("password must not be empty")
