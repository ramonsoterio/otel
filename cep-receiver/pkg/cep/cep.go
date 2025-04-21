package cep

import (
	"errors"
	"regexp"
)

var (
	errInvalidCep     = errors.New("invalid zipcode")
	errCompilingRegex = errors.New("error compiling regex")
)

func IsValid(cep string) error {
	reg, err := regexp.Compile("^[0-9]{5}-?([0-9]{3}$)")
	if err != nil {
		return errCompilingRegex
	}
	if !reg.Match([]byte(cep)) {
		return errInvalidCep
	}
	return nil
}
