package model

import (
	"errors"
	"regexp"
	"strings"
)

func (p *Post) ValidateBeforeCreate() error {
	err := []string{}
	if strings.TrimSpace(p.Title) == "" {
		err = append(err, "title is empty")
	}
	if strings.TrimSpace(p.Body) == "" {
		err = append(err, "body is empty")
	}
	if len(err) != 0 {
		return errors.New("validation: " + strings.Join(err, ", "))
	}
	return nil
}

func (u *User) ValidateBeforeCreate() error {
	err := []string{}
	if len(u.Username) < 2 {
		err = append(err, "username must be longer then 2")
	}
	if len(u.Password) < 8 {
		err = append(err, "password must have minimum 8 symbols")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_.]+@[a-z]+\.[a-z]{2,4}$`)
	if !re.MatchString(u.Email) {
		err = append(err, "wrong email format")
	}
	if len(err) != 0 {
		return errors.New("validation: " + strings.Join(err, ", "))
	}
	return nil
}

func (u *User) ValidateBeforeUpdate() error {
	err := []string{}
	if len(u.Username) < 2 {
		err = append(err, "username must be longer then 2")
	}
	if len(err) != 0 {
		return errors.New("validation: " + strings.Join(err, ", "))
	}
	return nil
}
