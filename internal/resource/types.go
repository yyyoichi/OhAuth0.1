package resource

import (
	"errors"
	"strings"
)

type HeaderRequest struct {
	Authorization string `header:"Authorization" binding:"required"`
}

var (
	ErrInvalidToken = errors.New("header token has invalid")
)

func (r HeaderRequest) FilterToken() (string, error) {
	if !strings.HasPrefix(r.Authorization, "Bearer") {
		return "", ErrInvalidToken
	}
	return r.Authorization[7:], nil
}

type ProfileGetResponse struct {
	UserId  string
	Name    string
	Age     uint32
	Profile string
}
