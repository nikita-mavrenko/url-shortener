package storage

import "errors"

var (
	ErrUrlNotFound      = errors.New("url not found")
	ErrUrlAlreadyExists = errors.New("url already exists")
)
