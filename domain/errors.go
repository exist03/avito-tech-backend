package domain

import "errors"

var ErrNotFound = errors.New("link not found")
var ErrNoContent = errors.New("no content")
var ErrInvalidArgument = errors.New("invalid argument")
