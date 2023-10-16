package models

import "errors"

var (
	ErrBadHash          = errors.New("incorrect hash")
	ErrImageNotFound    = errors.New("image not found")
	ErrMetaDataNotFound = errors.New("metadata not found")
)
