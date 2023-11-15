package utils

import (
	"FIO_App/pkg/dtos"
	"errors"
	"regexp"
)

var re = regexp.MustCompile("^[a-zA-Zа-яА-Я\\s]+$")

const (
	ErrEmptyName          = "ERROR: field name shouldn't be empty"
	ErrEmptySurname       = "ERROR: field surname shouldn't be empty"
	ErrInvalidName        = "ERROR: name should contain only digits a..z and A..Z"
	ErrInvalidSurname     = "ERROR: surname should contain only digits a..z and A..Z"
	ErrInvalidPatronymic  = "ERROR: patronymic should contain only digits a..z and A..Z"
	ErrInvalidAge         = "ERROR: age should be greater 0"
	ErrInvalidGender      = "ERROR: gender can be only male or female"
	ErrInvalidNationality = "ERROR: nationality should be in this format AZ"
)

func ValidateFIO(fio dtos.FIO) error {
	if fio.Name == "" {
		return errors.New(ErrEmptyName)
	}
	if fio.Surname == "" {
		return errors.New(ErrEmptySurname)
	}
	if !re.MatchString(fio.Name) {
		return errors.New(ErrInvalidName)
	}
	if !re.MatchString(fio.Surname) {
		return errors.New(ErrInvalidSurname)
	}
	if fio.Patronymic != "" && !re.MatchString(fio.Patronymic) {
		return errors.New(ErrInvalidPatronymic)
	}
	return nil
}
